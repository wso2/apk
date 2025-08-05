    isolated function gethost(string url) returns string {
        string host = "";
        if url.startsWith("https://") {
            host = url.substring(8, url.length());
        } else if url.startsWith("http://") {
            host = url.substring(7, url.length());
        } else {
            return "";
        }
        int? indexOfColon = host.indexOf(":", 0);
        if indexOfColon is int {
            return host.substring(0, indexOfColon);
        } else {
            int? indexOfSlash = host.indexOf("/", 0);
            if indexOfSlash is int {
                return host.substring(0, indexOfSlash);
            } else {
                return host;
            }
        }
    }