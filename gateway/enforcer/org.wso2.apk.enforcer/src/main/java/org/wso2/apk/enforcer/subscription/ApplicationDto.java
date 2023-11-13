package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Entity for keeping Application related information. Represents an Application in APK.
 */
public class ApplicationDto implements Serializable {

    private static final long serialVersionUID = 1L;

    private String uuid;
    private String name;
    private String owner;
    private Map<String, String> attributes = new ConcurrentHashMap<>();

    private String organizationId;

    public String getUuid() {

        return uuid;
    }

    public void setUuid(String uuid) {

        this.uuid = uuid;
    }

    public String getName() {

        return name;
    }

    public void setName(String name) {

        this.name = name;
    }

    public String getOwner() {

        return owner;
    }

    public void setOwner(String owner) {

        this.owner = owner;
    }

    public Map<String, String> getAttributes() {

        return attributes;
    }

    public void setAttributes(Map<String, String> attributes) {

        this.attributes = attributes;
    }

    public String getOrganizationId() {

        return organizationId;
    }

    public void setOrganizationId(String organizationId) {

        this.organizationId = organizationId;
    }
}
