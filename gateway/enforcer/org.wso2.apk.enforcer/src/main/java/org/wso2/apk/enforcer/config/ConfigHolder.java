/*
 * Copyright (c) 2020, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.config;

import com.nimbusds.jose.Algorithm;
import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.jwk.JWK;
import com.nimbusds.jose.jwk.KeyUse;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.util.X509CertUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.dto.JWTConfigurationDto;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.dto.APIKeyIssuerDto;
import org.wso2.apk.enforcer.config.dto.AnalyticsDTO;
import org.wso2.apk.enforcer.config.dto.AnalyticsPublisherConfigDTO;
import org.wso2.apk.enforcer.config.dto.AnalyticsReceiverConfigDTO;
import org.wso2.apk.enforcer.config.dto.AuthServiceConfigurationDto;
import org.wso2.apk.enforcer.config.dto.CacheDto;
import org.wso2.apk.enforcer.config.dto.ClientConfigDto;
import org.wso2.apk.enforcer.config.dto.FilterDTO;
import org.wso2.apk.enforcer.config.dto.ManagementCredentialsDto;
import org.wso2.apk.enforcer.config.dto.MetricsDTO;
import org.wso2.apk.enforcer.config.dto.MutualSSLDto;
import org.wso2.apk.enforcer.config.dto.SoapErrorResponseConfigDto;
import org.wso2.apk.enforcer.config.dto.ThreadPoolConfig;
import org.wso2.apk.enforcer.config.dto.TracingDTO;
import org.wso2.apk.enforcer.constants.Constants;
import org.wso2.apk.enforcer.constants.JwtConstants;
import org.wso2.apk.enforcer.discovery.config.enforcer.APIKeyEnforcer;
import org.wso2.apk.enforcer.discovery.config.enforcer.Analytics;
import org.wso2.apk.enforcer.discovery.config.enforcer.AnalyticsPublisher;
import org.wso2.apk.enforcer.discovery.config.enforcer.Cache;
import org.wso2.apk.enforcer.discovery.config.enforcer.Config;
import org.wso2.apk.enforcer.discovery.config.enforcer.Filter;
import org.wso2.apk.enforcer.discovery.config.enforcer.HttpClient;
import org.wso2.apk.enforcer.discovery.config.enforcer.JWTGenerator;
import org.wso2.apk.enforcer.discovery.config.enforcer.Keypair;
import org.wso2.apk.enforcer.discovery.config.enforcer.Management;
import org.wso2.apk.enforcer.discovery.config.enforcer.Metrics;
import org.wso2.apk.enforcer.discovery.config.enforcer.MutualSSL;
import org.wso2.apk.enforcer.discovery.config.enforcer.Service;
import org.wso2.apk.enforcer.discovery.config.enforcer.Soap;
import org.wso2.apk.enforcer.discovery.config.enforcer.Tracing;
import org.wso2.apk.enforcer.jmx.MBeanRegistrator;
import org.wso2.apk.enforcer.jwks.BackendJWKSDto;
import org.wso2.apk.enforcer.util.FilterUtils;
import org.wso2.apk.enforcer.util.JWTUtils;
import org.wso2.apk.enforcer.util.TLSUtils;

import java.io.IOException;
import java.lang.reflect.Field;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.interfaces.RSAPublicKey;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import javax.net.ssl.TrustManagerFactory;

/**
 * Configuration holder class for Microgateway.
 */
public class ConfigHolder {

    // TODO: Resolve default configs
    private static final Logger logger = LogManager.getLogger(ConfigHolder.class);

    private static ConfigHolder configHolder;
    private final EnvVarConfig envVarConfig = EnvVarConfig.getInstance();
    EnforcerConfig config = new EnforcerConfig();

    private KeyStore keyStore = null;
    private KeyStore trustStore = null;
    private KeyStore trustStoreForJWT = null;
    private KeyStore opaKeyStore = null;
    private TrustManagerFactory trustManagerFactory = null;
    private static final String dtoPackageName = EnforcerConfig.class.getPackageName();

