# Description
public type OrganizationResolver isolated object{
    public isolated function retrieveOrganizationFromIDPClaimValue(map<anydata> claims,string organizationClaim) returns Organization|APKError|();
    public isolated function retrieveOrganizationByName(string organizationName) returns Organization|APKError|();
};