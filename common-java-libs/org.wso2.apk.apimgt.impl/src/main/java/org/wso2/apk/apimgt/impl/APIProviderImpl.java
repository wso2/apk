/*
*  Copyright (c) 2005-2013, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/

package org.wso2.apk.apimgt.impl;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.json.simple.JSONObject;
import org.wso2.apk.apimgt.api.*;
import org.wso2.apk.apimgt.api.doc.model.APIResource;
import org.wso2.apk.apimgt.api.dto.*;
import org.wso2.apk.apimgt.api.model.*;
import org.wso2.apk.apimgt.api.model.policy.APIPolicy;
import org.wso2.apk.apimgt.api.model.policy.ApplicationPolicy;
import org.wso2.apk.apimgt.api.model.policy.GlobalPolicy;
import org.wso2.apk.apimgt.api.model.policy.SubscriptionPolicy;
import org.wso2.apk.apimgt.impl.dao.dto.DocumentContent;
import org.wso2.apk.apimgt.impl.dao.dto.DocumentSearchResult;
import org.wso2.apk.apimgt.impl.dao.dto.Organization;
import org.wso2.apk.apimgt.impl.dao.dto.PublisherAPI;
import org.wso2.apk.apimgt.impl.dao.dto.PublisherAPIInfo;
import org.wso2.apk.apimgt.impl.dao.dto.PublisherAPISearchResult;
import org.wso2.apk.apimgt.impl.dao.dto.UserContext;
import org.wso2.apk.apimgt.impl.dao.exceptions.DocumentationPersistenceException;
import org.wso2.apk.apimgt.impl.dao.exceptions.ThumbnailPersistenceException;
import org.wso2.apk.apimgt.impl.dao.mapper.APIMapper;
import org.wso2.apk.apimgt.impl.dao.mapper.DocumentMapper;
import org.wso2.apk.apimgt.impl.lifecycle.LCManagerFactory;
import org.wso2.apk.apimgt.impl.utils.APIUtil;

import java.io.ByteArrayInputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * This class provides the core API provider functionality. It is implemented in a very
 * self-contained and 'pure' manner, without taking requirements like security into account,
 * which are subject to frequent change. Due to this 'pure' nature and the significance of
 * the class to the overall API management functionality, the visibility of the class has
 * been reduced to package level. This means we can still use it for internal purposes and
 * possibly even extend it, but it's totally off the limits of the users. Users wishing to
 * pragmatically access this functionality should use one of the extensions of this
 * class which is visible to them. These extensions may add additional features like
 * security to this class.
 */
class APIProviderImpl extends AbstractAPIManager implements APIProvider {

    private static final Log log = LogFactory.getLog(APIProviderImpl.class);

    private final String userNameWithoutChange;

    public APIProviderImpl(String username) throws APIManagementException {
        super(username);
        this.userNameWithoutChange = username;
//        certificateManager = CertificateManagerImpl.getInstance();
//        this.artifactSaver = ServiceReferenceHolder.getInstance().getArtifactSaver();
//        this.importExportAPI = ServiceReferenceHolder.getInstance().getImportExportService();
//        this.gatewayArtifactsMgtDAO = GatewayArtifactsMgtDAO.getInstance();
//        this.recommendationEnvironment = ServiceReferenceHolder.getInstance().getAPIManagerConfigurationService()
//                .getAPIManagerConfiguration().getApiRecommendationEnvironment();
//        globalMediationPolicyImpl = new GlobalMediationPolicyImpl(organization);
    }


    @Override
    public String addComment(String uuid, Comment comment, String user) throws APIManagementException {
        return commentDAOImpl.addComment(uuid, comment, user);
    }

    @Override
    public Comment getComment(ApiTypeWrapper apiTypeWrapper, String commentId, Integer replyLimit, Integer replyOffset)
            throws APIManagementException {
        Comment comment = commentDAOImpl.getComment(apiTypeWrapper, commentId, replyLimit, replyOffset);
        if (comment != null) {
            return comment;
        } else {
            throw new APIManagementException(ExceptionCodes.COMMENT_NOT_FOUND);
        }
    }

    @Override
    public CommentList getComments(ApiTypeWrapper apiTypeWrapper, String parentCommentID,
                                   Integer replyLimit, Integer replyOffset) throws
            APIManagementException {
        return commentDAOImpl.getComments(apiTypeWrapper, parentCommentID, replyLimit, replyOffset);
    }

    @Override
    public boolean editComment(ApiTypeWrapper apiTypeWrapper, String commentId, Comment comment) throws
            APIManagementException {
        return commentDAOImpl.editComment(apiTypeWrapper, commentId, comment);
    }

    @Override
    public boolean deleteComment(ApiTypeWrapper apiTypeWrapper, String commentId) throws APIManagementException {
        return commentDAOImpl.deleteComment(apiTypeWrapper, commentId);
    }

