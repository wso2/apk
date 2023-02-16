import runtime_domain_service.model;
import ballerina/http;

public class OrgClient {
    public function retrieveAllOrganizationsAtStartup(map<model:Organization>? organizationsMap, string? continueValue) returns error? {
        string? resultValue = continueValue;
        model:OrganizationList|http:ClientError retrieveAllOrganizationsResult;
        if resultValue is string {
            retrieveAllOrganizationsResult = retrieveAllOrganizations(resultValue);
        } else {
            retrieveAllOrganizationsResult = retrieveAllOrganizations(());
        }

        if retrieveAllOrganizationsResult is model:OrganizationList {
            model:ListMeta metadata = retrieveAllOrganizationsResult.metadata;
            model:Organization[] organizations = retrieveAllOrganizationsResult.items;
            if organizationsMap is map<model:Organization> {
                lock {
                    putAllOrganizations(organizationsMap, organizations.clone());
                }
            } else {
                lock {
                    putAllOrganizations(organizationList, organizations.clone());
                }
            }
            string? continueElement = metadata.'continue;
            if continueElement is string {
                if continueElement.length() > 0 {
                    _ = check self.retrieveAllOrganizationsAtStartup(organizationsMap, continueElement);
                }
            }
            string? resourceVersion = metadata.'resourceVersion;
            if resourceVersion is string {
                setResourceVersion(resourceVersion);
            }
        }
    }
    public isolated function retrieveOrganizationFromIDPClaimValue(string orgClaimValue) returns model:Organization|(){
        lock{
            foreach model:Organization organization in organizationList {
                if organization.spec.organizationClaimValue==orgClaimValue{
                    return organization.cloneReadOnly();
                }
            }
        }    
        return ;    
    }
    public isolated function retrieveOrganizationByName(string orgName) returns model:Organization|(){
        lock{
            foreach model:Organization organization in organizationList {
                if organization.spec.name==orgName{
                    return organization.cloneReadOnly();
                }
            }
        }    
        return ;    
    }
}
