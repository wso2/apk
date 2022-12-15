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


import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

public class URITemplate implements Serializable{

    private static final long serialVersionUID = 1L;

    private String uriTemplate;
    private String resourceURI;
    private String resourceSandboxURI;
    private String httpVerb;
    private String authType;
    private String throttlingTier;
    private Scope scope;
    private List<Scope> scopes = new ArrayList<Scope>();
    private int id;
    private String amznResourceName;
    private int amznResourceTimeout;
    private List<OperationPolicy> operationPolicies = new ArrayList<>();


    public String getThrottlingTier() {
        return throttlingTier;
    }

    public void setThrottlingTier(String throttlingTier) {
        this.throttlingTier = throttlingTier;
    }

    public String getHTTPVerb() {
        return httpVerb;
    }

    public void setHTTPVerb(String httpVerb) {
        this.httpVerb = httpVerb;
    }

    public String getAuthType() {
        return authType;
    }

    public void setAuthType(String authType) {
        this.authType = authType;

    }

    public String getResourceURI() {
        return resourceURI;
    }

    public void setResourceURI(String resourceURI) {
        this.resourceURI = resourceURI;
    }

    public boolean isResourceURIExist(){
        return this.resourceURI != null;
    }

    public String getResourceSandboxURI() {
        return resourceSandboxURI;
    }

    public void setResourceSandboxURI(String resourceSandboxURI) {
        this.resourceSandboxURI = resourceSandboxURI;
    }

    public boolean isResourceSandboxURIExist(){
        return this.resourceSandboxURI != null;
    }

    public String getUriTemplate() {
        return uriTemplate;
    }

    public void setUriTemplate(String template) {
        this.uriTemplate = template;
    }

    public Scope getScope() {
        return scope;
    }
    public List<Scope> getScopes() {
        return scopes;
    }

    public void setScope(Scope scope) {
        this.scope = scope;
    }

    public void setScopes(Scope scope){
        this.scopes.add(scope);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }

        URITemplate that = (URITemplate) o;

        if (!uriTemplate.equals(that.uriTemplate)) {
            return false;
        }
        if (resourceURI != null ? !resourceURI.equals(that.resourceURI) : that.resourceURI != null) {
            return false;
        }
        if (resourceSandboxURI != null ? !resourceSandboxURI.equals(that.resourceSandboxURI) : that
                .resourceSandboxURI != null) {
            return false;
        }
        if (!httpVerb.equals(that.httpVerb)) {
            return false;
        }
        if (!authType.equals(that.authType)) {
            return false;
        }

        if (!throttlingTier.equals(that.throttlingTier)) {
            return false;
        }
        if (scope != null ? !scope.equals(that.scope) : that.scope != null) {
            return false;
        }
        return scopes != null ? scopes.equals(that.scopes) : that.scopes == null;
    }

    @Override
    public int hashCode() {
        int result = uriTemplate.hashCode();
        result = 31 * result + (resourceURI != null ? resourceURI.hashCode() : 0);
        result = 31 * result + (resourceSandboxURI != null ? resourceSandboxURI.hashCode() : 0);
        result = 31 * result + (httpVerb != null ? httpVerb.hashCode() : 0);
        result = 31 * result + (authType != null ? authType.hashCode() : 0);
        result = 31 * result + (throttlingTier != null ? throttlingTier.hashCode() : 0);
        result = 31 * result + (scope != null ? scope.hashCode() : 0);
        result = 31 * result + (scopes != null ? scopes.hashCode() : 0);
        return result;
    }

    public int getId() {
        return id;
    }

    public void setId(int id) {
        this.id = id;
    }

    public List<Scope> retrieveAllScopes() {
        return this.scopes;
    }

    public void addAllScopes(List<Scope> scopes) {

        this.scopes = scopes;
    }

    public void setAmznResourceName(String amznResourceName) {
        this.amznResourceName = amznResourceName;
    }

    public String getAmznResourceName() {
        return amznResourceName;
    }

    public void setAmznResourceTimeout(int amznResourceTimeout) {
        this.amznResourceTimeout = amznResourceTimeout;
    }

    public int getAmznResourceTimeout() {
        return amznResourceTimeout;
    }

    public void setOperationPolicies(List<OperationPolicy> operationPolicies) {
        this.operationPolicies = operationPolicies;
    }

    public List<OperationPolicy> getOperationPolicies() {
        return operationPolicies;
    }

    public void addOperationPolicy(OperationPolicy policy) {
        operationPolicies.add(policy);
    }
}
