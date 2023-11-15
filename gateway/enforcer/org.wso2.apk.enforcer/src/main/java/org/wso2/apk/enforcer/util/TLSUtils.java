/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package org.wso2.apk.enforcer.util;

import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.shaded.io.netty.handler.ssl.ClientAuth;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import org.apache.commons.io.IOUtils;
import org.apache.commons.lang3.RandomStringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.wso2.apk.enforcer.commons.exception.EnforcerException;
import org.wso2.apk.enforcer.config.ConfigHolder;

import javax.net.ssl.SSLException;
import javax.net.ssl.TrustManager;
import javax.net.ssl.TrustManagerFactory;
import javax.net.ssl.X509TrustManager;
import java.io.*;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.Key;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Utility Functions related to TLS Certificates.
 */
public class TLSUtils {

    private static final Logger log = LogManager.getLogger(TLSUtils.class);
    private static final String X509 = "X.509";
    private static final String crtExtension = ".crt";
    private static final String pemExtension = ".pem";
    private static final String endCertificateDelimiter = "-----END CERTIFICATE-----";

    /**
     * Read the certificate file and return the certificate.
     *
     * @param filePath Filepath of the corresponding certificate
     * @return Certificate
     */
    public static Certificate getCertificateFromFile(String filePath)
            throws CertificateException, IOException, EnforcerException {

        return getCertsFromFile(filePath, true).get(0);
    }

    /**
     * Read the pem encoded certificate content and generate certificate.
     *
     * @param certificateContent Pem Encoded certificate Content
     * @return Certificate
     */
    public static Certificate getCertificateFromString(String certificateContent)
            throws CertificateException, IOException {
        // A single certificate file is expected
        try (InputStream inputStream = new ByteArrayInputStream(certificateContent.getBytes())) {
            CertificateFactory fact = CertificateFactory.getInstance(X509);
            return fact.generateCertificate(inputStream);
        }
    }

    /**
     * Add the certificates to the truststore.
     *
     * @param filePath   Filepath of the corresponding certificate or directory containing the certificates
     * @param trustStore Keystore with trusted certificates
     */
    public static void addCertsToTruststore(KeyStore trustStore, String filePath) throws IOException {

        if (!Files.exists(Paths.get(filePath))) {
            log.error("The provided certificates directory/file path does not exist. : " + filePath);
            return;
        }
        if (Files.isDirectory(Paths.get(filePath))) {
            log.debug("Provided Path is a directory: " + filePath);
            Files.walk(Paths.get(filePath)).filter(path -> {
                Path fileName = path.getFileName();
                return fileName != null && (fileName.toString().endsWith(crtExtension) ||
                        fileName.toString().endsWith(pemExtension));
            }).forEach(path -> {
                updateTruststoreWithMultipleCertPem(trustStore, path.toAbsolutePath().toString());
            });
        } else {
            log.debug("Provided Path is a regular File Path : " + filePath);
            updateTruststoreWithMultipleCertPem(trustStore, filePath);
        }
    }

    public static void convertAndAddCertificatesToTrustStore(KeyStore trustStore, List<Certificate> certificates) {

        for (Certificate certificate : certificates) {
            try {
                trustStore.setCertificateEntry(RandomStringUtils.random(10, true, false),
                        certificate);
            } catch (KeyStoreException e) {
                log.error("Error while adding the trusted certificates to the trustStore.", e);
            }
        }
    }

    private static List<Certificate> getCertsFromFile(String filepath, boolean restrictToOne)
            throws CertificateException, IOException, EnforcerException {

        String content = new String(Files.readAllBytes(Paths.get(filepath)));

        if (!content.contains(endCertificateDelimiter)) {
            throw new EnforcerException("Content Provided within the certificate file:" + filepath + "is invalid.");
        }

        int endIndex = content.lastIndexOf(endCertificateDelimiter) + endCertificateDelimiter.length();
        // If there are any additional characters afterwards,
        if (endIndex < content.length()) {
            content = content.substring(0, endIndex);
        }

        List<Certificate> certList = new ArrayList<>();
        CertificateFactory cf = CertificateFactory.getInstance(X509);
        InputStream inputStream = new ByteArrayInputStream(content.getBytes());
        BufferedInputStream bufferedInputStream = new BufferedInputStream(inputStream);
        int count = 1;
        while (bufferedInputStream.available() > 0) {
            if (count > 1 && restrictToOne) {
                log.warn("Provided PEM file " + filepath +
                        "contains more than one certificate. Hence proceeding with" +
                        "the first certificate in the File for the JWT configuraion related certificate.");
                return certList;
            }
            Certificate cert = cf.generateCertificate(bufferedInputStream);
            certList.add(cert);
            count++;
        }
        return certList;
    }

