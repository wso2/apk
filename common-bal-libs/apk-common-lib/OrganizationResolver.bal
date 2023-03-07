# Description
public type OrganizationResolver object{
    public isolated function retrieveOrganizationFromIDPClaimValue(string organizationClaim) returns Organization|APKError|();
    public isolated function retrieveOrganizationByName(string organizationName) returns Organization|APKError|();
};