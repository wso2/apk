import ballerina/jwt;
import ballerina/http;

public service class JWTValidationInterceptor {
    *http:RequestInterceptor;
    private final IDPConfiguration idpConfiguration;
    private final jwt:ValidatorConfig jwtValidatorConfig;
    private final OrganizationResolver organizationResolver;
    public function init(IDPConfiguration idpConfiguration, OrganizationResolver organizationResolver) {
        self.idpConfiguration = idpConfiguration.clone();
        self.organizationResolver = organizationResolver;
        self.jwtValidatorConfig = initializeJWTValidator(idpConfiguration.clone());
    }
    resource function 'default [string... path](http:RequestContext ctx, http:Request request, http:Caller caller) returns http:NextService|error? {
        if path[0] == "health" {
            return ctx.next();
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
            if (validatedJWT.hasKey(self.idpConfiguration.organizationClaim)) {
                string organizationClaim = <string>validatedJWT.get(self.idpConfiguration.organizationClaim);
                Organization? retrievedorg = check self.organizationResolver.retrieveOrganizationFromIDPClaimValue(organizationClaim);
                if retrievedorg is Organization {
                    if retrievedorg.enabled {
                        UserContext userContext = {username: <string>validatedJWT.sub, organization: retrievedorg};
                        userContext.claims = self.extractCustomClaims(validatedJWT);
                        return userContext;
                    } else {
                        APKError apkError = error("inactive organization", code = 900951, description = "organization is enactive", statusCode = 401, message = "organization is enactive");
                        return apkError;
                    }
                } else {
                    APKError apkError = error("Organization not found in APK system", code = 900952, description = "Organization not found in APK system", statusCode = 401, message = "Organization not found in APK system");
                    return apkError;
                }
            } else {
                // find default organization
                Organization? retrievedorg = check self.organizationResolver.retrieveOrganizationByName(DEFAULT_ORGANIZATION_NAME);
                if retrievedorg is Organization {
                    UserContext userContext = {username: <string>validatedJWT.sub, organization: retrievedorg};
                    userContext.claims = self.extractCustomClaims(validatedJWT);
                    return userContext;
                } else {
                    APKError apkError = error("Organization not found in APK system", code = 900952, description = "Organization not found in APK system", statusCode = 401, message = "Organization not found in APK system");
                    return apkError;  
                }
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
isolated function initializeJWTValidator(IDPConfiguration idpConfiguration) returns jwt:ValidatorConfig {
    jwt:ValidatorConfig validatorConfig = {issuer: idpConfiguration.issuer};
    string? jwksUrl = idpConfiguration.jwksUrl;
    jwt:ValidatorSignatureConfig signatureConfig = {certFile: idpConfiguration.publicKey.path};
    if jwksUrl is string {
        signatureConfig = {jwksConfig: {url: jwksUrl}};
    }
    validatorConfig.signatureConfig = signatureConfig;
    return validatorConfig;
}