    @Override
    public Set<Subscriber> getSubscribersOfProvider(String providerId) throws APIManagementException {
        return null;
    }

    @Override
    public Usage getUsageByAPI(APIIdentifier apiIdentifier) {
        return null;
    }

    @Override
    public Usage getAPIUsageByUsers(String providerId, String apiName) {
        return null;
    }

    @Override
    public UserApplicationAPIUsage[] getAllAPIUsageByProvider(String providerId) throws APIManagementException {
        return new UserApplicationAPIUsage[0];
    }

    @Override
    public List<SubscribedAPI> getAPIUsageByAPIId(String uuid, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public List<SubscribedAPI> getAPIProductUsageByAPIProductId(APIProductIdentifier apiProductId) throws APIManagementException {
        return null;
    }

    @Override
    public Usage getAPIUsageBySubscriber(APIIdentifier apiIdentifier, String consumerEmail) {
        return null;
    }

    @Override
    public List<SubscribedAPI> getSubscriptionsOfAPI(String apiName, String apiVersion, String provider) throws APIManagementException {
        return null;
    }

    @Override
    public String getSubscriber(String subscriptionId) throws APIManagementException {
        return null;
    }

    @Override
    public Map getSubscriberClaims(String subscriber) throws APIManagementException {
        return null;
    }

    @Override
    public void deleteSubscriptionBlockCondition(String conditionValue) throws APIManagementException {

    }

    @Override
    public String getAPIContext(String uuid) throws APIManagementException {
        return null;
    }

    @Override
    public APIPolicy getAPIPolicy(String username, String policyName) throws APIManagementException {
        return null;
    }

    @Override
    public ApplicationPolicy getApplicationPolicy(String username, String policyName) throws APIManagementException {
        return null;
    }

    @Override
    public SubscriptionPolicy getSubscriptionPolicy(String username, String policyName) throws APIManagementException {
        return null;
    }

    @Override
    public GlobalPolicy getGlobalPolicy(String policyName) throws APIManagementException {
        return null;
    }

    @Override
    public boolean isGlobalPolicyKeyTemplateExists(GlobalPolicy policy) throws APIManagementException {
        return false;
    }

    @Override
    public API addAPI(API api) throws APIManagementException {
        return null;
    }

    @Override
    public ApiTypeWrapper getAPIorAPIProductByUUID(String uuid, String requestedTenantDomain) throws APIManagementException {
        APIInfo apiInfo = apiDAOImpl.getAPIInfoByUUID(uuid);
        if (apiInfo != null) {
            if (apiInfo.getOrganization().equals(requestedTenantDomain)) {
                if (APIConstants.API_PRODUCT.equals(apiInfo.getApiType())) {
                    return new ApiTypeWrapper(getAPIProductbyUUID(uuid, requestedTenantDomain));
                } else {
                    return new ApiTypeWrapper(getAPIbyUUID(uuid, requestedTenantDomain));
                }
            } else {
                throw new APIManagementException(
                        "User " + username + " does not have permission to view API Product : " + uuid,
                        ExceptionCodes.NO_READ_PERMISSIONS);
            }
        } else {
            String msg = "Failed to get API. API artifact corresponding to artifactId " + uuid + " does not exist";
            throw new APIManagementException(msg, ExceptionCodes.NO_API_ARTIFACT_FOUND);
        }
    }

    public APIProduct getAPIProductbyUUID(String uuid, String organization) throws APIManagementException {
        //TODO:APK
        return null;
//        try {
//            Organization org = new Organization(organization);
//            PublisherAPIProduct publisherAPIProduct = apiPersistenceInstance.getPublisherAPIProduct(org, uuid);
//            if (publisherAPIProduct != null) {
//                APIProduct product = APIProductMapper.INSTANCE.toApiProduct(publisherAPIProduct);
//                product.setID(new APIProductIdentifier(publisherAPIProduct.getProviderName(),
//                        publisherAPIProduct.getApiProductName(), publisherAPIProduct.getVersion(), uuid));
//                checkAccessControlPermission(userNameWithoutChange, product.getAccessControl(),
//                        product.getAccessControlRoles());
//                populateRevisionInformation(product, uuid);
//                populateAPIProductInformation(uuid, organization, product);
//                populateAPIStatus(product);
//                populateAPITier(product);
//                return product;
//            } else {
//                String msg = "Failed to get API Product. API Product artifact corresponding to artifactId " + uuid
//                        + " does not exist";
//                throw new APIManagementException(msg, ExceptionCodes.from(ExceptionCodes.API_PRODUCT_NOT_FOUND, uuid));
//            }
//        } catch (APIPersistenceException | OASPersistenceException | ParseException e) {
//            String msg = "Failed to get API Product";
//            throw new APIManagementException(msg, e, ExceptionCodes.INTERNAL_ERROR);
//        }
    }

    @Override
    public void deleteAPIProduct(APIProduct apiProduct) throws APIManagementException {

    }

    @Override
    public String addAPIRevision(APIRevision apiRevision, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public APIRevision getAPIRevision(String revisionUUID) throws APIManagementException {
        return null;
    }

    @Override
    public String getAPIRevisionUUID(String revisionNum, String apiUUID) throws APIManagementException {
        return null;
    }

    @Override
    public String getAPIRevisionUUIDByOrganization(String revisionNum, String apiUUID, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public String getEarliestRevisionUUID(String apiUUID) throws APIManagementException {
        return null;
    }

    @Override
    public String getLatestRevisionUUID(String apiUUID) throws APIManagementException {
        return null;
    }

    @Override
    public List<APIRevision> getAPIRevisions(String apiUUID) throws APIManagementException {
        return null;
    }

    @Override
    public void removeUnDeployedAPIRevision(String apiId, String apiRevisionUUID, String environment) throws APIManagementException {

    }

    @Override
    public void updateAPIDisplayOnDevportal(String apiId, String apiRevisionId, APIRevisionDeployment apiRevisionDeployment) throws APIManagementException {

    }

    @Override
    public void updateAPIProductDisplayOnDevportal(String apiProductId, String apiRevisionId, APIRevisionDeployment apiRevisionDeployment) throws APIManagementException {

    }

    @Override
    public APIRevisionDeployment getAPIRevisionDeployment(String name, String revisionId) throws APIManagementException {
        return null;
    }


    @Override
    public API getAPIbyUUID(String uuid, String organization) throws APIManagementException {
        //TODO:APK
        return null;
    }

    @Override
    public APISearchResult searchPaginatedAPIsByFQDN(String endpoint, String tenantDomain, int start, int end) throws APIManagementException {
        return null;
    }

    @Override
    public Map<String, Object> searchPaginatedAPIs(String searchQuery, String organization, int start, int end,
                                                   String sortBy, String sortOrder) throws APIManagementException {
        Map<String, Object> result = new HashMap<String, Object>();
        if (log.isDebugEnabled()) {
            log.debug("Original search query received : " + searchQuery);
        }

        Organization org = new Organization(organization);
        String[] roles = APIUtil.getFilteredUserRoles(userNameWithoutChange);
        Map<String, Object> properties = APIUtil.getUserProperties(userNameWithoutChange);
        UserContext userCtx = new UserContext(userNameWithoutChange, org, properties, roles);
        try {
            PublisherAPISearchResult searchAPIs = apiDAOImpl.searchAPIsForPublisher(org, searchQuery,
                    start, end, userCtx, sortBy, sortOrder);
            if (log.isDebugEnabled()) {
                log.debug("searched APIs for query : " + searchQuery + " :-->: " + searchAPIs.toString());
            }
            Set<Object> apiSet = new LinkedHashSet<>();
            if (searchAPIs != null) {
                List<PublisherAPIInfo> list = searchAPIs.getPublisherAPIInfoList();
                List<Object> apiList = new ArrayList<>();
                for (PublisherAPIInfo publisherAPIInfo : list) {
                    // API mappedAPI = APIMapper.INSTANCE.toApi(publisherAPIInfo);
                    APIIdentifier apiIdentifier = new APIIdentifier(publisherAPIInfo.getProviderName(), publisherAPIInfo.getApiName(), publisherAPIInfo.getVersion(), publisherAPIInfo.getId());
                    API mappedAPI = new API(apiIdentifier);
                    mappedAPI.setContext(publisherAPIInfo.getContext());
                    mappedAPI.setId(apiIdentifier);
                    mappedAPI.setUuid(publisherAPIInfo.getId());
                    mappedAPI.setContextTemplate(publisherAPIInfo.getContext());
                    populateAPIStatus(mappedAPI);
                    populateDefaultVersion(mappedAPI);
                    apiList.add(mappedAPI);
                }
                apiSet.addAll(apiList);
                result.put("apis", apiSet);
                result.put("length", searchAPIs.getTotalAPIsCount());
                result.put("isMore", true);
            } else {
                result.put("apis", apiSet);
                result.put("length", 0);
                result.put("isMore", false);
            }
        } catch (APIManagementException e) {
            throw new APIManagementException("Error while searching the api ", e, ExceptionCodes.INTERNAL_ERROR);
        }
        return result ;
    }

    @Override
    public Map<String, Object> searchPaginatedContent(String searchQuery, String orgId, int start, int end) throws APIManagementException {
        return null;
    }

    private void populateAPIStatus(API api) throws APIManagementException {
        if (api.isRevision()) {
            api.setStatus(apiDAOImpl.getAPIStatusFromAPIUUID(api.getRevisionedApiId()));
        } else {
            api.setStatus(apiDAOImpl.getAPIStatusFromAPIUUID(api.getUuid()));
        }
    }

    protected void populateDefaultVersion(API api) throws APIManagementException {

        apiMgtDAO.setDefaultVersion(api);
    }

    @Override
    public Documentation addDocumentation(String uuid, Documentation documentation, String organization) throws APIManagementException {
        if (documentation != null) {
            org.wso2.apk.apimgt.impl.dao.dto.Documentation mappedDoc = DocumentMapper.INSTANCE
                    .toDocumentation(documentation);
            try {
                org.wso2.apk.apimgt.impl.dao.dto.Documentation addedDoc = apiDAOImpl.addDocumentation(
                        new Organization(organization), uuid, mappedDoc);
                if (addedDoc != null) {
                    return DocumentMapper.INSTANCE.toDocumentation(addedDoc);
                }
            } catch (DocumentationPersistenceException e) {
                handleExceptionWithCode("Failed to add documentation", e, ExceptionCodes.INTERNAL_ERROR);
            }
        }
        return null;
    }

    @Override
    public boolean isDocumentationExist(String uuid, String docName, String organization) throws APIManagementException {
        boolean exist = false;
        UserContext ctx = null;
        try {
            DocumentSearchResult result = apiDAOImpl.searchDocumentation(new Organization(organization), uuid, 0, 0,
                    "name:" + docName, ctx);
            if (result != null && result.getDocumentationList() != null && !result.getDocumentationList().isEmpty()) {
                String returnDocName = result.getDocumentationList().get(0).getName();
                if (returnDocName != null && returnDocName.equals(docName)) {
                    exist = true;
                }
            }
        } catch (DocumentationPersistenceException e) {
            handleExceptionWithCode("Failed to search documentation for name " + docName, e,
                    ExceptionCodes.INTERNAL_ERROR);
        }
        return exist;
    }

    @Override
    public void addWSDLResource(String apiId, ResourceFile resource, String url, String organization) throws APIManagementException {

    }

    /**
     * Updates a given documentation
     *
     * @param apiId         id of the document
     * @param documentation Documentation
     * @param organization identifier of the organization
     * @return updated documentation Documentation
     * @throws APIManagementException if failed to update docs
     */
    @Override
    public Documentation updateDocumentation(String apiId, Documentation documentation, String organization) throws APIManagementException {

        if (documentation != null) {
            org.wso2.apk.apimgt.impl.dao.dto.Documentation mappedDoc = DocumentMapper.INSTANCE
                    .toDocumentation(documentation);
            try {
                org.wso2.apk.apimgt.impl.dao.dto.Documentation updatedDoc = apiDAOImpl
                        .updateDocumentation(new Organization(organization), apiId, mappedDoc);
                if (updatedDoc != null) {
                    return DocumentMapper.INSTANCE.toDocumentation(updatedDoc);
                }
            } catch (DocumentationPersistenceException e) {
                handleExceptionWithCode("Failed to add documentation", e, ExceptionCodes.INTERNAL_ERROR);
            }
        }
        return null;
    }

    @Override
    public void addDocumentationContent(String uuid, String docId, String organization, DocumentationContent content)
            throws APIManagementException {
        DocumentContent mappedContent = null;
        try {
            mappedContent = DocumentMapper.INSTANCE.toDocumentContent(content);
            DocumentContent doc = apiDAOImpl.addDocumentationContent(new Organization(organization), uuid, docId,
                    mappedContent);
        } catch (DocumentationPersistenceException e) {
            throw new APIManagementException("Error while adding content to doc " + docId,
                    ExceptionCodes.INTERNAL_ERROR);
        }
    }

    @Override
    public void removeDocumentation(String apiId, String docId, String organization) throws APIManagementException {
        try {
            apiDAOImpl.deleteDocumentation(new Organization(organization), apiId, docId);
        } catch (DocumentationPersistenceException e) {
            throw new APIManagementException("Error while deleting the document " + docId,
                    ExceptionCodes.INTERNAL_ERROR);
        }

    }

    @Override
    public List<ResourcePath> getResourcePathsOfAPI(APIIdentifier apiId) throws APIManagementException {
        return apiDAOImpl.getResourcePathsOfAPI(apiId);
    }

    @Override
    public void deleteWorkflowTask(Identifier identifier) throws APIManagementException {

    }

    @Override
    public JSONObject getSecurityAuditAttributesFromConfig(String userId) throws APIManagementException {
        return null;
    }

    @Override
    public List<APIResource> getRemovedProductResources(Set<URITemplate> updatedUriTemplates, API existingAPI) {
        return null;
    }

    @Override
    public boolean isSharedScopeNameExists(String scopeName, int tenantId) throws APIManagementException {
        return false;
    }

    @Override
    public String addSharedScope(Scope scope, String tenantDomain) throws APIManagementException {
        return null;
    }

    @Override
    public List<Scope> getAllSharedScopes(String tenantDomain) throws APIManagementException {
        return null;
    }

    @Override
    public Set<String> getAllSharedScopeKeys(String tenantDomain) throws APIManagementException {
        return null;
    }

    @Override
    public Scope getSharedScopeByUUID(String sharedScopeId, String tenantDomain) throws APIManagementException {
        return null;
    }

    @Override
    public void deleteSharedScope(String scopeName, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void updateSharedScope(Scope sharedScope, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void validateSharedScopes(Set<Scope> scopes, String tenantDomain) throws APIManagementException {

    }

    @Override
    public SharedScopeUsage getSharedScopeUsage(String uuid, int tenantId) throws APIManagementException {
        return null;
    }

    @Override
    public List<APIResource> getUsedProductResources(String uuid) throws APIManagementException {
        return null;
    }

    @Override
    public void deleteAPI(String apiUuid, String organization) throws APIManagementException {

    }

    @Override
    public List<APIRevisionDeployment> getAPIRevisionDeploymentList(String revisionUUID) throws APIManagementException {
        return apiDAOImpl.getAPIRevisionDeploymentByRevisionUUID(revisionUUID);
    }

    @Override
    public void saveAsyncApiDefinition(API api, String jsonText) throws APIManagementException {

    }

    @Override
    public String generateApiKey(String apiId) throws APIManagementException {
        return null;
    }

    @Override
    public List<APIRevisionDeployment> getAPIRevisionsDeploymentList(String apiId) throws APIManagementException {
        return null;
    }

    @Override
    public void addEnvironmentSpecificAPIProperties(String apiUuid, String envUuid, EnvironmentPropertiesDTO environmentPropertyDTO) throws APIManagementException {

    }

    @Override
    public EnvironmentPropertiesDTO getEnvironmentSpecificAPIProperties(String apiUuid, String envUuid) throws APIManagementException {
        return null;
    }

    @Override
    public Environment getEnvironment(String organization, String uuid) throws APIManagementException {
        return null;
    }

    @Override
    public void setOperationPoliciesToURITemplates(String apiId, Set<URITemplate> uriTemplates) throws APIManagementException {

    }

    @Override
    public String importOperationPolicy(OperationPolicyData operationPolicyData, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public String addAPISpecificOperationPolicy(String apiUUID, OperationPolicyData operationPolicyData, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public String addCommonOperationPolicy(OperationPolicyData operationPolicyData, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public OperationPolicyData getAPISpecificOperationPolicyByPolicyName(String policyName, String policyVersion, String apiUUID, String revisionUUID, String organization, boolean isWithPolicyDefinition) throws APIManagementException {
        return null;
    }

    @Override
    public OperationPolicyData getCommonOperationPolicyByPolicyName(String policyName, String policyVersion, String organization, boolean isWithPolicyDefinition) throws APIManagementException {
        return null;
    }

    @Override
    public OperationPolicyData getAPISpecificOperationPolicyByPolicyId(String policyId, String apiUUID, String organization, boolean isWithPolicyDefinition) throws APIManagementException {
        return null;
    }

    @Override
    public OperationPolicyData getCommonOperationPolicyByPolicyId(String policyId, String organization, boolean isWithPolicyDefinition) throws APIManagementException {
        return null;
    }

    @Override
    public void updateOperationPolicy(String operationPolicyId, OperationPolicyData operationPolicyData, String organization) throws APIManagementException {

    }

    @Override
    public List<OperationPolicyData> getAllCommonOperationPolicies(String organization) throws APIManagementException {
        return null;
    }

    @Override
    public List<OperationPolicyData> getAllAPISpecificOperationPolicies(String apiUUID, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public void deleteOperationPolicyById(String policyId, String organization) throws APIManagementException {

    }

    @Override
    public void loadMediationPoliciesToAPI(API api, String organization) throws APIManagementException {

    }

    @Override
    public APIRevision checkAPIUUIDIsARevisionUUID(String apiUUID) throws APIManagementException {
        return null;
    }

    @Override
    public APIProduct getAPIProduct(APIProductIdentifier identifier) throws APIManagementException {
        return null;
    }

    @Override
    public Map<String, Object> searchPaginatedAPIProducts(String searchQuery, String tenantDomain, int start, int end) throws APIManagementException {
        return null;
    }


    /**
     * Get minimal details of API by registry artifact id
     *
     * @param uuid Registry artifact id
     * @param organization identifier of the organization
     * @return API of the provided artifact id
     * @throws APIManagementException
     */
    @Override
    public API getLightweightAPIByUUID(String uuid, String organization) throws APIManagementException {
        try {
            Organization org = new Organization(organization);
            PublisherAPI publisherAPI = apiDAOImpl.getPublisherAPI(org, uuid);
            if (publisherAPI != null) {
                API api = APIMapper.INSTANCE.toApi(publisherAPI);
                api.setOrganization(organization);
                String tiers = null;
                return api;
            } else {
                String msg = "Failed to get API. API artifact corresponding to artifactId " + uuid + " does not exist";
                throw new APIMgtResourceNotFoundException(msg, ExceptionCodes.NO_API_ARTIFACT_FOUND);
            }
        } catch (APIManagementException e) {
            String msg = "Failed to get API with uuid " + uuid;
            throw new APIManagementException(msg, e, ExceptionCodes.INTERNAL_ERROR);
        }
    }

    public API updateAPI(API api, API existingAPI) throws APIManagementException {

        //TODO:APK
        return api;
    }

    @Override
    public API createNewAPIVersion(String apiId, String newVersion, Boolean defaultVersion, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public String retrieveServiceKeyByApiId(int apiId, int tenantId) throws APIManagementException {
        return null;
    }

    @Override
    public void setThumbnailToAPI(String apiId, ResourceFile resource, String organization) throws APIManagementException {

        try {
            org.wso2.apk.apimgt.impl.dao.dto.ResourceFile iconResourceFile = new org.wso2.apk.apimgt.impl.dao.dto.ResourceFile(
                    resource.getContent(), resource.getContentType());
            apiDAOImpl.saveThumbnail(new Organization(organization), apiId, iconResourceFile);
        } catch (ThumbnailPersistenceException e) {
            if (e.getErrorHandler() == ExceptionCodes.API_NOT_FOUND) {
                throw new APIMgtResourceNotFoundException(e);
            } else {
                String errorMessage = "Error while saving thumbnail ";
                throw new APIManagementException(errorMessage, e,
                        ExceptionCodes.from(ExceptionCodes.INTERNAL_ERROR_WITH_SPECIFIC_MESSAGE, errorMessage));
            }
        }
    }

    @Override
    public void saveGraphqlSchemaDefinition(String apiId, String definition, String orgId) throws APIManagementException {

    }

    /**
     * This method returns the lifecycle data for an API including current state,next states.
     *
     * @param uuid  ID of the API
     * @param orgId Identifier of an organization
     * @return Map<String, Object> a map with lifecycle data
     */
    public Map<String, Object> getAPILifeCycleData(String uuid, String orgId) throws APIManagementException {

        API api = getLightweightAPIByUUID(uuid, orgId);
        return getApiOrApiProductLifecycleData(api.getStatus());
    }

    @Override
    public String[] getPolicyNames(String username, String level) throws APIManagementException {
        return new String[0];
    }

    @Override
    public BlockConditionsDTO getBlockCondition(int conditionId) throws APIManagementException {
        return null;
    }

    @Override
    public BlockConditionsDTO getBlockConditionByUUID(String uuid) throws APIManagementException {
        return null;
    }

    @Override
    public boolean updateBlockCondition(int conditionId, String state) throws APIManagementException {
        return false;
    }

    @Override
    public String addBlockCondition(String conditionType, String conditionValue) throws APIManagementException {
        return null;
    }

    @Override
    public boolean deleteBlockConditionByUUID(String uuid) throws APIManagementException {
        return false;
    }

    @Override
    public String getExternalWorkflowReferenceId(int subscriptionId) throws APIManagementException {
        return null;
    }

    @Override
    public boolean isConfigured() {
        return false;
    }

    @Override
    public CertificateMetadataDTO getCertificate(String alias) throws APIManagementException {
        return null;
    }

    @Override
    public List<CertificateMetadataDTO> getCertificates(String userName) throws APIManagementException {
        return null;
    }

    @Override
    public List<CertificateMetadataDTO> searchCertificates(int tenantId, String alias, String endpoint) throws APIManagementException {
        return null;
    }

    @Override
    public List<ClientCertificateDTO> searchClientCertificates(int tenantId, String alias, APIIdentifier apiIdentifier, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public List<ClientCertificateDTO> searchClientCertificates(int tenantId, String alias, APIProductIdentifier apiProductIdentifier, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public int getCertificateCountPerTenant(int tenantId) throws APIManagementException {
        return 0;
    }

    @Override
    public int getClientCertificateCount(int tenantId) throws APIManagementException {
        return 0;
    }

    @Override
    public boolean isCertificatePresent(int tenantId, String alias) throws APIManagementException {
        return false;
    }

    @Override
    public ClientCertificateDTO getClientCertificate(String alias, ApiTypeWrapper apiTypeWrapper, String organization) throws APIManagementException {
        return null;
    }

    @Override
    public CertificateInformationDTO getCertificateStatus(String alias) throws APIManagementException {
        return null;
    }

    @Override
    public int updateCertificate(String certificateString, String alias) throws APIManagementException {
        return 0;
    }

    @Override
    public int updateClientCertificate(String certificate, String alias, ApiTypeWrapper apiTypeWrapper, String tier, int tenantId, String organization) throws APIManagementException {
        return 0;
    }

    @Override
    public ByteArrayInputStream getCertificateContent(String alias) throws APIManagementException {
        return null;
    }

    @Override
    public Map<API, List<APIProductResource>> addAPIProductWithoutPublishingToGateway(APIProduct product) throws APIManagementException {
        return null;
    }

    @Override
    public void saveToGateway(APIProduct product) throws FaultGatewaysException, APIManagementException {

    }

    @Override
    public void deleteAPIProduct(APIProductIdentifier identifier, String apiProductUUID, String organization) throws APIManagementException {

    }

    @Override
    public Map<API, List<APIProductResource>> updateAPIProduct(APIProduct product) throws APIManagementException, FaultGatewaysException {
        return null;
    }

    private Map<String, Object> getApiOrApiProductLifecycleData(String status) throws APIManagementException {

        Map<String, Object> lcData = new HashMap<String, Object>();
        List<String> actionsList;
        actionsList = LCManagerFactory.getInstance().getLCManager().getAllowedActionsForState(status);
        if (actionsList != null) {
            String[] actionsArray = new String[actionsList.size()];
            actionsArray = actionsList.toArray(actionsArray);
            lcData.put(APIConstants.LC_NEXT_STATES, actionsArray);
        }
        status = status.substring(0, 1).toUpperCase() + status.substring(1).toLowerCase(); // First letter capital
        lcData.put(APIConstants.LC_STATUS, status);
        return lcData;
    }


    /**
     * This method is to change registry lifecycle states for an API or API Product artifact
     *
     * @param orgId          UUID of the organization
     * @param apiTypeWrapper API Type Wrapper
     * @param action         Action which need to execute from registry lifecycle
     * @param checklist      checklist items
     * @return APIStateChangeResponse API workflow state and WorkflowResponse
     */
    @Override
    public APIStateChangeResponse changeLifeCycleStatus(String orgId, ApiTypeWrapper apiTypeWrapper, String action,
                                                        Map<String, Boolean> checklist) throws APIManagementException{
        APIStateChangeResponse response = new APIStateChangeResponse();
        //TODO:APK
//        try {
//            String targetStatus;
//            String providerName;
//            String apiName;
//            String apiContext;
//            String apiType;
//            String apiVersion;
//            String currentStatus;
//            String uuid;
//            int apiOrApiProductId;
//            boolean isApiProduct = apiTypeWrapper.isAPIProduct();
//            String workflowType;
//
//            if (isApiProduct) {
//                APIProduct apiProduct = apiTypeWrapper.getApiProduct();
//                providerName = apiProduct.getId().getProviderName();
//                apiName = apiProduct.getId().getName();
//                apiContext = apiProduct.getContext();
//                apiType = apiProduct.getType();
//                apiVersion = apiProduct.getId().getVersion();
//                currentStatus = apiProduct.getState();
//                uuid = apiProduct.getUuid();
//                apiOrApiProductId = apiDAOImpl.getAPIProductId(apiTypeWrapper.getApiProduct().getId());
//                workflowType = WorkflowConstants.WF_TYPE_AM_API_PRODUCT_STATE;
//            } else {
//                API api = apiTypeWrapper.getApi();
//                providerName = api.getId().getProviderName();
//                apiName = api.getId().getApiName();
//                apiContext = api.getContext();
//                apiType = api.getType();
//                apiVersion = api.getId().getVersion();
//                currentStatus = api.getStatus();
//                uuid = api.getUuid();
//                apiOrApiProductId = apiDAOImpl.getAPIID(uuid);
//                workflowType = WorkflowConstants.WF_TYPE_AM_API_STATE;
//            }
//            String gatewayVendor = apiDAOImpl.getGatewayVendorByAPIUUID(uuid);
//
//            WorkflowStatus apiWFState = null;
//            WorkflowDTO wfDTO = workflowDAOImpl.retrieveWorkflowFromInternalReference(Integer.toString(apiOrApiProductId),
//                    workflowType);
//            if (wfDTO != null) {
//                apiWFState = wfDTO.getStatus();
//            }
//
//            // if the workflow has started, then executor should not fire again
//            if (!WorkflowStatus.CREATED.equals(apiWFState)) {
//                response = executeStateChangeWorkflow(currentStatus, action, apiName, apiContext, apiType,
//                        apiVersion, providerName, apiOrApiProductId, uuid, gatewayVendor, workflowType);
//                // get the workflow state once the executor is executed.
//                wfDTO = workflowDAOImpl.retrieveWorkflowFromInternalReference(Integer.toString(apiOrApiProductId),
//                        workflowType);
//                if (wfDTO != null) {
//                    apiWFState = wfDTO.getStatus();
//                    response.setStateChangeStatus(apiWFState.toString());
//                } else {
//                    response.setStateChangeStatus(WorkflowStatus.APPROVED.toString());
//                }
//            }
//
//            // only change the lifecycle if approved
//            // apiWFState is null when simple wf executor is used because wf state is not stored in the db.
//            if (WorkflowStatus.APPROVED.equals(apiWFState) || apiWFState == null) {
//                //TODO:APK
////                LifeCycleUtils.changeLifecycle(this.username, this, orgId, apiTypeWrapper, action, checklist);
//                JSONObject apiLogObject = new JSONObject();
//                apiLogObject.put(APIConstants.AuditLogConstants.NAME, apiName);
//                apiLogObject.put(APIConstants.AuditLogConstants.CONTEXT, apiContext);
//                apiLogObject.put(APIConstants.AuditLogConstants.VERSION, apiVersion);
//                apiLogObject.put(APIConstants.AuditLogConstants.PROVIDER, providerName);
//                APIUtil.logAuditMessage(APIConstants.AuditLogConstants.API, apiLogObject.toString(),
//                        APIConstants.AuditLogConstants.LIFECYCLE_CHANGED, this.username);
//            }
//        } catch (APIPersistenceException e) {
//            handleExceptionWithCode("Error while accessing persistence layer", e, ExceptionCodes.INTERNAL_ERROR);
//        }
        return response;
    }

    /**
     * Returns the details of all the life-cycle changes done per API or API Product
     *
     * @param      uuid Unique UUID of the API or API Product
     * @return List of lifecycle events per given API or API Product
     * @throws APIManagementException if failed to copy docs
     */
    public List<LifeCycleEvent> getLifeCycleEvents(String uuid) throws APIManagementException {

        return apiDAOImpl.getLifeCycleEvents(uuid);
    }

    @Override
    public void updateSubscription(APIIdentifier apiId, String subStatus, int appId, String organization) throws APIManagementException {

    }

    @Override
    public void updateSubscription(SubscribedAPI subscribedAPI) throws APIManagementException {

    }

    @Override
    public void updateTierPermissions(String tierName, String permissionType, String roles) throws APIManagementException {

    }

    @Override
    public Set getTierPermissions() throws APIManagementException {
        return null;
    }

    @Override
    public Set getThrottleTierPermissions() throws APIManagementException {
        return null;
    }

    @Override
    public boolean publishToExternalAPIStores(API api, List<String> externalStoreIds) throws APIManagementException {
        return false;
    }

    @Override
    public void publishToExternalAPIStores(API api, Set<APIStore> apiStoreSet, boolean apiOlderVersionExist) throws APIManagementException {

    }

    @Override
    public boolean updateAPIsInExternalAPIStores(API api, Set<APIStore> apiStoreSet, boolean apiOlderVersionExist) throws APIManagementException {
        return false;
    }

    @Override
    public Set<APIStore> getExternalAPIStores(String apiId) throws APIManagementException {
        return null;
    }

    @Override
    public Set<APIStore> getPublishedExternalAPIStores(String apiId) throws APIManagementException {
        return null;
    }

    @Override
    public boolean isSynapseGateway() throws APIManagementException {
        return false;
    }

    @Override
    public void saveSwaggerDefinition(API api, String jsonText, String organization) throws APIManagementException {

    }

    @Override
    public void saveSwaggerDefinition(String apiId, String jsonText, String orgId) throws APIManagementException {

    }

    @Override
    public void addAPIProductSwagger(String apiProductId, Map<API, List<APIProductResource>> apiToProductResourceMapping, APIProduct apiProduct, String orgId) throws APIManagementException {

    }

    @Override
    public void updateAPIProductSwagger(String apiProductId, Map<API, List<APIProductResource>> apiToProductResourceMapping, APIProduct apiProduct, String orgId) throws APIManagementException, FaultGatewaysException {

    }

    @Override
    public void validateResourceThrottlingTiers(API api, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void validateResourceThrottlingTiers(String swaggerContent, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void validateAPIThrottlingTier(API api, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void validateProductThrottlingTier(APIProduct apiProduct, String tenantDomain) throws APIManagementException {

    }

    @Override
    public void configureMonetizationInAPIArtifact(API api) throws APIManagementException {

    }

    @Override
    public Monetization getMonetizationImplClass() throws APIManagementException {
        return null;
    }
}
