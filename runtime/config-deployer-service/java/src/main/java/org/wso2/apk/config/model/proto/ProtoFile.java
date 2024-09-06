package org.wso2.apk.config.model.proto;

import java.util.ArrayList;
import java.util.List;

public class ProtoFile {

    String apiName;
     String packageName;
     String basePath;
     String version;
     List<Service> services;

    public String getApiName() {
        if (apiName == null) {
            return "";
        }
        return apiName;
    }

    public void setApiName(String apiName) {
        this.apiName = apiName;
    }

    public List<Service> getServices() {
        return services;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public void setBasePath(String basePath) {
        this.basePath = basePath;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public void setServices(List<Service> services) {
        if (this.services == null){
            this.services = new ArrayList<>();
        }
        this.services.addAll(services);
    }

    public void addService(Service service) {
        if (this.services == null){
            this.services = new ArrayList<>();
        }
        this.services.add(service);
    }

    public String getBasePath() {
        return basePath;
    }

    public String getVersion() {
        return version;
    }

    public ProtoFile(String packageName, String basePath, String version, List<Service> services) {
        this.packageName = packageName;
        this.basePath = basePath;
        this.version = version;
        this.services = services;
    }

    public ProtoFile() {
        this.packageName = "packageName";
        this.basePath = "basePath";
        this.version = "version";
        this.services = new ArrayList<>();
    }

    @Override
    public String toString() {
        return "ProtoFile{" +
                "packageName='" + packageName + '\'' +
                ", basePath='" + basePath + '\'' +
                ", version='" + version + '\'' +
                ", services=" + services +
                '}';
    }

}