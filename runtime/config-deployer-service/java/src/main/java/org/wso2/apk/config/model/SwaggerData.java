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
package org.wso2.apk.config.model;

import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.Set;

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
        private boolean authType;
        private String policy;
        private String[] scopes;

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

        public boolean isAuthType() {

            return authType;
        }

        public void setAuthType(boolean authType) {

            this.authType = authType;
        }

        public String getPolicy() {

            return policy;
        }

        public void setPolicy(String policy) {

            this.policy = policy;
        }



        public String[] getScopes() {

            return scopes;
        }

        public void setScopes(String[] scopes) {

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
    private final Set<String> scopes = new HashSet<>();

    public SwaggerData(API api) {

        title = api.getName();
        version = api.getVersion();

        URITemplate[] uriTemplates = api.getUriTemplates();

        for (URITemplate uriTemplate : uriTemplates) {
            Resource resource = new Resource();
            resource.path = uriTemplate.getUriTemplate();
            resource.verb = uriTemplate.getHTTPVerb();
            resource.authType = uriTemplate.isAuthEnabled();
            resource.scopes = uriTemplate.retrieveAllScopes();
            resources.add(resource);
        }

        transportType = api.getType();
    }

    public Set<Resource> getResources() {

        return resources;
    }

    public Set<String> getScopes() {

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
