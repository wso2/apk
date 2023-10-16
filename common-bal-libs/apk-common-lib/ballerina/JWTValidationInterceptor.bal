import ballerina/jwt;
import ballerina/regex;
import ballerina/http;

public isolated service class JWTValidationInterceptor {
    *http:RequestInterceptor;
    private final IDPConfiguration & readonly idpConfiguration;
    private final jwt:ValidatorConfig jwtValidatorConfig;
    private final OrganizationResolver organizationResolver;
    private final string[] & readonly ignoredPaths;
    public isolated function init(IDPConfiguration idpConfiguration, OrganizationResolver organizationResolver, string[] ignoredPaths) {
        self.idpConfiguration = idpConfiguration.cloneReadOnly();
        self.organizationResolver = organizationResolver;
        self.jwtValidatorConfig = initializeJWTValidator(idpConfiguration.cloneReadOnly()).cloneReadOnly();
        self.ignoredPaths = ignoredPaths.cloneReadOnly();
    }
    isolated resource function 'default [string... path](http:RequestContext ctx, http:Request request, http:Caller caller) returns http:NextService|error? {
        string concatPath = "/" + string:'join("/", ...path);
        foreach string ignoredPath in self.ignoredPaths {
            if regex:matches(concatPath, ignoredPath) {
                return ctx.next();
            }
        }
        string|http:HeaderNotFoundError authorizationHeader = request.getHeader(self.idpConfiguration.authorizationHeader);
        if authorizationHeader is string {
            UserContext userContext = check self.validateJWT(authorizationHeader);
            ctx.set(VALIDATED_USER_CONTEXT, userContext.clone());
            return ctx.next();
        } else {
            http:Response response = new ();
            response.statusCode = 401;
            check caller->respond(response);
        }
        return;
    }
    isolated function validateJWT(string header) returns UserContext|APKError {
        jwt:Payload|jwt:Error validatedJWT;
        lock {
            validatedJWT = jwt:validate(header, self.jwtValidatorConfig.clone());
        }
        if validatedJWT is jwt:Payload {
            map<anydata> claims = self.extractCustomClaims(validatedJWT);
            if (validatedJWT.hasKey(self.idpConfiguration.organizationClaim)) {
                string organizationClaim = <string>validatedJWT.get(self.idpConfiguration.organizationClaim);
                Organization? retrievedorg = check self.organizationResolver.retrieveOrganizationFromIDPClaimValue(claims, organizationClaim);
                if retrievedorg is Organization {
                    if retrievedorg.enabled {
                        UserContext userContext = {username: <string>validatedJWT.get(self.idpConfiguration.userClaim), organization: retrievedorg};
                        userContext.claims = claims;
                        return userContext;
                    } else {
                        APKError apkError = error("Inactive Organization", code = 900951, description = "Organization is inactive", statusCode = 401, message = "Organization is inactive");
                        return apkError;
                    }
                } else {
                    APKError apkError = error("Organization not found in APK system", code = 900952, description = "Organization not found in APK system", statusCode = 401, message = "Organization not found in APK system");
                    return apkError;
                }
            } else {
                Organization org = {uuid: "a3b58ccf-6ecc-4557-b5bb-0a35cce38256", name: DEFAULT_ORGANIZATION_NAME, displayName: DEFAULT_ORGANIZATION_NAME, organizationClaimValue: "", enabled: true};
                UserContext userContext = {username: <string>validatedJWT.get(self.idpConfiguration.userClaim), organization: org};
                userContext.claims = claims;
                return userContext;
            }
        }
        else {
            APKError apkError = error("invalid Token", validatedJWT, code = 900954, description = "invalid Token", statusCode = 401, message = "invalid Token");
            return apkError;
        }
    }

    isolated function extractCustomClaims(jwt:Payload jwt) returns map<anydata> {
        map<anydata> customClaims = {};
        map<[string, anydata]> claims = jwt.clone().entries();
        foreach [string, anydata] [claimKey, claimValue] in claims {
            if jwtClaims.indexOf(claimKey) is () {
                customClaims[claimKey] = claimValue.clone();
            }
        }
        return customClaims;
    }

}

# This function used to initialize JWTValidator.
#
# + idpConfiguration - Parameter Description
# + return - Return Value Description
isolated function initializeJWTValidator(IDPConfiguration & readonly idpConfiguration) returns jwt:ValidatorConfig {
    jwt:ValidatorConfig validatorConfig = {issuer: idpConfiguration.issuer};
    string? jwksUrl = idpConfiguration.jwksUrl;
    jwt:ValidatorSignatureConfig signatureConfig = {certFile: idpConfiguration.publicKey.certFilePath};
    if jwksUrl is string {
        signatureConfig = {jwksConfig: {url: jwksUrl}};
    }
    validatorConfig.signatureConfig = signatureConfig;
    return validatorConfig;
}
