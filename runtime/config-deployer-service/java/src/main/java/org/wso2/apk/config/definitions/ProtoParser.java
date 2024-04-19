package org.wso2.apk.config.definitions;

import java.util.List;
import java.util.ArrayList;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ProtoParser {

    public ProtoFile protoFile;
    public ProtoParser(String content) {
        protoFile = parseProtoContent(content);
    }

    public String getPackageString(String content) {
        Pattern packagePattern = Pattern.compile("package\\s+([\\w\\.]+);");
        Matcher packageMatcher = packagePattern.matcher(content);
        if (packageMatcher.find()) {
            return packageMatcher.group(1);
        }
        return null;
    }
    public String getVersion(String packageString) {
        Pattern versionPattern = Pattern.compile("v\\d+(\\.\\d+)*");
        Matcher versionMatcher = versionPattern.matcher(packageString);
        if (versionMatcher.find()) {
            return versionMatcher.group(0);
        }
        System.out.println("No version found");
        return null;
    }
    public String getPackageName(String packageString){
        Pattern namePattern = Pattern.compile("v\\d+(\\.\\d+)*\\.(\\w+)");
        Matcher nameMatcher = namePattern.matcher(packageString);
        if (nameMatcher.find()) {
            return nameMatcher.group(2);
        }
        System.out.println("No name found");
        return null;
    }

    public String getBasePath(String packageString){
        Pattern basePathPattern = Pattern.compile("^(.*?)v\\d");

        Matcher basePathMatcher = basePathPattern.matcher(packageString);
        if (basePathMatcher.find()) {
            String basePath =  basePathMatcher.group(1);
            if(basePath.charAt(basePath.length()-1) == '.'){
                basePath = basePath.substring(0, basePath.length()-1);
            }
            return basePath;
        }
        System.out.println("No base path found");
        return null;
    }

    // Method to extract service blocks from a given text
    public List<String> extractServiceBlocks(String text) {
        // Regular expression pattern to match the service blocks
        String patternString = "service\\s+\\w+\\s*\\{[^{}]*(?:\\{[^{}]*\\}[^{}]*)*\\}";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString, Pattern.DOTALL);
        Matcher matcher = pattern.matcher(text);

        // Find all matches and append them to the result
        List<String> result = new ArrayList<>();
        while (matcher.find()) {
            result.add(matcher.group());
        }
        return result;
    }

    public List<String> extractMethodNames(String serviceBlock) {
        // Regular expression pattern to match the method names
        String patternString = "(?<=rpc\\s)\\w+";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString);
        Matcher matcher = pattern.matcher(serviceBlock);

        // Find all matches and append them to the result
        List<String> result = new ArrayList<>();
        while (matcher.find()) {
            result.add(matcher.group());
        }
        return result;
    }

    public String getServiceName(String serviceBlock) {
        // Regular expression pattern to match the service name
        String patternString = "(?<=service\\s)\\w+";

        // Compile the regular expression
        Pattern pattern = Pattern.compile(patternString);
        Matcher matcher = pattern.matcher(serviceBlock);

        // Find the first match and return it
        if (matcher.find()) {
            return matcher.group();
        }
        return null;
    }

    public ProtoFile parseProtoContent(String content) {
        ProtoFile protoFile = new ProtoFile();
        protoFile.services = new ArrayList<>();

        List<String> serviceBlocks = extractServiceBlocks(content);
        for (String serviceBlock : serviceBlocks) {
            Service service = new Service();
            service.name = getServiceName(serviceBlock);
            service.methods = new ArrayList<>();
            service.methods.addAll(extractMethodNames(serviceBlock));
            protoFile.services.add(service);
        }

        // Extract package name
        String packageName = getPackageString(content);
        protoFile.packageName = getPackageName(packageName);
        protoFile.version = getVersion(packageName);
        protoFile.basePath = getBasePath(packageName);

//        System.out.println(protoFile);

        return protoFile;
    }

    public List<String> getMethods(Service Service){
        return Service.methods;
    }
    public List<Service> getServices(){
        return this.protoFile.services;
    }
    public class ProtoFile {
        public String packageName;
        public String basePath;
        public String version;
        public List<Service> services;

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

    public class Service {
        public String name;
        public List<String> methods;

        @Override
        public String toString() {
            return " Service{" +
                    "name='" + name + '\'' +
                    ", methods=" + methods +
                    '}';
        }
    }
}