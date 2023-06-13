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


import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public class URITemplate implements Serializable{

    private static final long serialVersionUID = 1L;

    private String uriTemplate;
    private String resourceURI;
    private String httpVerb;
    private boolean authEnabled = true;
    private List<String> scopes = new ArrayList<String>();
    private int id;
    private String endpoint;

    public String getHTTPVerb() {
        return httpVerb;
    }

    public void setHTTPVerb(String httpVerb) {
        this.httpVerb = httpVerb;
    }

    public boolean isAuthEnabled() {
        return authEnabled;
    }

    public void setAuthEnabled(boolean authEnabled) {
        this.authEnabled = authEnabled;
    }

    public String getResourceURI() {
        return resourceURI;
    }

    public void setResourceURI(String resourceURI) {
        this.resourceURI = resourceURI;
    }

    public String getUriTemplate() {
        return uriTemplate;
    }

    public void setUriTemplate(String template) {
        this.uriTemplate = template;
    }

    public String[] getScopes() {
        return scopes.toArray(new String[scopes.size()]);
    }


    public void setScopes(String scope){
        this.scopes.add(scope);
    }

    @Override
    public boolean equals(Object o) {

        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        URITemplate that = (URITemplate) o;
        return authEnabled == that.authEnabled && id == that.id && Objects.equals(uriTemplate, that.uriTemplate) && Objects.equals(resourceURI, that.resourceURI) && Objects.equals(httpVerb, that.httpVerb) && Objects.equals(scopes, that.scopes) && Objects.equals(endpoint, that.endpoint);
    }

    @Override
    public int hashCode() {

        return Objects.hash(uriTemplate, resourceURI, httpVerb, authEnabled, scopes, id, endpoint);
    }

    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public String[] retrieveAllScopes() {
        return this.scopes.toArray(new String[scopes.size()]);
    }

    public void addAllScopes(List<String> scopes) {

        this.scopes = scopes;
    }

    public String getEndpoint() {

        return endpoint;
    }

    public void setEndpoint(String endpoint) {

        this.endpoint = endpoint;
    }
}
