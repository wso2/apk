/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.server;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;
import org.wso2.apk.enforcer.util.JWTUtils;
import org.wso2.apk.enforcer.util.TLSUtils;
import redis.clients.jedis.DefaultJedisClientConfig;
import redis.clients.jedis.HostAndPort;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisClientConfig;
import redis.clients.jedis.JedisPool;
import redis.clients.jedis.JedisPubSub;
import redis.clients.jedis.params.ScanParams;
import redis.clients.jedis.resps.ScanResult;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.nio.file.Paths;
import java.security.GeneralSecurityException;
import java.security.Key;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.util.Collections;
import java.util.HashSet;
import java.util.Map;
import java.util.Queue;
import java.util.Set;
import java.util.concurrent.Executors;
import java.util.concurrent.PriorityBlockingQueue;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import javax.net.ssl.KeyManagerFactory;
import javax.net.ssl.ManagerFactoryParameters;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSocketFactory;
import javax.net.ssl.TrustManagerFactory;

public class RevokedTokenRedisClient {

    private JedisPool jedisPool;
    private Set<String> revokedTokens;
    private Queue<Map.Entry<Long, String>> expiryQueue;
    private static Set<String> revokedTokensStatic;
    private String redisRevokedTokensChannel;
    private final ScheduledExecutorService revokedTokensCleanupScheduler = Executors.newScheduledThreadPool(1);
    private int revokedTokenCleanupInterval;
    private static volatile boolean isAlreadyStarted = false;
    private static final Logger logger = LogManager.getLogger(RevokedTokenRedisClient.class);
    private static final String TOKEN_EXPIRY_DIVIDER = "_##_";
    private static final String REVOKED_TOKEN_REDIS_KEY_PATTERN = "wso2:apk:revoked_token:*";
    private RevokedTokenRedisClient(Set<String> revokedTokens, Queue<Map.Entry<Long, String>> expiryQueue) throws EnforcerException {
        this.revokedTokens = revokedTokens;
        this.expiryQueue = expiryQueue;

        String userName = ConfigHolder.getInstance().getEnvVarConfig().getRedisUsername();
        String password = ConfigHolder.getInstance().getEnvVarConfig().getRedisPassword();
        String host = ConfigHolder.getInstance().getEnvVarConfig().getRedisHost();
        int port = ConfigHolder.getInstance().getEnvVarConfig().getRedisPort();
        boolean isSSLEnabled = ConfigHolder.getInstance().getEnvVarConfig().isRedisTlsEnabled();
        this.redisRevokedTokensChannel = ConfigHolder.getInstance().getEnvVarConfig().getRevokedTokensRedisChannel();
        this.revokedTokenCleanupInterval = ConfigHolder.getInstance().getEnvVarConfig().getRevokedTokenCleanupInterval();
        String caCert = ConfigHolder.getInstance().getEnvVarConfig().getRedisCaCertFile();
        DefaultJedisClientConfig.Builder builder = DefaultJedisClientConfig.builder()
                .password(password);
        if (!userName.isEmpty()) {
            builder.user(userName);
        }
        if (isSSLEnabled) {
            SSLSocketFactory sslFactory = createSslSocketFactory(caCert);
            builder = builder
                    .ssl(true)
                    .sslSocketFactory(sslFactory);
        }
        JedisClientConfig config = builder.build();

        HostAndPort hostAndPort = new HostAndPort(host, port);
        this.jedisPool = new JedisPool(hostAndPort, config);

    }

    public static void retrieveAndSubscribe() throws EnforcerException {
        if (isAlreadyStarted) {
            logger.debug("Already a token retreival task is running");
            return;
        }
        logger.debug("Starting redis revoked token client...");
        isAlreadyStarted = true;

        HashSet<String> revokedTokens = new HashSet<>();
        Set<String> synchronizedRevokedTokens = Collections.synchronizedSet(revokedTokens);
        PriorityBlockingQueue<Map.Entry<Long, String>> expiryQueue =
                new PriorityBlockingQueue<>(10, Map.Entry.comparingByKey());
        RevokedTokenRedisClient revokedTokenRedisClient =
                new RevokedTokenRedisClient(synchronizedRevokedTokens, expiryQueue);

        revokedTokensStatic = revokedTokens;
        revokedTokenRedisClient.subscribe();
        revokedTokenRedisClient.retrieveAllRevokedTokens();
        revokedTokenRedisClient.scheduleCleanup();
    }

    private void scheduleCleanup() {
        this.revokedTokensCleanupScheduler.scheduleAtFixedRate(this::startCleanupTask, this.revokedTokenCleanupInterval,
                this.revokedTokenCleanupInterval, TimeUnit.SECONDS);
    }

    private void subscribe() {
        Thread jedisThread = new Thread(new RevokedTokenRedisSubscriber(this.jedisPool,
                this.revokedTokens, this.expiryQueue, this.redisRevokedTokensChannel));
        jedisThread.start();
    }

