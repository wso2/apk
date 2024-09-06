package org.wso2.apk.config.model.proto;

import java.util.List;

public class Service {
    String serviceName;
    List<String> methods;

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public void setMethods(List<String> methods) {
        this.methods = methods;
    }

    public String getServiceName() {
        return serviceName;
    }

    public List<String> getMethods() {
        return methods;
    }

    public Service() { }

    public Service(String serviceName, List<String> methods) {
        this.serviceName = serviceName;
        this.methods = methods;
    }

    @Override
    public String toString() {
        return "Service{" + "serviceName='" + serviceName + '\'' + ", methods=" + methods + '}';
    }
}