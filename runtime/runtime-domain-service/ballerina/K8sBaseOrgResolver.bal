import runtime_domain_service.model;
import wso2/apk_common_lib as commons;

public class K8sBaseOrgResolver {
    *commons:OrganizationResolver;

    public isolated function retrieveOrganizationFromIDPClaimValue(string organizationClaim) returns commons:Organization|commons:APKError|() {
        OrgClient orgClient = new;
        model:Organization? retrievedOrg = orgClient.retrieveOrganizationFromIDPClaimValue(organizationClaim);
        if retrievedOrg is model:Organization {
            return self.convertK8sOrgToGenericOrg(retrievedOrg);
        }

        return;
    }

    public isolated function retrieveOrganizationByName(string organizationName) returns commons:Organization|commons:APKError|() {
        OrgClient orgClient = new;
        model:Organization? organization = orgClient.retrieveOrganizationByName(organizationName);
        if organization is model:Organization {
            return self.convertK8sOrgToGenericOrg(organization);
        }
        return;
    }
    public isolated function convertK8sOrgToGenericOrg(model:Organization k8sOrganization) returns commons:Organization {
        commons:Organization organization = {
            displayName: k8sOrganization.spec.displayName,
            name: k8sOrganization.spec.name,
            organizationClaimValue: k8sOrganization.spec.organizationClaimValue,
            uuid: k8sOrganization.spec.uuid,
            enabled: k8sOrganization.spec.enabled,
            properties: k8sOrganization.spec.properties,
            serviceListingNamespaces: k8sOrganization.spec.serviceListingNamespaces
        };
        return organization;
    }
}
