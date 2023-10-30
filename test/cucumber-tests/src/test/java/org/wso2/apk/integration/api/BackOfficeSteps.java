package org.wso2.apk.integration.api;

import java.util.HashMap;
import java.util.Map;
import io.cucumber.java.en.When;

import org.apache.http.HttpResponse;
import org.wso2.apk.integration.utils.Constants;
import org.wso2.apk.integration.utils.Utils;
import org.wso2.apk.integration.utils.clients.SimpleHTTPClient;

public class BackOfficeSteps {

    private final SharedContext sharedContext;

    public BackOfficeSteps(SharedContext sharedContext) {

        this.sharedContext = sharedContext;
    }

    @When("I make the GET APIs call to the backoffice")
    public void make_a_deployment_request() throws Exception {
        Map<String, String> headers = new HashMap<>();
        headers.put(Constants.REQUEST_HEADERS.AUTHORIZATION, "Bearer " + sharedContext.getAccessToken());
        headers.put(Constants.REQUEST_HEADERS.HOST, Constants.DEFAULT_API_HOST);

        HttpResponse response = sharedContext.getHttpClient().doGet(Utils.getBackOfficeAPIURL(),
                headers);

        sharedContext.setResponse(response);
        sharedContext.setResponseBody(SimpleHTTPClient.responseEntityBodyToString(sharedContext.getResponse()));
    }
}
