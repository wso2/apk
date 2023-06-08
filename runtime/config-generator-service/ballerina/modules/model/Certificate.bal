public type Certificate record{|
    string certificateId;
    string srtialNumber;
    string issuer;
    string subject;
    string notBefore;
    string notAfter;
    string certificateContent;
    string hostname;
    string 'version;
    boolean active;
|};