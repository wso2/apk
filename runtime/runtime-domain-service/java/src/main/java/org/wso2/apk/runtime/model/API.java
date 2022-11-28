package org.wso2.apk.runtime.model;

import org.wso2.apk.apimgt.api.model.URITemplate;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

public class API {
    private String name;
    private String version;
    private String type;
    private Set<URITemplate> uriTemplates = new HashSet<>();

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
}
