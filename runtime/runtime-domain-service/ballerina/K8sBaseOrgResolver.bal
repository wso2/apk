import runtime_domain_service.model;
import wso2/apk_common_lib as commons;

public class K8sBaseOrgResolver {
    *commons:OrganizationResolver;

    public isolated function retrieveOrganizationFromIDPClaimValue(string organizationClaim) returns commons:Organization? {
        OrgClient orgClient = new;
        model:Organization? retrievedOrg = orgClient.retrieveOrganizationFromIDPClaimValue(organizationClaim);
        if retrievedOrg is model:Organization {
            commons:Organization organization = {
                displayName: retrievedOrg.spec.displayName,
                name: retrievedOrg.spec.name,
                organizationClaimValue: retrievedOrg.spec.organizationClaimValue,
                uuid: retrievedOrg.spec.uuid,
                enabled: retrievedOrg.spec.enabled,
                properties: retrievedOrg.spec.properties
            };
            return organization;
        }

        return;
    }
}
