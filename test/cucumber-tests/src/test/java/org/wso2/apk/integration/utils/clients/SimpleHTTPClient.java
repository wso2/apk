/*
 * Copyright (c) 2023, WSO2 LLC (http://www.wso2.com).
 *
 * WSO2 LLC licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.wso2.apk.integration.utils.clients;

import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.apache.http.HttpEntity;
import org.apache.http.HttpEntityEnclosingRequest;
import org.apache.http.HttpHeaders;
import javax.net.ssl.TrustManager;
import org.apache.http.HttpResponse;
import org.apache.http.client.HttpClient;
import org.apache.http.client.config.RequestConfig;
import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.client.methods.HttpDelete;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpHead;
import org.apache.http.client.methods.HttpOptions;
import org.apache.http.client.methods.HttpPatch;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.client.methods.HttpPut;
import org.apache.http.client.methods.HttpUriRequest;
import org.apache.http.conn.HttpClientConnectionManager;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustAllStrategy;
import org.apache.http.entity.ContentProducer;
import org.apache.http.entity.ContentType;
import org.apache.http.entity.EntityTemplate;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;
import org.apache.http.ssl.SSLContexts;
import org.wso2.apk.integration.utils.MultipartFilePart;
import org.wso2.apk.integration.utils.exceptions.TimeoutException;

import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.util.HashMap;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;
import javax.net.ssl.SSLContext;
import javax.net.ssl.X509TrustManager;

public class SimpleHTTPClient {

    protected Log log = LogFactory.getLog(getClass());
    private CloseableHttpClient client;
    private HttpUriRequest lastRequest;
    private static final int EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS = 15;

    public SimpleHTTPClient() throws NoSuchAlgorithmException, KeyStoreException, KeyManagementException {
        String httpClientSetup = System.getProperty("http.client.setup", "apk");
        log.info(httpClientSetup);

        if ("apk".equals(httpClientSetup)) {
            final SSLContext sslcontext = SSLContexts.custom()
            .loadTrustMaterial(null, new TrustAllStrategy())
            .build();

            final SSLConnectionSocketFactory csf = new SSLConnectionSocketFactory(sslcontext);
            RequestConfig requestConfig = RequestConfig.custom()
                    .setRedirectsEnabled(false) // Disable redirects
                    .build();
            this.client = HttpClients.custom()
                    .setDefaultRequestConfig(requestConfig)
                    .setSSLSocketFactory(csf)
                    .evictExpiredConnections()
                    .setMaxConnPerRoute(100)
                    .setMaxConnTotal(1000)
                    .build();
            this.lastRequest = null;
        }
        else if ("apim-apk".equals(httpClientSetup)) {
            // Create SSL context that trusts all certificates
            SSLContext sslContext = createAcceptAllSSLContext();

            // Create a socket factory with custom SSL context and hostname verifier that accepts all hostnames
            SSLConnectionSocketFactory sslSocketFactory = new SSLConnectionSocketFactory(sslContext,
                    NoopHostnameVerifier.INSTANCE);

            // Create HttpClient with custom SSL socket factory
            this.client = HttpClientBuilder.create().setSSLSocketFactory(sslSocketFactory).build();
            this.lastRequest = null;
            }
    }

    private SSLContext createAcceptAllSSLContext() throws NoSuchAlgorithmException, KeyManagementException {
        // Create a TrustManager that trusts all certificates
        TrustManager[] trustAllCerts = new TrustManager[]{
                new X509TrustManager() {
                    public java.security.cert.X509Certificate[] getAcceptedIssuers() {
                        return null;
                    }

                    public void checkClientTrusted(X509Certificate[] certs, String authType) {
                    }

                    public void checkServerTrusted(X509Certificate[] certs, String authType) {
                    }
                }
        };

        // Create SSL context with the TrustManager that trusts all certificates
        SSLContext sslContext = SSLContext.getInstance("TLS");
        sslContext.init(null, trustAllCerts, new java.security.SecureRandom());
        return sslContext;
    }

    /**
     * Function to extract response body as a string
     *
     * @param response org.apache.http.HttpResponse object containing response entity body
     * @return returns the response entity body as a string
     * @throws IOException
     */
    public static String responseEntityBodyToString(HttpResponse response) throws IOException {

        if (response != null && response.getEntity() != null) {
            try (InputStream inputStreamContent = response.getEntity().getContent()) {
                return IOUtils.toString(inputStreamContent);
            }
        }
        return null;
    }

    /**
     * Send a HTTP GET request to the specified URL
     *
     * @param url     Target endpoint URL
     * @param headers Any HTTP headers that should be added to the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doGet(String url, Map<String, String> headers) throws IOException {

        HttpUriRequest request = new HttpGet(url);
        setHeaders(headers, request);
        this.lastRequest = request;
        return client.execute(request);
    }

    /**
     * Send a HTTP POST request to the specified URL
     *
     * @param url         Target endpoint URL
     * @param headers     Any HTTP headers that should be added to the request
     * @param payload     Content payload that should be sent
     * @param contentType Content-type of the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doPost(String url, final Map<String, String> headers, final String payload, String contentType)
            throws IOException {

        HttpUriRequest request = new HttpPost(url);
        setHeaders(headers, request);
        HttpEntityEnclosingRequest entityEncReq = (HttpEntityEnclosingRequest) request;
        final boolean zip = headers != null && "gzip".equals(headers.get(HttpHeaders.CONTENT_ENCODING));

        EntityTemplate ent = new EntityTemplate(new ContentProducer() {
            public void writeTo(OutputStream outputStream) throws IOException {
                OutputStream out = outputStream;
                if (zip) {
                    out = new GZIPOutputStream(outputStream);
                }
                out.write(payload.getBytes());
                out.flush();
                out.close();
            }
        });
        if (contentType != null) {
            ent.setContentType(contentType);
        } else {
            ent.setContentType(MediaType.JSON.getValue());
        }
        if (zip) {
            ent.setContentEncoding("gzip");
        }
        entityEncReq.setEntity(ent);
        this.lastRequest = request;
        log.info("Request: " + request);
        return client.execute(request);
    }

    /**
     * Send a HTTP POST with multipart request to the specified URL
     *
     * @param url Target endpoint URL
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doPostWithMultipart(String url, HttpEntity httpEntity)
            throws IOException {

        return doPostWithMultipart(url, httpEntity, new HashMap<>());
    }

    public HttpResponse doPostWithMultipart(String url, HttpEntity httpEntity, Map<String, String> header)
            throws IOException {

        HttpPost request = new HttpPost(url);
        for (String headerKey : header.keySet()) {
            request.addHeader(headerKey, header.get(headerKey));
        }
        request.setEntity(httpEntity);
        this.lastRequest = request;
        return client.execute(request);
    }

    public HttpResponse doPostWithMultipart(String url, List<MultipartFilePart> fileParts, Map<String, String> header)
            throws IOException {

        MultipartEntityBuilder entitybuilder = MultipartEntityBuilder.create();
        entitybuilder.setMode(HttpMultipartMode.BROWSER_COMPATIBLE);
        for (MultipartFilePart filePart : fileParts) {
            entitybuilder.addPart(filePart.getName(), new FileBody(filePart.getFile()));
        }
        HttpPost request = new HttpPost(url);
        for (String headerKey : header.keySet()) {
            request.addHeader(headerKey, header.get(headerKey));
        }
        HttpEntity mutiPartHttpEntity = entitybuilder.build();
        request.setEntity(mutiPartHttpEntity);
        this.lastRequest = request;
        return client.execute(request);
    }

    public HttpResponse doPutWithMultipart(String url, File file, Map<String, String> header)
            throws IOException {

        MultipartEntityBuilder entitybuilder = MultipartEntityBuilder.create();
        entitybuilder.setMode(HttpMultipartMode.BROWSER_COMPATIBLE);
        entitybuilder.addBinaryBody("file", file, ContentType.APPLICATION_OCTET_STREAM, file.getName());
        HttpPut request = new HttpPut(url);
        for (String headerKey : header.keySet()) {
            request.addHeader(headerKey, header.get(headerKey));
        }
        HttpEntity mutiPartHttpEntity = entitybuilder.build();
        request.setEntity(mutiPartHttpEntity);
        this.lastRequest = request;
        return client.execute(request);
    }

    /**
     * Extracts the payload from a HTTP response. For a given HttpResponse object, this
     * method can be called only once.
     *
     * @param response HttpResponse instance to be extracted
     * @return Content payload
     * @throws IOException If an error occurs while reading from the response
     */
    public String getResponsePayload(HttpResponse response) throws IOException {

        if (response.getEntity() != null) {
            InputStream in = response.getEntity().getContent();
            int length;
            byte[] tmp = new byte[2048];
            StringBuilder buffer = new StringBuilder();
            while ((length = in.read(tmp)) != -1) {
                buffer.append(new String(tmp, 0, length));
            }
            return buffer.toString();
        }
        return null;
    }

    /**
     * Send a HTTP PATCH request to the specified URL
     *
     * @param url         Target endpoint URL
     * @param headers     Any HTTP headers that should be added to the request
     * @param payload     Content payload that should be sent
     * @param contentType Content-type of the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doPatch(String url, final Map<String, String> headers, final String payload, String contentType)
            throws IOException {

        HttpUriRequest request = new HttpPatch(url);
        setHeaders(headers, request);
        HttpEntityEnclosingRequest entityEncReq = (HttpEntityEnclosingRequest) request;
        final boolean zip = headers != null && "gzip".equals(headers.get(HttpHeaders.CONTENT_ENCODING));

        EntityTemplate ent = new EntityTemplate(new ContentProducer() {
            public void writeTo(OutputStream outputStream) throws IOException {

                OutputStream out = outputStream;
                if (zip) {
                    out = new GZIPOutputStream(outputStream);
                }
                out.write(payload.getBytes());
                out.flush();
                out.close();
            }
        });
        ent.setContentType(contentType);
        if (zip) {
            ent.setContentEncoding("gzip");
        }
        entityEncReq.setEntity(ent);
        return client.execute(request);
    }

    /**
     * Send a HTTP OPTIONS request to the specified URL
     *
     * @param url         Target endpoint URL
     * @param headers     Any HTTP headers that should be added to the request
     * @param payload     Content payload that should be sent
     * @param contentType Content-type of the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doOptions(String url, final Map<String, String> headers, final String payload,
                                  String contentType) throws IOException {

        HttpUriRequest request = new HttpOptions(url);
        setHeaders(headers, request);
        if (payload != null) {
            HttpEntityEnclosingRequest entityEncReq = (HttpEntityEnclosingRequest) request;
            final boolean zip = headers != null && "gzip".equals(headers.get(HttpHeaders.CONTENT_ENCODING));

            EntityTemplate ent = new EntityTemplate(new ContentProducer() {
                public void writeTo(OutputStream outputStream) throws IOException {

                    OutputStream out = outputStream;
                    if (zip) {
                        out = new GZIPOutputStream(outputStream);
                    }
                    out.write(payload.getBytes());
                    out.flush();
                    out.close();
                }
            });
            ent.setContentType(contentType);
            if (zip) {
                ent.setContentEncoding("gzip");
            }
            entityEncReq.setEntity(ent);
        }
        return client.execute(request);
    }

    /**
     * Send a HTTP Head request to the specified URL
     *
     * @param url     Target endpoint URL
     * @param headers Any HTTP headers that should be added to the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doHead(String url, final Map<String, String> headers) throws IOException {

        HttpUriRequest request = new HttpHead(url);
        setHeaders(headers, request);
        return client.execute(request);
    }

    /**
     * Send a HTTP DELETE request to the specified URL
     *
     * @param url     Target endpoint URL
     * @param headers Any HTTP headers that should be added to the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doDelete(String url, final Map<String, String> headers) throws IOException {

        HttpUriRequest request = new HttpDelete(url);
        setHeaders(headers, request);
        this.lastRequest = lastRequest;
        return client.execute(request);
    }

    /**
     * Send a HTTP PUT request to the specified URL
     *
     * @param url         Target endpoint URL
     * @param headers     Any HTTP headers that should be added to the request
     * @param payload     Content payload that should be sent
     * @param contentType Content-type of the request
     * @return Returned HTTP response
     * @throws IOException If an error occurs while making the invocation
     */
    public HttpResponse doPut(String url, final Map<String, String> headers, final String payload, String contentType)
            throws IOException {

        HttpUriRequest request = new HttpPut(url);
        setHeaders(headers, request);
        HttpEntityEnclosingRequest entityEncReq = (HttpEntityEnclosingRequest) request;
        final boolean zip = headers != null && "gzip".equals(headers.get(HttpHeaders.CONTENT_ENCODING));

        EntityTemplate ent = new EntityTemplate(new ContentProducer() {
            public void writeTo(OutputStream outputStream) throws IOException {

                OutputStream out = outputStream;
                if (zip) {
                    out = new GZIPOutputStream(outputStream);
                }
                out.write(payload.getBytes());
                out.flush();
                out.close();
            }
        });
        ent.setContentType(contentType);
        if (zip) {
            ent.setContentEncoding("gzip");
        }
        entityEncReq.setEntity(ent);
        this.lastRequest = lastRequest;
        return client.execute(request);
    }

    private void setHeaders(Map<String, String> headers, HttpUriRequest request) {

        if (headers != null && headers.size() > 0) {
            for (Map.Entry<String, String> header : headers.entrySet()) {
                request.setHeader(header.getKey(), header.getValue());
            }
        }
    }

    public HttpResponse executeLastRequestForEventualConsistentResponse(int successResponseCode,
                                                                        List<Integer> nonAcceptableCodes) throws IOException, InterruptedException {

        int counter = 1;
        int responseCode = -1;
        String lastResponseBody = null;
        while (counter < EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS) {
            counter++;
            Thread.sleep(1000);
            HttpResponse httpResponse = getClient().execute(lastRequest);
            responseCode = httpResponse.getStatusLine().getStatusCode();
            if (responseCode == successResponseCode || nonAcceptableCodes.contains(responseCode)) {
                return httpResponse;
            } else {
                if (counter == EVENTUAL_SUCCESS_RESPONSE_TIMEOUT_IN_SECONDS) {
                    lastResponseBody = responseEntityBodyToString(httpResponse);
                }
                ((CloseableHttpResponse) httpResponse).close();
            }
        }
        throw new TimeoutException("Could not receive expected response within time. Last received code: "
                + responseCode + ", last response body: " + lastResponseBody);
    }

    private HttpClient getClient() {

        final SSLContext sslcontext;
        try {
            sslcontext = SSLContexts.custom()
                    .loadTrustMaterial(null, new TrustAllStrategy())
                    .build();
        } catch (NoSuchAlgorithmException | KeyManagementException | KeyStoreException e) {
            throw new RuntimeException(e);
        }
        final SSLConnectionSocketFactory csf = new SSLConnectionSocketFactory(sslcontext);

        return HttpClients.custom()
                .setSSLSocketFactory(csf)
                .evictExpiredConnections()
                .build();
    }
}

enum MediaType {
    JSON("application/json"),
    XML("application/xml"),
    FORM("application/x-www-form-urlencoded");
    // Add more Content-Type values as needed

    private final String value;

    MediaType(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }
}
