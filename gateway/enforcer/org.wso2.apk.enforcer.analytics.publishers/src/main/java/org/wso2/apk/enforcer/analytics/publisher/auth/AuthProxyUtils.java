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

package org.wso2.apk.enforcer.analytics.publisher.auth;

import feign.Client;
import feign.httpclient.ApacheHttpClient;
import org.apache.commons.lang3.StringUtils;
import org.apache.http.HttpHost;
import org.apache.http.auth.AuthScope;
import org.apache.http.auth.UsernamePasswordCredentials;
import org.apache.http.client.CredentialsProvider;
import org.apache.http.config.RegistryBuilder;
import org.apache.http.conn.socket.ConnectionSocketFactory;
import org.apache.http.conn.socket.PlainConnectionSocketFactory;
import org.apache.http.conn.ssl.DefaultHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.SSLContexts;
import org.apache.http.impl.client.BasicCredentialsProvider;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.client.ProxyAuthenticationStrategy;
import org.apache.http.impl.conn.DefaultProxyRoutePlanner;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.wso2.apk.enforcer.analytics.publisher.exception.HttpClientException;
import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyManagementException;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.util.Map;
import javax.net.ssl.SSLContext;

import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.HTTPS_PROTOCOL;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.HTTP_PROTOCOL;
import static org.wso2.apk.enforcer.analytics.publisher.util.Constants.KEYSTORE_TYPE;

/**
 * Util class to generate http client with proxy configurations.
 */
public class AuthProxyUtils {

    private static final Logger log = LoggerFactory.getLogger(AuthProxyUtils.class);

    public static Client getClient(Map<String, String> properties) {
        return getFeignHttpClient(properties);
    }

    private static ApacheHttpClient getFeignHttpClient(Map<String, String> properties) {
        String proxyHost = properties.get(Constants.PROXY_HOST);
        int proxyPort = Integer.parseInt(properties.get(Constants.PROXY_PORT));
        String proxyUsername = properties.get(Constants.PROXY_USERNAME);
        String proxyPassword = properties.get(Constants.PROXY_PASSWORD);
        String proxyProtocol = properties.get(Constants.PROXY_PROTOCOL);

        if (StringUtils.isEmpty(proxyProtocol)) {
            proxyProtocol = HTTP_PROTOCOL;
        }

        PoolingHttpClientConnectionManager pool = null;
        try {
            pool = getPoolingHttpClientConnectionManager(properties);
        } catch (HttpClientException e) {
            log.error("Error while getting http client connection manager", e);
        }

        HttpHost host = new HttpHost(proxyHost, proxyPort, proxyProtocol);
        DefaultProxyRoutePlanner routePlanner = new DefaultProxyRoutePlanner(host);
        HttpClientBuilder clientBuilder = HttpClientBuilder.create()
                .setRoutePlanner(routePlanner)
                .setConnectionManager(pool);

        if (!StringUtils.isBlank(proxyUsername) && !StringUtils.isBlank(proxyPassword)) {
            CredentialsProvider credentialsProvider = new BasicCredentialsProvider();
            credentialsProvider.setCredentials(new AuthScope(proxyHost, proxyPort),
                    new UsernamePasswordCredentials(proxyUsername, proxyPassword));
            clientBuilder
                    .setProxyAuthenticationStrategy(new ProxyAuthenticationStrategy())
                    .setDefaultCredentialsProvider(credentialsProvider);
        }

        return new ApacheHttpClient(clientBuilder.build());
    }

    private static PoolingHttpClientConnectionManager getPoolingHttpClientConnectionManager(Map<String,
            String> properties) throws HttpClientException {
        SSLConnectionSocketFactory socketFactory = createSocketFactory(properties);
        ConnectionSocketFactory httpSocketFactory = new PlainConnectionSocketFactory();
        org.apache.http.config.Registry<ConnectionSocketFactory> socketFactoryRegistry =
                RegistryBuilder.<ConnectionSocketFactory>create()
                        .register(HTTP_PROTOCOL, httpSocketFactory)
                        .register(HTTPS_PROTOCOL, socketFactory)
                        .build();
        return new PoolingHttpClientConnectionManager(socketFactoryRegistry);
    }

    private static SSLConnectionSocketFactory createSocketFactory(Map<String, String> properties)
            throws HttpClientException {
        SSLContext sslContext;

        String keyStorePassword = properties.get(Constants.KEYSTORE_PASSWORD);
        String keyStorePath = properties.get(Constants.KEYSTORE_LOCATION);
        try {
            KeyStore trustStore = KeyStore.getInstance(KEYSTORE_TYPE);
            trustStore.load(Files.newInputStream(Paths.get(keyStorePath)), keyStorePassword.toCharArray());
            sslContext = SSLContexts.custom().loadTrustMaterial(trustStore).build();
            return new SSLConnectionSocketFactory(sslContext, new DefaultHostnameVerifier());
        } catch (KeyStoreException e) {
            throw new HttpClientException("Failed to read from Key Store", e);
        } catch (IOException e) {
            throw new HttpClientException("Key Store not found in " + keyStorePath, e);
        } catch (CertificateException e) {
            throw new HttpClientException("Failed to read Certificate", e);
        } catch (NoSuchAlgorithmException e) {
            throw new HttpClientException("Failed to load Key Store from " + keyStorePath, e);
        } catch (KeyManagementException e) {
            throw new HttpClientException("Failed to load key from" + keyStorePath, e);
        }
    }
}
