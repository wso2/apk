package org.wso2.apk.runtime.model;


import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class API {
    private String name;
    private String version;
    private String type;
    private Set<URITemplate> uriTemplates = new HashSet<>();
    private String apiSecurity;
    private Set<Scope> scopes;
    private String graphQLSchema;
    private String swaggerDefinition;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public Set<URITemplate> getUriTemplates() {
        return uriTemplates;
    }

    public void setUriTemplates(Set<URITemplate> uriTemplates) {
        this.uriTemplates = uriTemplates;
    }

    public API(String name, String version, Set<URITemplate> uriTemplates) {
        this.name = name;
        this.version = version;
        this.uriTemplates = uriTemplates;
    }

    public API() {
    }

    public void setApiSecurity(String apiSecurity) {
        this.apiSecurity = apiSecurity;
    }

    public String getApiSecurity() {
        return apiSecurity;
    }

    public Set<Scope> getScopes() {
        return scopes;
    }

    public void setScopes(Set<Scope> scopes) {
        this.scopes = scopes;
    }

    public String getGraphQLSchema() {
        return graphQLSchema;
    }

    public void setGraphQLSchema(String graphQLSchema) {
        this.graphQLSchema = graphQLSchema;
    }

    public String getSwaggerDefinition() {
        return swaggerDefinition;
    }

    public void setSwaggerDefinition(String swaggerDefinition) {
        this.swaggerDefinition = swaggerDefinition;
    }
}
