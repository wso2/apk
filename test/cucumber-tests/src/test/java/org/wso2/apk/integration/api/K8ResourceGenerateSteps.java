package org.wso2.apk.integration.api;

import com.google.common.io.Resources;
import io.cucumber.java.en.Then;
import io.cucumber.java.en.When;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.entity.mime.HttpMultipartMode;
import org.apache.http.entity.mime.MultipartEntityBuilder;
import org.apache.http.entity.mime.content.FileBody;
import org.testng.Assert;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.util.zip.ZipEntry;
import java.util.zip.ZipInputStream;

import org.apache.commons.io.IOUtils;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;

public class K8ResourceGenerateSteps {
    
    private final SharedContext sharedContext;
    private static final Log logger = LogFactory.getLog(BaseSteps.class);
    private File definitionFile;
    private File apkConfFile;

    public K8ResourceGenerateSteps(SharedContext sharedContext) {
        this.sharedContext = sharedContext;
    }

    @When("I use the definition file {string}")
    public void i_use_the_definition_file(String definitionFilePath) {

        URL url = Resources.getResource(definitionFilePath);
        definitionFile = new File(url.getPath());
    }

    @When("I use the apk conf file {string} in resources")
    public void i_use_the_apkconf_file_in_resources(String confFilePath) {

        URL url = Resources.getResource(confFilePath);
        apkConfFile = new File(url.getPath());
    }

    @When("I generate and apply the K8Artifacts belongs to that API")
    public void generate_the_k8artifacts_set() throws Exception {

        // Create a MultipartEntityBuilder to build the request entity
        MultipartEntityBuilder builder = MultipartEntityBuilder.create()
                .setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
                .addPart("apkConfiguration", new FileBody(apkConfFile))
                .addPart("definitionFile", new FileBody(definitionFile));

        HttpEntity multipartEntity = builder.build();
        HttpResponse httpResponse = sharedContext.getHttpClient().doPostWithMultipart(Utils.getK8ResourceGeneratorURL(),
            multipartEntity);
        
        sharedContext.setResponse(httpResponse);

        // Process the HTTP response
        if (httpResponse.getStatusLine().getStatusCode() == 200) {
            HttpEntity entity = httpResponse.getEntity();
            if (entity != null) {
                try (InputStream inputStream = entity.getContent()) {
                    // Create a temporary file and store the zip content
                    File tempFile = File.createTempFile("k8s_artifacts", ".zip");
                    try (FileOutputStream outputStream = new FileOutputStream(tempFile)) {
                        inputStream.transferTo(outputStream);
                        logger.info("Data transformation complete");
                    }
    
                    // Process the zip file
                    try (ZipInputStream zipInputStream = new ZipInputStream(new FileInputStream(tempFile))) {
                        ZipEntry entry;
                        while ((entry = zipInputStream.getNextEntry()) != null) {
                            if (!entry.isDirectory()) {
                                if (entry.getName().endsWith(".yaml") || entry.getName().endsWith(".yml")) {
                                    // Create a temporary file and store the YAML content into it
                                    File yamlFile = File.createTempFile("k8s_yaml", ".yaml");
                                    try (FileOutputStream yamlOutputStream = new FileOutputStream(yamlFile)) {
                                        byte[] buffer = new byte[1024];
                                        int bytesRead;
                                        while ((bytesRead = zipInputStream.read(buffer)) != -1) {
                                            yamlOutputStream.write(buffer, 0, bytesRead);
                                        }
                                    }

                                    String command = "kubectl apply -f " + yamlFile.getAbsolutePath() + " -n apk";
                                    Process process = Runtime.getRuntime().exec(command);
                                    int exitCode = process.waitFor();
                                    if (exitCode != 0) {
                                        logger.error("Error applying YAML file: " + entry.getName());
                                    } else {
                                        logger.info("File: " + entry.getName() + " applied successfully.");
                                    }
                                }
                            }
                        }
                    } catch (IOException e) {
                        e.printStackTrace();
                    }
                }
            }
        } else {
            logger.info("Failed to generate K8s artifacts. HTTP status code: " + httpResponse.getStatusLine().getStatusCode());
        }
    }

}
