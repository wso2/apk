# Description
public type OrganizationResolver object{
    public isolated function retrieveOrganizationFromIDPClaimValue(string organizationClaim) returns Organization?;
    public isolated function retrieveOrganizationByName(string organizationName) returns Organization?;
};