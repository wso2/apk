package org.wso2.apk.enforcer.analytics.publisher.retriever;

/**
 *Interface for moesif KeyRetriever.
 */
public interface MoesifKeyRetriever {

    /**
     * Returns Key for given organization and environment
     * @param organization
     * @param environment
     * @return APiKey.
     */
    public String getKey(String organization, String environment);
}
