/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.enforcer.util;

import com.nimbusds.jose.jwk.JWKSet;
import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.http.HttpEntity;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;

import java.io.IOException;
import java.io.InputStream;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.text.ParseException;
import java.util.List;

/**
 * This class used to create JWKS Client.
 */
public class JWKSClient {
    private static final Log log = LogFactory.getLog(JWKSClient.class.getName());
    private HttpClient httpClient;
    private String jwksEndpoint;

    public JWKSClient(String jwksEndpoint, List<Certificate> certificates) throws EnforcerException {
        this.jwksEndpoint = jwksEndpoint;
        try {
            KeyStore trustStore = ConfigHolder.getInstance().getTrustStore();
            if (certificates.size() > 0) {
                trustStore = TLSUtils.getDefaultCertTrustStore();
            }
            TLSUtils.convertAndAddCertificatesToTrustStore(trustStore, certificates);
            httpClient = FilterUtils.getHttpClient(null, trustStore, null);
        } catch (EnforcerException e) {
            log.error("Error occured while inializing JWKS Client", e);
            throw new EnforcerException("Error occured while inializing JWKS Client", e);
        }
    }

    public JWKSet getJWKSet() throws EnforcerException {
        try {
            System.out.print(jwksEndpoint + "haha");
            HttpGet httpGet = new HttpGet(jwksEndpoint);
            try (CloseableHttpResponse response = (CloseableHttpResponse) httpClient.execute(httpGet)) {
                if (response.getStatusLine().getStatusCode() == 200) {
                    HttpEntity entity = response.getEntity();
                    try (InputStream content = entity.getContent()) {
                        String stringContent = IOUtils.toString(content);
                        return JWKSet.parse(stringContent);
                    }
                } else {
                    throw new EnforcerException("Error occurred when calling JWKS Endpoint");
                }
            }
        } catch (IOException | ParseException e) {
            throw new EnforcerException("Error occurred when calling JWKS Endpoint", e);
        }
    }
}