    private ConfigHolder() {

        loadTrustStore();
        loadOpaClientKeyStore();
        loadKeyStore();
    }

    private void loadKeyStore() {
        String certPath = getEnvVarConfig().getEnforcerPublicKeyPath();
        String keyPath = getEnvVarConfig().getEnforcerPrivateKeyPath();
        keyStore = TLSUtils.getKeyStore(certPath, keyPath);
    }

    public static ConfigHolder getInstance() {

        if (configHolder != null) {
            return configHolder;
        }

        configHolder = new ConfigHolder();
        return configHolder;
    }

    /**
     * Initialize the configuration provider class by parsing the cds configuration.
     *
     * @param cdsConfig configuration fetch from CDS
     */
    public static ConfigHolder load(Config cdsConfig) {

        configHolder.parseConfigs(cdsConfig);
        return configHolder;
    }

    /**
     * Parse configurations received from the CDS to internal configuration DTO.
     * This is done inorder to prevent complicated code changes during the initial development
     * of the mgw. Later we can switch to CDS data models directly.
     */
    private void parseConfigs(Config config) {
        // load auth service
        populateAuthService(config.getAuthService());

        // Read jwt token configuration

        // Read backend jwt generation configurations
        populateJWTGeneratorConfigurations(config.getJwtGenerator());

        // Read tracing configurations
        populateTracingConfig(config.getTracing());

        // Read tracing configurations
        populateMetricsConfig(config.getMetrics());

        // Read token caching configs
        populateCacheConfigs(config.getCache());

        // Populate Analytics Configuration Values
        populateAnalyticsConfig(config.getAnalytics());

        populateMTLSConfigurations(config.getSecurity().getMutualSSL());

        populateManagementCredentials(config.getManagement());

        // Populates the SOAP error response related configs (SoapErrorInXMLEnabled).
        populateSoapErrorResponseConfigs(config.getSoap());

        // Populates the custom filter configurations applied along with enforcer filters.
        populateCustomFilters(config.getFiltersList());
        populateAPIKeyIssuer(config.getSecurity().getApiKey());
        populateInternalTokenIssuer(config.getSecurity().getRuntimeToken());
        populateMandateSubscriptionValidationConfig(config.getMandateSubscriptionValidation());
        populateMandateInternalKeyValidationConfig(config.getMandateInternalKeyValidation());
        populateHttpClientConfig(config.getHttpClient());
        // resolve string variables provided as environment variables.
        resolveConfigsWithEnvs(this.config);
    }

    private void populateHttpClientConfig(HttpClient httpClient) {

        ClientConfigDto clientConfigDto = new ClientConfigDto();
        clientConfigDto.setEnableSslVerification(httpClient.getSkipSSl());
        clientConfigDto.setHostnameVerifier(httpClient.getHostnameVerifier());
        clientConfigDto.setConnectionTimeout(httpClient.getConnectTimeout());
        clientConfigDto.setSocketTimeout(httpClient.getSocketTimeout());
        clientConfigDto.setMaxConnections(httpClient.getMaxTotalConnections());
        clientConfigDto.setMaxConnectionsPerRoute(httpClient.getMaxConnectionsPerRoute());
        config.setHttpClientConfigDto(clientConfigDto);
    }

    private void populateInternalTokenIssuer(APIKeyEnforcer runtimeToken) {

        APIKeyIssuerDto apiKeyIssuerDto = new APIKeyIssuerDto();
        apiKeyIssuerDto.setEnabled(runtimeToken.getEnabled());
        try {
            apiKeyIssuerDto.setPublicCertificate(TLSUtils.getCertificate(runtimeToken.getCertificateFilePath()));
            config.setRuntimeTokenIssuerDto(apiKeyIssuerDto);
        } catch (CertificateException | IOException e) {
            logger.error("Error occurred while configuring RuntimeToken Issuer", e);
        }
    }