    private void retrieveAllRevokedTokens() {
        try (Jedis jedis = this.jedisPool.getResource()) {
            String cursor = "0";
            Set<String> keysAndValues = new HashSet<>();
            do {
                ScanResult<String> scanResult = jedis.scan(cursor,
                        new ScanParams().match(REVOKED_TOKEN_REDIS_KEY_PATTERN));
                keysAndValues.addAll(scanResult.getResult());
                cursor = scanResult.getCursor();
            } while (!cursor.equals("0"));

            for (String key : keysAndValues) {
                try {
                    String value = jedis.get(key);
                    Long expiry = Long.valueOf(value);
                    String token = key.substring(REVOKED_TOKEN_REDIS_KEY_PATTERN.length()-1);
                    revokedTokens.add(token);
                    expiryQueue.offer(Map.entry(expiry, token));
                    logger.debug("New token added. Token : " + token + " expiry: " + expiry);
                } catch(Exception e) {
                    logger.warn("Error while processing key: " + key, e);
                }
            }
        } catch (Exception e) {
            logger.error("Error while creating redis connection.", e);
        }
    }

    private void startCleanupTask() {
        long currentTime = System.currentTimeMillis() / 1000L;
        while (!this.expiryQueue.isEmpty() && this.expiryQueue.peek().getKey() <= currentTime) {
            Map.Entry<Long, String> entry = this.expiryQueue.poll();
            String token = entry.getValue();
            this.revokedTokens.remove(token);
            logger.debug("Token removed: " + token + " expiry: " + entry.getKey());
        }
    }

    static class RevokedTokenRedisSubscriber implements Runnable {
        JedisPool jedisPool;
        Set<String> revokedTokens;
        Queue<Map.Entry<Long, String>> expiryQueue;
        private String redisRevokedTokensChannel;

        public RevokedTokenRedisSubscriber(JedisPool pool,
                                           Set<String> revokedTokens,
                                           Queue<Map.Entry<Long, String>> expiryQueue,
                                           String channel) {
            this.jedisPool = pool;
            this.revokedTokens = revokedTokens;
            this.expiryQueue = expiryQueue;
            this.redisRevokedTokensChannel = channel;
        }
        @Override
        public void run() {
            try(Jedis jedis = this.jedisPool.getResource()) {
                jedis.connect();
                JedisPubSub jedisPubSub = new JedisPubSub() {
                    @Override
                    public void onMessage(String channel, String message) {
                        try {
                            logger.debug("Received message: " + message);
                            String[] tokenAndExpiry = message.split(TOKEN_EXPIRY_DIVIDER);
                            Long expiry = Long.valueOf(tokenAndExpiry[1]);
                            String token = tokenAndExpiry[0];
                            revokedTokens.add(token);
                            expiryQueue.offer(Map.entry(expiry, token));
                        } catch (Exception e) {
                            logger.error("Error while processing the token message in the redis " +
                                    "subscriber. Exception: ", e);
                        }
                    }

                    @Override
                    public void onUnsubscribe(String channel, int subscribedChannels) {
                        logger.info("Unsubscribed Channel: {} subscribed channels: {}", channel, subscribedChannels);
                    }

                    @Override
                    public void onSubscribe(String channel, int subscribedChannels) {
                        logger.info("Subscribed Channel: {} subscribed channels: {}", channel, subscribedChannels);
                    }
                };
                jedis.subscribe(jedisPubSub, this.redisRevokedTokensChannel);
            } catch (Exception e) {
                logger.error("Error occured in the subscription connection trying to connect again."+ e);
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException ex) {
                    logger.error("Exception while sleeping. Exception: " + ex);
                }
                run();
            }

        }
    }

    public static Set<String> getRevokedTokens() {
        return revokedTokensStatic;
    }

    public static void setRevokedTokens(Set<String>revokedTokensSet) {
        revokedTokensStatic = revokedTokensSet;
    }

    private static SSLSocketFactory createSslSocketFactory(String redisCaCertPath) throws EnforcerException {
        try {
            KeyStore trustStore = TLSUtils.getDefaultCertTrustStore();
            TLSUtils.addCertsToTruststore(trustStore, redisCaCertPath);
            TrustManagerFactory trustManagerFactory = TrustManagerFactory.getInstance("X509");
            trustManagerFactory.init(trustStore);
            Certificate cert = TLSUtils.getCertificateFromFile(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPublicKeyPath());
            Key key = JWTUtils.getPrivateKey(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPrivateKeyPath());
            KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
            keyStore.load(null, null);
            keyStore.setKeyEntry("client-keys", key, null, new Certificate[]{cert});
            KeyManagerFactory keyManagerFactory = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm());
            keyManagerFactory.init(keyStore, null);
            SSLContext sslContext = SSLContext.getInstance("TLS");
            sslContext.init(keyManagerFactory.getKeyManagers(), trustManagerFactory.getTrustManagers(), null);
            return sslContext.getSocketFactory();
        } catch (Exception e) {
            throw new EnforcerException("Error while creating SSL socket factory.", e);
        }
    }

}