    private static void updateTruststoreWithMultipleCertPem(KeyStore trustStore, String filePath) {

        try {
            List<Certificate> certificateList = getCertsFromFile(filePath, false);
            certificateList.forEach(certificate -> {
                try {
                    trustStore.setCertificateEntry(RandomStringUtils.random(10, true, false),
                            certificate);
                } catch (KeyStoreException e) {
                    log.error("Error while adding the trusted certificates to the trustStore.", e);
                }
            });
            log.debug("Certificate Added to the truststore : " + filePath);
        } catch (CertificateException | IOException | EnforcerException e) {
            log.error("Error while adding certificates to the truststore.", e);
        }
    }

    public static Certificate getCertificate(String filePath) throws CertificateException, IOException {

        try (FileInputStream fileInputStream = new FileInputStream(filePath)) {
            String content = IOUtils.toString(fileInputStream);
            return getCertificateFromContent(content);
        }
    }

    public static Certificate getCertificateFromContent(String content) throws CertificateException, IOException {

        CertificateFactory fact = CertificateFactory.getInstance(X509);
        try (InputStream is = new ByteArrayInputStream(content.getBytes())) {
            X509Certificate cert = (X509Certificate) fact.generateCertificate(is);
            return cert;
        }
    }

    /**
     * Generate the gRPC Server SSL Context where the mutual SSL is also enabled.
     *
     * @return {@code SsLContext} generated SSL Context
     * @throws SSLException
     */
    public static SslContext buildGRPCServerSSLContext() throws SSLException {

        File certFile = Paths.get(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPublicKeyPath()).toFile();
        File keyFile = Paths.get(ConfigHolder.getInstance().getEnvVarConfig().getEnforcerPrivateKeyPath()).toFile();

        return GrpcSslContexts.forServer(certFile, keyFile)
                .trustManager(ConfigHolder.getInstance().getTrustManagerFactory())
                .clientAuth(ClientAuth.REQUIRE)
                .build();
    }

    public static KeyStore getDefaultCertTrustStore() throws EnforcerException {

        try {

            KeyStore trustStore = KeyStore.getInstance(KeyStore.getDefaultType());
            trustStore.load(null);
            loadDefaultCertsToTrustStore(trustStore);
            return trustStore;
        } catch (KeyStoreException | CertificateException | NoSuchAlgorithmException | IOException e) {
            throw new EnforcerException("Error while generating Default trustStore", e);
        }
    }

    public static void loadDefaultCertsToTrustStore(KeyStore trustStore) throws
            NoSuchAlgorithmException, KeyStoreException {

        TrustManagerFactory tmf = TrustManagerFactory
                .getInstance(TrustManagerFactory.getDefaultAlgorithm());
        // Using null here initialises the TMF with the default trust store.
        tmf.init((KeyStore) null);

        // Get hold of the default trust manager
        X509TrustManager defaultTm = null;
        for (TrustManager tm : tmf.getTrustManagers()) {
            if (tm instanceof X509TrustManager) {
                defaultTm = (X509TrustManager) tm;
                break;
            }
        }

        // Get the certs from defaultTm and add them to our trustStore
        if (defaultTm != null) {
            X509Certificate[] trustedCerts = defaultTm.getAcceptedIssuers();
            Arrays.stream(trustedCerts)
                    .forEach(cert -> {
                        try {
                            trustStore.setCertificateEntry(RandomStringUtils.random(10, true, false),
                                    cert);
                        } catch (KeyStoreException e) {
                            log.error("Error while adding default trusted ca cert", e);
                        }
                    });
        }
    }

    public static KeyStore getKeyStore(String certPath, String keyPath) {
        KeyStore keyStore = null;
        try {
            Certificate cert =
                    TLSUtils.getCertificateFromFile(certPath);
            Key key = JWTUtils.getPrivateKey(keyPath);
            keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
            keyStore.load(null, null);
            keyStore.setKeyEntry("client-keys", key, null, new Certificate[]{cert});
        } catch (EnforcerException | CertificateException | IOException | KeyStoreException |
                 NoSuchAlgorithmException e) {
            log.error("Error occurred while configuring KeyStore", e);
        }
        return keyStore;
    }
}