    private void populateAPIKeyIssuer(APIKeyEnforcer apiKey) {

        APIKeyIssuerDto apiKeyIssuerDto = new APIKeyIssuerDto();
        apiKeyIssuerDto.setEnabled(apiKey.getEnabled());
        try {
            apiKeyIssuerDto.setPublicCertificate(TLSUtils.getCertificate(apiKey.getCertificateFilePath()));
            config.setApiKeyIssuerDto(apiKeyIssuerDto);
        } catch (CertificateException | IOException e) {
            logger.error("Error occurred while configuring APIKey Issuer", e);
        }
    }

    private void populateSoapErrorResponseConfigs(Soap soap) {

        SoapErrorResponseConfigDto soapErrorResponseConfigDto = new SoapErrorResponseConfigDto();
        soapErrorResponseConfigDto.setEnable(soap.getSoapErrorInXMLEnabled());
        config.setSoapErrorResponseConfigDto(soapErrorResponseConfigDto);
    }

    private void populateMandateSubscriptionValidationConfig(boolean mandateSubscriptionValidation) {
        config.setMandateSubscriptionValidation(mandateSubscriptionValidation);
    }

    private void populateMandateInternalKeyValidationConfig(boolean mandateInternalKeyValidation) {
        config.setMandateInternalKeyValidation(mandateInternalKeyValidation);
    }

    private void populateManagementCredentials(Management management) {

        ManagementCredentialsDto managementCredentialsDto = new ManagementCredentialsDto();
        managementCredentialsDto.setPassword(management.getPassword().toCharArray());
        managementCredentialsDto.setUserName(management.getUsername());
        config.setManagement(managementCredentialsDto);
    }

    private void populateMTLSConfigurations(MutualSSL mtlsInfo) {

        MutualSSLDto mutualSSLDto = new MutualSSLDto();
        mutualSSLDto.setCertificateHeader(mtlsInfo.getCertificateHeader());
        mutualSSLDto.setEnableClientValidation(mtlsInfo.getEnableClientValidation());
        mutualSSLDto.setClientCertificateEncode(mtlsInfo.getClientCertificateEncode());
        mutualSSLDto.setEnableOutboundCertificateHeader(mtlsInfo.getEnableOutboundCertificateHeader());
        config.setMtlsInfo(mutualSSLDto);
    }

    private void populateAuthService(Service cdsAuth) {

        AuthServiceConfigurationDto authDto = new AuthServiceConfigurationDto();
        authDto.setKeepAliveTime(cdsAuth.getKeepAliveTime());
        authDto.setPort(cdsAuth.getPort());
        authDto.setMaxHeaderLimit(cdsAuth.getMaxHeaderLimit());
        authDto.setMaxMessageSize(cdsAuth.getMaxMessageSize());

        ThreadPoolConfig threadPool = new ThreadPoolConfig();
        MBeanRegistrator.registerMBean(threadPool);

        threadPool.setCoreSize(cdsAuth.getThreadPool().getCoreSize());
        threadPool.setKeepAliveTime(cdsAuth.getThreadPool().getKeepAliveTime());
        threadPool.setMaxSize(cdsAuth.getThreadPool().getMaxSize());
        threadPool.setQueueSize(cdsAuth.getThreadPool().getQueueSize());
        authDto.setThreadPool(threadPool);

        config.setAuthService(authDto);
    }

    private void populateTracingConfig(Tracing tracing) {

        TracingDTO tracingConfig = new TracingDTO();
        tracingConfig.setTracingEnabled(tracing.getEnabled());
        tracingConfig.setExporterType(tracing.getType());
        tracingConfig.setConfigProperties(tracing.getConfigPropertiesMap());
        config.setTracingConfig(tracingConfig);
    }

    private void populateMetricsConfig(Metrics metrics) {

        MetricsDTO metricsConfig = new MetricsDTO();
        metricsConfig.setMetricsEnabled(metrics.getEnabled());
        metricsConfig.setMetricsType(metrics.getType());
        config.setMetricsConfig(metricsConfig);
    }

