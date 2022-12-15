/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
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
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.runtime.model;


import java.util.*;

/**
 * Intermediate model used to store data required for swagger processing
 */
public class SwaggerData {
    /**
     * Maps to Swagger PathItem/Operation
     */
    public static class Resource {
        private String path;
        private String verb;
        private String authType;
        private String policy;
        private Scope scope;
        private List<Scope> scopes = new ArrayList<>();
        private String amznResourceName;
        private int amznResourceTimeout;

        public String getPath() {
            return path;
        }

        public void setPath(String path) {
            this.path = path;
        }

        public String getVerb() {
            return verb;
        }

        public void setVerb(String verb) {
            this.verb = verb;
        }

        public String getAuthType() {
            return authType;
        }

        public void setAuthType(String authType) {
            this.authType = authType;
        }

        public String getPolicy() {
            return policy;
        }

        public void setPolicy(String policy) {
            this.policy = policy;
        }

        public Scope getScope() {
            return scope;
        }

        public void setScope(Scope scope) {
            this.scope = scope;
        }

        public String getAmznResourceName() {
            return amznResourceName;
        }

        public void setAmznResourceName(String amznResourceName) {
            this.amznResourceName = amznResourceName;
        }

        public int getAmznResourceTimeout() {
            return amznResourceTimeout;
        }

        public void setAmznResourceTimeout(int amznResourceTimeout) {
            this.amznResourceTimeout = amznResourceTimeout;
        }

        public List<Scope> getScopes() {

            return scopes;
        }

        public void setScopes(List<Scope> scopes) {

            this.scopes = scopes;
        }

    }

    private final String title;
    private String description;
    private final String version;
    private String contactName;
    private String contactEmail;
    private final String transportType;
    private String security;
    private String apiLevelPolicy;
    private final Set<Resource> resources = new LinkedHashSet<>();
    private final Set<Scope> scopes = new HashSet<>();

    public SwaggerData(API api) {
        title = api.getName();
        version = api.getVersion();


        Set<URITemplate> uriTemplates = api.getUriTemplates();

        for (URITemplate uriTemplate : uriTemplates) {
            Resource resource = new Resource();
            resource.path = uriTemplate.getUriTemplate();
            resource.verb = uriTemplate.getHTTPVerb();
            resource.authType = uriTemplate.getAuthType();
            resource.policy = uriTemplate.getThrottlingTier();
            resource.scope = uriTemplate.getScope();
            resource.scopes = uriTemplate.retrieveAllScopes();
            resource.amznResourceName = uriTemplate.getAmznResourceName();
            resource.amznResourceTimeout = uriTemplate.getAmznResourceTimeout();
            resources.add(resource);
        }

        transportType = api.getType();
    }

    public Set<Resource> getResources() {
        return resources;
    }

    public Set<Scope> getScopes() {
        return scopes;
    }

    public String getTitle() {
        return title;
    }

    public String getDescription() {
        return description;
    }

    public String getVersion() {
        return version;
    }

    public String getContactName() {
        return contactName;
    }

    public String getContactEmail() {
        return contactEmail;
    }

    public String getTransportType() {
        return transportType;
    }

    public String getSecurity() {
        return security;
    }

    public String getApiLevelPolicy() {
        return apiLevelPolicy;
    }
}
