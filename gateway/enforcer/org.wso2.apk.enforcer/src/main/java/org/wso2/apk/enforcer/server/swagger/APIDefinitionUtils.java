package org.wso2.apk.enforcer.server.swagger;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.grpc.netty.shaded.io.netty.handler.codec.http.HttpResponseStatus;
import org.wso2.apk.enforcer.models.ResponsePayload;

public class APIDefinitionUtils {
    public static ResponsePayload buildResponsePayload(Object dataModel, HttpResponseStatus status, boolean isError)
            throws JsonProcessingException {

        String jsonPayload;
        ObjectMapper objectMapper = new ObjectMapper();
        if (!(dataModel instanceof String)) {
            jsonPayload = objectMapper.writeValueAsString(dataModel);
        } else {
            jsonPayload = (String) dataModel;
        }
        ResponsePayload responsePayload = new ResponsePayload();
        responsePayload.setContent(jsonPayload);
        responsePayload.setError(isError);
        responsePayload.setStatus(status);
        return responsePayload;
    }
}
