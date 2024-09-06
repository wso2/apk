package org.wso2.apk.config.model;

public class API {
    private String name;
    private String basePath;
    private String version;
    private String type;
    private String endpoint;
    private URITemplate[] uriTemplates = new URITemplate[0];
    private String apiSecurity;
    private String[] scopes;
    private String graphQLSchema;
    private String protoDefinition;
    private String swaggerDefinition;
    private String environment;

    public API(String name, String version, URITemplate[] uriTemplates) {
        this.name = name;
        this.version = version;
        this.uriTemplates = uriTemplates;
    }

    public API() {
    }

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

    public URITemplate[] getUriTemplates() {
        return uriTemplates;
    }

    public void setUriTemplates(URITemplate[] uriTemplates) {
        this.uriTemplates = uriTemplates;
    }

    public String getApiSecurity() {
        return apiSecurity;
    }

    public void setApiSecurity(String apiSecurity) {
        this.apiSecurity = apiSecurity;
    }

    public String[] getScopes() {
        return scopes;
    }

    public void setScopes(String[] scopes) {
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

    public String getProtoDefinition() {
        return protoDefinition;
    }

    public void setProtoDefinition(String protoDefinition) {
        this.protoDefinition = protoDefinition;
    }

    public String getBasePath() {

        return basePath;
    }

    public void setBasePath(String basePath) {

        this.basePath = basePath;
    }

    public String getEndpoint() {

        return endpoint;
    }

    public void setEndpoint(String endpoint) {

        this.endpoint = endpoint;
    }

    public String getEnvironment() {
        return environment;
    }

    public void setEnvironment(String environment) {
        this.environment = environment;
    }
}
