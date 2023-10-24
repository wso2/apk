package org.wso2.apk.enforcer.analytics.publisher.retriever;

import org.wso2.apk.enforcer.analytics.publisher.util.Constants;

import java.util.Map;

/**
 * Factory Class to get Moesif Key Retriever.
 */
public class MoesifKeyRetrieverFactory {

    private MoesifKeyRetrieverFactory() {

    }

    public static MoesifKeyRetriever getMoesifKeyRetriever(Map<String, String> properties) {

        String clientType = properties.get(Constants.MOESIF_KEY_RETRIEVER_CLIENT_TYPE);

        if (Constants.MOESIF_KEY_RETRIEVER_CHOREO_CLIENT.equals(clientType)) {
            return MoesifKeyRetrieverChoreoClient.getInstance(properties);
        }
        return new ConfigBaseMoesifKeyRetriever(properties);
    }
}
