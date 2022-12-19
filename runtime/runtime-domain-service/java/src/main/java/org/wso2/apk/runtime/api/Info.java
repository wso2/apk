package org.wso2.apk.runtime.api;

import java.util.List;

public class Info {
    private String openAPIVersion;
    private String name;
    private String version;
    private String context;
    private String description;
    private List<String> endpoints;

    public String getOpenAPIVersion() {
        return openAPIVersion;
    }

    public void setOpenAPIVersion(String openAPIVersion) {
        this.openAPIVersion = openAPIVersion;
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

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getContext() {
        return context;
    }

    public void setContext(String context) {
        this.context = context;
    }

    public List<String> getEndpoints() { return endpoints; }

    public void setEndpoints(List<String> endpoints) { this.endpoints = endpoints; }
}
