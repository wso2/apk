package org.wso2.apk.enforcer.subscription;

import java.io.Serializable;

public class SubscribedAPIDto implements Serializable {

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

    private String name;
    private String version;
}
