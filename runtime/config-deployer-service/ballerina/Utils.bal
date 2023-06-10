    isolated function getDomain(string url) returns string {
        string hostPort = "";
        string protocol = "";
        if url.startsWith("https://") {
            hostPort = url.substring(8, url.length());
            protocol = "https";
        } else if url.startsWith("http://") {
            hostPort = url.substring(7, url.length());
            protocol = "http";
        } else {
            return "";
        }
        int? indexOfSlash = hostPort.indexOf("/", 0);
        if indexOfSlash is int {
            return protocol + "://" + hostPort.substring(0, indexOfSlash);
        } else {
            return protocol + "://" + hostPort;
        }
    }
        isolated function getPath(string url) returns string {
        string hostPort = "";
        if url.startsWith("https://") {
            hostPort = url.substring(8, url.length());
        } else if url.startsWith("http://") {
            hostPort = url.substring(7, url.length());
        } else {
            return "";
        }
        int? indexOfSlash = hostPort.indexOf("/", 0);
        if indexOfSlash is int {
            return hostPort.substring(indexOfSlash, hostPort.length());
        } else {
            return "";
        }
    }