    private void loadTrustStore() {

        try {

            trustStore = KeyStore.getInstance(KeyStore.getDefaultType());
            trustStore.load(null);

            if (getEnvVarConfig().isTrustDefaultCerts()) {
                TLSUtils.loadDefaultCertsToTrustStore(trustStore);
            }
            loadTrustedCertsToTrustStore();

            trustManagerFactory = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
            trustManagerFactory.init(trustStore);

        } catch (KeyStoreException | CertificateException | NoSuchAlgorithmException | IOException e) {
            logger.error("Error in loading certs to the trust store.", e);
        }
    }

    private void loadTrustedCertsToTrustStore() throws IOException {

        String truststoreFilePath = getEnvVarConfig().getTrustedAdapterCertsPath();
        TLSUtils.addCertsToTruststore(trustStore, truststoreFilePath);
    }

    private void loadOpaClientKeyStore() {

        String certPath = getEnvVarConfig().getOpaClientPublicKeyPath();
        String keyPath = getEnvVarConfig().getOpaClientPrivateKeyPath();
        opaKeyStore = FilterUtils.createClientKeyStore(certPath, keyPath);
    }

    private void populateJWTGeneratorConfigurations(JWTGenerator jwtGenerator) {

        List<Keypair> keypairs = jwtGenerator.getKeypairsList();
        // Validation is done at the adapter to ensure that only one signing keypair is available
        Keypair signingKey = getSigningKey(keypairs);
        JWTConfigurationDto jwtConfigurationDto = new JWTConfigurationDto();
        try {
            jwtConfigurationDto.setPublicCert(TLSUtils.getCertificate(signingKey.getPublicCertificatePath()));
            jwtConfigurationDto.setPrivateKey(JWTUtils.getPrivateKey(signingKey.getPrivateKeyPath()));
        } catch (EnforcerException | CertificateException | IOException e) {
            String err = "Error in loading keypair for Backend JWTs: " + e;
            logger.error(err);
        }
        jwtConfigurationDto.setUseKid(true);
        config.setJwtConfigurationDto(jwtConfigurationDto);
        populateBackendJWKSConfiguration(jwtGenerator);
    }

    public KeyStore getKeyStore() {

        return keyStore;
    }

    public void setKeyStore(KeyStore keyStore) {

        this.keyStore = keyStore;
    }

    private void populateBackendJWKSConfiguration(JWTGenerator jwtGenerator) {

        BackendJWKSDto backendJWKSDto = new BackendJWKSDto();
        List<Keypair> keypairs = jwtGenerator.getKeypairsList();
        ArrayList<JWK> jwks = new ArrayList<>();
        try {
            for (Keypair keypair : keypairs) {
                X509Certificate cert = X509CertUtils
                        .parse(TLSUtils.getCertificate(keypair.getPublicCertificatePath()).getEncoded());
                RSAPublicKey publicKey = RSAKey.parse(cert).toRSAPublicKey();
                RSAKey jwk = new RSAKey.Builder(publicKey)
                        .keyUse(KeyUse.SIGNATURE)
                        .algorithm(JWSAlgorithm.RS256)
                        .keyIDFromThumbprint()
                        .build().toPublicJWK();
                jwks.add(jwk);
                if (keypair.getUseForSigning()) {
                    config.getJwtConfigurationDto().setKidValue(jwk.getKeyID());
                }
            }
        } catch (JOSEException | CertificateException | IOException e) {
            logger.error("Error in loading additional public certificates for JWKS: " + e);
        }
        backendJWKSDto.setJwks(jwks);
        config.setBackendJWKSDto(backendJWKSDto);
    }

    private Keypair getSigningKey(List<Keypair> keypairs) {

        for (Keypair keypair : keypairs) {
            if (keypair.getUseForSigning()) {
                return keypair;
            }
        }
        return null;
    }

