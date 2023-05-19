//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

import ballerina/http;
import ballerina/log;
import wso2/apk_common_lib as commons;

service /api/am/admin on ep0 {
    # Retrieve/Search Policies
    #
    # + query - **Search**. You can search by providing a keyword. Allowed to search by type and name only. 
    # + return - OK. List of qualifying Policies is returned. 
    // resource function get policies/search(string? query) returns PolicyDetailsList {
    // }
    # Get all Application Rate Plans
    #
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + return - returns can be any of following types
    # ApplicationRatePlanList (OK. Policies returned)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'application\-rate\-plans(http:RequestContext requestContext, @http:Header string? accept = "application/json") returns ApplicationRatePlanList|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        ApplicationRatePlanList|commons:APKError appPolicyList = getApplicationUsagePlans(organization);
        if appPolicyList is ApplicationRatePlanList {
            log:printDebug(appPolicyList.toString());
        }
        return appPolicyList;
    }
    # Add an Application Rate Plan
    #
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Application level policy object that should to be added 
    # + return - returns can be any of following types
    # ApplicationRatePlan (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post 'application\-rate\-plans(http:RequestContext requestContext, @http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns ApplicationRatePlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        ApplicationRatePlan|commons:APKError createdAppPol = addApplicationUsagePlan(payload, organization);
        if createdAppPol is ApplicationRatePlan {
            log:printDebug(createdAppPol.toString());
        }
        return createdAppPol;
    }
    # Get an Application Rate Plan
    #
    # + planId - Policy UUID 
    # + return - returns can be any of following types
    # ApplicationRatePlan (OK. Plan returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'application\-rate\-plans/[string planId](http:RequestContext requestContext) returns ApplicationRatePlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        ApplicationRatePlan|commons:APKError appPolicy = getApplicationUsagePlanById(planId, organization);
        if appPolicy is ApplicationRatePlan {
            log:printDebug(appPolicy.toString());
        }
        return appPolicy;
    }
    # Update an Application Rate Plan
    #
    # + planId - Policy UUID 
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Policy object that needs to be modified 
    # + return - returns can be any of following types
    # ApplicationRatePlan (OK. Plan updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function put 'application\-rate\-plans/[string planId](http:RequestContext requestContext, @http:Payload ApplicationRatePlan payload, @http:Header string 'content\-type = "application/json") returns ApplicationRatePlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        ApplicationRatePlan|commons:APKError appPolicy = updateApplicationUsagePlan(planId, payload, organization);
        if appPolicy is ApplicationRatePlan {
            log:printDebug(appPolicy.toString());
        }
        return appPolicy;
    }
    # Delete an Application Rate Plan
    #
    # + planId - Policy UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function delete 'application\-rate\-plans/[string planId](http:RequestContext requestContext) returns http:Ok|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string|commons:APKError ex = removeApplicationUsagePlan(planId, organization);
        if ex is commons:APKError {
            return ex;
        } else {
            return http:OK;
        }
    }
    # Get all Business Plans
    #
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + return - returns can be any of following types
    # BusinessPlanList (OK. Plans returned)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'business\-plans(http:RequestContext requestContext, @http:Header string? accept = "application/json") returns BusinessPlanList|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BusinessPlanList|commons:APKError subPolicyList = getBusinessPlans(organization);
        if subPolicyList is BusinessPlanList {
            log:printDebug(subPolicyList.toString());
        }
        return subPolicyList;
    }
    # Add a Business Plan
    #
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Business Plan object that should to be added 
    # + return - returns can be any of following types
    # BusinessPlan (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post 'business\-plans(http:RequestContext requestContext, @http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns BusinessPlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BusinessPlan|commons:APKError createdSubPol = addBusinessPlan(payload, organization);
        if createdSubPol is BusinessPlan {
            log:printDebug(createdSubPol.toString());
        }
        return createdSubPol;
    }
    # Get a Business Plan
    #
    # + planId - Policy UUID 
    # + return - returns can be any of following types
    # BusinessPlan (OK. Plan returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'business\-plans/[string planId](http:RequestContext requestContext) returns BusinessPlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BusinessPlan|commons:APKError subPolicy = getBusinessPlanById(planId, organization);
        if subPolicy is BusinessPlan {
            log:printDebug(subPolicy.toString());
        }
        return subPolicy;
    }
    # Update a Business Plan
    #
    # + planId - Policy UUID 
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Plan object that needs to be modified 
    # + return - returns can be any of following types
    # BusinessPlan (OK. Plan updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function put 'business\-plans/[string planId](http:RequestContext requestContext, @http:Payload BusinessPlan payload, @http:Header string 'content\-type = "application/json") returns BusinessPlan|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return updateBusinessPlan(planId, payload, organization);
    }
    # Delete a Business Plan
    #
    # + planId - Policy UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function delete 'business\-plans/[string planId](http:RequestContext requestContext) returns http:Ok|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string|commons:APKError ex = removeBusinessPlan(planId, organization);
        if ex is commons:APKError {
            return ex;
        } else {
            return http:OK;
        }
    }
    # Export a Throttling Policy
    #
    # + policyId - UUID of the ThrottlingPolicy 
    # + name - Throttling Policy Name 
    # + 'type - Type of the Throttling Policy 
    # + format - Format of output documents. Can be YAML or JSON. 
    # + return - returns can be any of following types
    # ExportPolicy (OK. Export Successful.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function get throttling/policies/export(string? policyId, string? name, string? 'type, string? format) returns ExportPolicy|NotFoundError|InternalServerErrorError {
    // }
    # Import a Throttling Policy
    #
    # + overwrite - Update an existing throttling policy with the same name. 
    # + request - parameter description 
    # + return - returns can be any of following types
    # http:Ok (Created. Throttling Policy Imported Successfully.)
    # ForbiddenError (Forbidden. The request must be conditional but no condition has been specified.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # ConflictError (Conflict. Specified resource already exists.)
    # InternalServerErrorError (Internal Server Error.)
    // resource function post throttling/policies/'import(boolean? overwrite, http:Request request) returns http:Ok|ForbiddenError|NotFoundError|ConflictError|InternalServerErrorError {
    // }
    # Get all Deny Policies
    #
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + return - returns can be any of following types
    # BlockingConditionList (OK. Deny Policies returned)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'deny\-policies(http:RequestContext requestContext, @http:Header string? accept = "application/json") returns BlockingConditionList|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return getAllDenyPolicies(organization);
    }
    # Add a deny policy
    #
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Blocking condition object that should to be added 
    # + return - returns can be any of following types
    # BlockingCondition (Created. Successful response with the newly created object as entity in the body. Location header contains URL of newly created entity.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # UnsupportedMediaTypeError (Unsupported Media Type. The entity of the request was not in a supported format.)
    isolated resource function post 'deny\-policies(http:RequestContext requestContext, @http:Payload BlockingCondition payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BlockingCondition|commons:APKError createdDenyPol = addDenyPolicy(payload, organization);
        if createdDenyPol is BlockingCondition {
            log:printDebug(createdDenyPol.toString());
        }
        return createdDenyPol;
    }
    # Get a Deny Policy
    #
    # + policyId - Policy UUID 
    # + return - returns can be any of following types
    # BlockingCondition (OK. Condition returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'deny\-policies/[string policyId](http:RequestContext requestContext) returns BlockingCondition|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BlockingCondition|commons:APKError denyPolicy = getDenyPolicyById(policyId, organization);
        if denyPolicy is BlockingCondition {
            log:printDebug(denyPolicy.toString());
        }
        return denyPolicy;
    }
    # Delete a Deny Policy
    #
    # + policyId - Policy UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function delete 'deny\-policies/[string policyId](http:RequestContext requestContext) returns http:Ok|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string|commons:APKError ex = removeDenyPolicy(policyId, organization);
        if ex is commons:APKError {
            return ex;
        } else {
            return http:OK;
        }
    }
    # Update a Deny Policy
    #
    # + policyId - Policy UUID 
    # + 'content\-type - Media type of the entity in the body. Default is application/json. 
    # + payload - Blocking condition with updated status 
    # + return - returns can be any of following types
    # BlockingCondition (OK. Resource successfully updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function patch 'deny\-policies/[string policyId](http:RequestContext requestContext, @http:Payload BlockingConditionStatus payload, @http:Header string 'content\-type = "application/json") returns BlockingCondition|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        BlockingCondition|commons:APKError updatedPolicy = updateDenyPolicy(policyId, payload, organization);
        if updatedPolicy is BlockingCondition {
            log:printDebug(updatedPolicy.toString());
        }
        return updatedPolicy;
    }
    # Retrieve/Search Applications
    #
    # + user - username of the application creator 
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + name - Application Name 
    # + tenantDomain - Tenant domain of the applications to get. This has to be specified only if it is required to get applications of a tenant other than the requester's tenant. So, if not specified, the default will be set as the requester's tenant domain. This cross tenant Application access is allowed only for super tenant admin users **only at a migration process**. 
    # + sortBy - parameter description 
    # + sortOrder - parameter description 
    # + return - returns can be any of following types
    # ApplicationList (OK. Application list returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get applications(string? user, string? name, string? tenantDomain, int 'limit = 25, int offset = 0, @http:Header string? accept = "application/json", string sortBy = "name", string sortOrder = "asc") returns ApplicationList|BadRequestError|NotAcceptableError {
    // }
    # Get the details of an Application
    #
    # + applicationId - Application UUID 
    # + return - returns can be any of following types
    # Application (OK. Application details returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get applications/[string applicationId]() returns Application|BadRequestError|NotFoundError|NotAcceptableError {
    // }
    # Delete an Application
    #
    # + applicationId - Application UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Resource successfully deleted.)
    # AcceptedWorkflowResponse (Accepted. The request has been accepted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function delete applications/[string applicationId]() returns http:Ok|AcceptedWorkflowResponse|NotFoundError {
    // }
    # Change Application Owner
    #
    # + applicationId - Application UUID 
    # + owner - parameter description 
    # + return - returns can be any of following types
    # http:Ok (OK. Application owner changed successfully.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function post applications/[string applicationId]/'change\-owner(string owner) returns http:Ok|BadRequestError|NotFoundError {
    // }
    # Get all registered Environments
    #
    # + return - OK. Environments returned 
    // resource function get environments() returns EnvironmentList {
    // }
    # Add an Environment
    #
    # + payload - Environment object that should to be added 
    # + return - returns can be any of following types
    # Environment (Created. Successful response with the newly created environment as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    // resource function post environments(@http:Payload Environment payload) returns Environment|BadRequestError {
    // }
    # Update an Environment
    #
    # + environmentId - Environment UUID (or Environment name defined in config) 
    # + payload - Environment object with updated information 
    # + return - returns can be any of following types
    # Environment (OK. Environment updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function put environments/[string environmentId](@http:Payload Environment payload) returns Environment|BadRequestError|NotFoundError {
    // }
    # Delete an Environment
    #
    # + environmentId - Environment UUID (or Environment name defined in config) 
    # + return - returns can be any of following types
    # http:Ok (OK. Environment successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function delete environments/[string environmentId]() returns http:Ok|NotFoundError {
    // }
    # Get Tenant Id of User
    #
    # + username - The state represents the current state of the tenant. Supported states are [ active, inactive] 
    # + return - returns can be any of following types
    # TenantInfo (OK. Tenant id of the user retrieved.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get 'tenant\-info/[string username]() returns TenantInfo|NotFoundError|NotAcceptableError {
    // }
    # Get Custom URL Info of a Tenant Domain
    #
    # + tenantDomain - The tenant domain name. 
    # + return - returns can be any of following types
    # CustomUrlInfo (OK. Custom url info of the tenant is retrieved.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    // resource function get 'custom\-urls/[string tenantDomain]() returns CustomUrlInfo|NotFoundError|NotAcceptableError {
    // }
    # Get all API Categories
    #
    # + return - OK. Categories returned 
    isolated resource function get 'api\-categories(http:RequestContext requestContext) returns APICategoryList|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return getAllCategoryList(organization);
    }
    # Add API Category
    #
    # + payload - API Category object that should to be added 
    # + return - returns can be any of following types
    # APICategory (Created. Successful response with the newly created object as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    isolated resource function post 'api\-categories(http:RequestContext requestContext, @http:Payload APICategory payload) returns APICategory|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return addAPICategory(payload, organization);
    }
    # Update an API Category
    #
    # + apiCategoryId - API Category UUID 
    # + payload - API Category object with updated information 
    # + return - returns can be any of following types
    # APICategory (OK. Label updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function put 'api\-categories/[string apiCategoryId](http:RequestContext requestContext, @http:Payload APICategory payload) returns APICategory|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        return updateAPICategory(apiCategoryId, payload, organization);
    }
    # Delete an API Category
    #
    # + apiCategoryId - API Category UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. API Category successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function delete 'api\-categories/[string apiCategoryId](http:RequestContext requestContext) returns http:Ok|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        string|commons:APKError ex = removeAPICategory(apiCategoryId, organization);
        if ex is commons:APKError {
            return ex;
        } else {
            return http:OK;
        }
    }
    # Retrieve Admin Settings
    #
    # + return - returns can be any of following types
    # Settings (OK. Settings returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    resource function get settings(http:RequestContext requestContext) returns Settings|NotFoundError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;

        SettingsClient settingsClient = new;
        return settingsClient.getSettings(organization);
    }
    # Get all Key managers
    #
    # + return - OK. KeyManagers returned 
    isolated resource function get 'key\-managers(http:RequestContext requestContext) returns KeyManagerList|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;

        KeyManagerClient keyManagerClient = new ();
        return keyManagerClient.getAllKeyManagersByOrganization(organization);

    }
    # Add a new API Key Manager
    #
    # + payload - Key Manager object that should to be added 
    # + return - returns can be any of following types
    # KeyManager (Created. Successful response with the newly created object as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    isolated resource function post 'key\-managers(@http:Payload KeyManager payload, http:RequestContext requestContext) returns KeyManager|BadRequestError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;

        KeyManagerClient keyManagerClient = new ();
        return check keyManagerClient.addKeyManagerEntryToOrganization(payload, organization);
    }
    # Get a Key Manager Configuration
    #
    # + keyManagerId - Key Manager UUID 
    # + return - returns can be any of following types
    # KeyManager (OK. KeyManager Configuration returned)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get 'key\-managers/[string keyManagerId](http:RequestContext requestContext) returns KeyManager|NotFoundError|NotAcceptableError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        KeyManagerClient keyManagerClient = new ();
        return check keyManagerClient.getKeyManagerById(keyManagerId, organization);
    }
    # Update a Key Manager
    #
    # + keyManagerId - Key Manager UUID 
    # + payload - Key Manager object with updated information 
    # + return - returns can be any of following types
    # KeyManager (OK. Label updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function put 'key\-managers/[string keyManagerId](@http:Payload KeyManager payload, http:RequestContext requestContext) returns KeyManager|BadRequestError|NotFoundError|commons:APKError {

        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        KeyManagerClient keyManagerClient = new ();
        return check keyManagerClient.updateKeyManager(keyManagerId, payload, organization);
    }
    # Delete a Key Manager
    #
    # + keyManagerId - Key Manager UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Key Manager successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function delete 'key\-managers/[string keyManagerId](http:RequestContext requestContext) returns http:Ok|NotFoundError|commons:APKError {
        commons:UserContext authenticatedUserContext = check commons:getAuthenticatedUserContext(requestContext);
        commons:Organization organization = authenticatedUserContext.organization;
        KeyManagerClient keyManagerClient = new ();
        check keyManagerClient.deleteKeyManager(keyManagerId, organization);
        http:Ok okResponse = {};
        return okResponse;
    }
    # Retrieve Well-known information from Key Manager Well-known Endpoint
    #
    # + request - parameter description 
    # + return - OK. KeyManagers returned 
    // resource function post 'key\-managers/discover(http:Request request) returns OkKeyManagerWellKnownResponse {
    // }


    # Retrieve All Pending Workflow Processes
    #
    # + 'limit - Maximum size of resource array to return. 
    # + offset - Starting point within the complete list of items qualified. 
    # + accept - Media types acceptable for the response. Default is application/json. 
    # + workflowType - We need to show the values of each workflow process separately .for that we use workflow type. Workflow type can be APPLICATION_CREATION, SUBSCRIPTION_CREATION etc. 
    # + return - returns can be any of following types
    # WorkflowList (OK. Workflow pendding process list returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get workflows(string? workflowType, int 'limit = 25, int offset = 0, @http:Header string? accept = "application/json") returns WorkflowList|commons:APKError {
        return getWorkflowList(workflowType, 'limit, offset, accept);
    }
    # Update Workflow Status
    #
    # + workflowReferenceId - Workflow reference id 
    # + payload - Workflow event that need to be updated 
    # + return - returns can be any of following types
    # OkWorkflowInfo (OK. Workflow request information is returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    // resource function post workflows/'update\-workflow\-status(string workflowReferenceId, @http:Payload WorkflowInfo payload) returns OkWorkflowInfo|BadRequestError|NotFoundError {
    // }

    # Get all Organization
    #
    # + return - OK. Organization returned 
    isolated resource function get organizations() returns OrganizationList|commons:APKError {
        return getAllOrganization();
    }
    # Add Organization
    #
    # + payload - Organization object that should to be added 
    # + return - returns can be any of following types
    # Organization (Created. Successful response with the newly created object as entity in the body.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    isolated resource function post organizations(@http:Payload Organization payload) returns Organization|commons:APKError {
        return addOrganization(payload);
    }
    # Get the details of an Organization
    #
    # + organizationId - Organization UUID 
    # + return - returns can be any of following types
    # Organization (OK. Application details returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    # NotAcceptableError (Not Acceptable. The requested media type is not supported.)
    isolated resource function get organizations/[string organizationId]() returns Organization|commons:APKError {
        return getOrganizationById(organizationId);
    }
    # Update an Organization
    #
    # + organizationId - Organization UUID 
    # + payload - Organization object with updated information 
    # + return - returns can be any of following types
    # Organization (OK. Label updated.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    isolated resource function put organizations/[string organizationId](@http:Payload Organization payload) returns Organization|commons:APKError {
        return updatedOrganization(organizationId, payload);
    }
    # Delete an Organization
    #
    # + organizationId - Organization UUID 
    # + return - returns can be any of following types
    # http:Ok (OK. Organization successfully deleted.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    resource function delete organizations/[string organizationId]() returns http:Ok|commons:APKError {
        boolean|commons:APKError deleteOrganization = removeOrganization(organizationId);
        if deleteOrganization is commons:APKError {
            return deleteOrganization;
        } else {
            return http:OK;
        }
    }
    # Authenticate Organization info
    #
    # + return - returns can be any of following types
    # Organization (OK. Application details returned.)
    # BadRequestError (Bad Request. Invalid request or validation error.)
    # NotFoundError (Not Found. The specified resource does not exist.)
    resource function get 'organization\-info() returns Organization|commons:APKError {
        return getOrganizationByOrganizationClaim();
    }
}
