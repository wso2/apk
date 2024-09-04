package org.wso2.apk.config.model.proto;

import java.util.List;

public class Service {
    String serviceName;
    String packageName;
    List<String> methods;

    public String getServiceName() {
        return serviceName;
    }

    public String getPackageName() {
        return packageName;
    }

    public List<String> getMethods() {
        return methods;
    }

    public Service(String serviceName, List<String> methods, String packageName) {
        this.serviceName = serviceName;
        this.methods = methods;
        this.packageName = packageName;
    }

    @Override
    public String toString() {
        return "Service{" + "serviceName='" + serviceName + '\'' + ", packageName='" + packageName + '\'' + ", methods=" + methods + '}';
    }
}