    private Algorithm getJWKSAlgorithm(String alg) {

        switch (alg) {
            case JwtConstants.RS384:
                return JWSAlgorithm.RS384;
            case JwtConstants.RS512:
                return JWSAlgorithm.RS512;
            default:
                return JWSAlgorithm.RS256;
        }
    }

    private void populateCacheConfigs(Cache cache) {

        CacheDto cacheDto = new CacheDto();
        cacheDto.setEnabled(cache.getEnable());
        cacheDto.setMaximumSize(cache.getMaximumSize());
        cacheDto.setExpiryTime(cache.getExpiryTime());
        config.setCacheDto(cacheDto);
    }

    private void populateAnalyticsConfig(Analytics analyticsConfig) {

        AnalyticsReceiverConfigDTO serverConfig = new AnalyticsReceiverConfigDTO();
        serverConfig.setKeepAliveTime(analyticsConfig.getService().getKeepAliveTime());
        serverConfig.setMaxHeaderLimit(analyticsConfig.getService().getMaxHeaderLimit());
        serverConfig.setMaxMessageSize(analyticsConfig.getService().getMaxMessageSize());
        serverConfig.setPort(analyticsConfig.getService().getPort());

        ThreadPoolConfig threadPoolConfig = new ThreadPoolConfig();
        threadPoolConfig.setCoreSize(analyticsConfig.getService().getThreadPool().getCoreSize());
        threadPoolConfig.setMaxSize(analyticsConfig.getService().getThreadPool().getMaxSize());
        threadPoolConfig.setKeepAliveTime(analyticsConfig.getService().getThreadPool().getKeepAliveTime());
        threadPoolConfig.setQueueSize(analyticsConfig.getService().getThreadPool().getQueueSize());
        serverConfig.setThreadPoolConfig(threadPoolConfig);

        AnalyticsDTO analyticsDTO = new AnalyticsDTO();
        analyticsDTO.setServerConfig(serverConfig);
        analyticsDTO.setEnabled(analyticsConfig.getEnabled());
        Map<String, String> propertiesMap = analyticsConfig.getPropertiesMap();
        Map<String, Object> resolvedProperties = new HashMap<>();
        for (Map.Entry<String, String> propertiesEntry : propertiesMap.entrySet()) {
            resolvedProperties.put(propertiesEntry.getKey(), getEnvValue(propertiesEntry.getValue()));
        }
        analyticsDTO.setProperties(resolvedProperties);
        for (AnalyticsPublisher analyticsPublisher : analyticsConfig.getAnalyticsPublisherList()) {
            Map<String, String> resolvedConfigMap = new HashMap<>();
            Map<String, String> configPropertiesMap = analyticsPublisher.getConfigPropertiesMap();
            for (Map.Entry<String, String> config : configPropertiesMap.entrySet()) {
                resolvedConfigMap.put(config.getKey(), getEnvValue(config.getValue()).toString());
            }
            analyticsDTO.addAnalyticsPublisherConfig(new AnalyticsPublisherConfigDTO(analyticsPublisher.getEnabled(),
                    analyticsPublisher.getType(), resolvedConfigMap));
        }
        config.setAnalyticsConfig(analyticsDTO);

    }

    /**
     * This method recursively looks for the string type config values in the {@link EnforcerConfig} object ,
     * which have the prefix `$env{` and reads the respective value from the environment variable and set it to
     * the config object.
     *
     * @param config - Enforcer config object.
     */
    private void resolveConfigsWithEnvs(Object config) {
        //extended config class env variables should also be resolved
        if (config.getClass().getSuperclass() != null && (
                config.getClass().getSuperclass().getPackageName().contains(dtoPackageName))) {
            processRecursiveObject(config, config.getClass().getSuperclass().getDeclaredFields());
        }
        processRecursiveObject(config, config.getClass().getDeclaredFields());
    }

