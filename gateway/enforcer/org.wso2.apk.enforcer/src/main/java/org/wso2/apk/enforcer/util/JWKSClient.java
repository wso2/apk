package org.wso2.apk.enforcer.util;

import com.nimbusds.jose.jwk.JWKSet;
import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.http.HttpEntity;
import org.apache.http.client.ClientProtocolException;
import org.apache.http.client.HttpClient;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.impl.client.CloseableHttpClient;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;

import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URL;
import java.nio.charset.Charset;
import java.security.KeyStore;
import java.security.cert.Certificate;
import java.text.ParseException;
import java.util.List;

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
            httpClient = FilterUtils.getHttpClient("https", null, trustStore, null);
        } catch (EnforcerException e) {
            log.error("Error occured while inializing JWKS Client", e);
            throw new EnforcerException("Error occured while inializing JWKS Client", e);
        }
    }

    public JWKSet getJWKSet() throws EnforcerException {
        try {
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