    private void processRecursiveObject(Object config, Field[] classFields) {

        for (Field field : classFields) {
            try {
                field.setAccessible(true);
                // handles the string and char array objects
                if (field.getType().isAssignableFrom(String.class) || field.getType().isAssignableFrom(char[].class)) {
                    field.set(config, getEnvValue(field.get(config)));
                } else if (field.getName().contains(Constants.OBJECT_THIS_NOTATION)) {
                    // skip the java internal objects created, inside the recursion to avoid stack overflow.
                    continue;
                } else if (Map.class.isAssignableFrom(field.getType())) {
                    // handles the config objects saved as Maps
                    Map<Object, Object> objectMap = (Map<Object, Object>) field.get(config);
                    for (Map.Entry<Object, Object> entry : objectMap.entrySet()) {
                        if (entry.getValue().getClass().isAssignableFrom(String.class) || entry.getValue().getClass()
                                .isAssignableFrom(char[].class)) {
                            field.set(config, getEnvValue(field.get(config)));
                            continue;
                        } else if (entry.getValue().getClass().getPackageName().contains(dtoPackageName)) {
                            resolveConfigsWithEnvs(entry.getValue());
                        }
                    }
                } else if (field.getType().isArray() && field.getType().getPackageName().contains(dtoPackageName)) {
                    // handles the config objects saved as arrays
                    Object[] objectArray = (Object[]) field.get(config);
                    for (Object arrayObject : objectArray) {
                        if (arrayObject.getClass().getPackageName().contains(dtoPackageName)) {
                            resolveConfigsWithEnvs(arrayObject);
                        } else if (arrayObject.getClass().isAssignableFrom(String.class) || arrayObject.getClass()
                                .isAssignableFrom(char[].class)) {
                            field.set(config, getEnvValue(arrayObject));
                        }
                    }
                } else if (field.getType().getPackageName().contains(dtoPackageName)) { //recursively call the dto
                    // objects in the same package
                    resolveConfigsWithEnvs(field.get(config));
                }
            } catch (IllegalAccessException e) {
                //log and continue
                logger.error("Error while reading the config value : " + field.getName(), e);
            }
        }
    }

    private Object getEnvValue(Object configValue) {

        if (configValue instanceof String) {
            String value = (String) configValue;
            return replaceEnvRegex(value);
        } else if (configValue instanceof char[]) {
            String value = String.valueOf((char[]) configValue);
            return replaceEnvRegex(value).toCharArray();
        }
        return configValue;
    }

    private String replaceEnvRegex(String value) {

        Matcher m = Pattern.compile("\\$env\\{(.*?)\\}").matcher(value);
        if (value.contains(Constants.ENV_PREFIX)) {
            while (m.find()) {
                String envName = value.substring(m.start() + 5, m.end() - 1);
                if (System.getenv(envName) != null) {
                    value = value.replace(value.substring(m.start(), m.end()), System.getenv(envName));
                }
            }
        }
        return value;
    }

    private void populateCustomFilters(List<Filter> filterList) {

        FilterDTO[] filterArray = new FilterDTO[filterList.size()];
        int index = 0;
        for (Filter filter : filterList) {
            FilterDTO filterDTO = new FilterDTO();
            filterDTO.setClassName(filter.getClassName());
            filterDTO.setPosition(filter.getPosition());
            filterDTO.setConfigProperties(filter.getConfigPropertiesMap());
            filterArray[index] = filterDTO;
            index++;
        }
        config.setCustomFilters(filterArray);
    }

    public EnforcerConfig getConfig() {

        return config;
    }

    public void setConfig(EnforcerConfig config) {

        this.config = config;
    }

    public KeyStore getTrustStore() {

        return trustStore;
    }

    public KeyStore getTrustStoreForJWT() {

        return trustStoreForJWT;
    }

    public KeyStore getOpaKeyStore() {

        return opaKeyStore;
    }

    public TrustManagerFactory getTrustManagerFactory() {

        return trustManagerFactory;
    }

    public EnvVarConfig getEnvVarConfig() {

        return envVarConfig;
    }

}
