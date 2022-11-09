/*
 *  Copyright (c) 2005-2011, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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
package org.wso2.apk.apimgt.api;

import org.wso2.apk.apimgt.api.model.API;
import org.wso2.apk.apimgt.api.model.APIIdentifier;
import org.wso2.apk.apimgt.api.model.APIRevisionDeployment;
import org.wso2.apk.apimgt.api.model.APIStateChangeResponse;
import org.wso2.apk.apimgt.api.model.ApiTypeWrapper;
import org.wso2.apk.apimgt.api.model.Documentation;
import org.wso2.apk.apimgt.api.model.DocumentationContent;
import org.wso2.apk.apimgt.api.model.LifeCycleEvent;
import org.wso2.apk.apimgt.api.model.ResourceFile;
import org.wso2.apk.apimgt.api.model.ResourcePath;

import java.util.List;
import java.util.Map;

/**
 * APIProvider responsible for providing helper functionality
 */
public interface APIProvider extends APIManager {

    /**
     * Returns details of an API
     *
     * @param uuid         UUID of the API's registry artifact
     * @param organization Identifier of an organization
     * @return An API object related to the given artifact id or null
     * @throws APIManagementException if failed get API from APIIdentifier
     */
    API getAPIbyUUID(String uuid, String organization) throws APIManagementException;

    /**
     * Checks whether the given document already exists for the given api/product
     *
     * @param uuid         API/Product id
     * @param docName      Name of the document
     * @param organization Identifier of the organization
     * @return true if document already exists for the given api/product
     * @throws APIManagementException if failed to check existence of the documentation
     */
    boolean isDocumentationExist(String uuid, String docName, String organization) throws APIManagementException;

    /**
     * Updates a given documentation
     *
     * @param apiId         id of the document
     * @param documentation Documentation
     * @param organization  Identifier of an organization
     * @return updated documentation Documentation
     * @throws APIManagementException if failed to update docs
     */
    Documentation updateDocumentation(String apiId, Documentation documentation, String organization)
            throws APIManagementException;

    /**
     * Removes a given documentation
     *
     * @param apiId        api uuid
     * @param documentId   ID of the documentation
     * @param organization Identifier of an organization
     * @throws APIManagementException if failed to remove documentation
     */
    void removeDocumentation(String apiId, String documentId, String organization) throws APIManagementException;

    /**
     * Adds Documentation to an API/Product
     *
     * @param uuid                API/Product Identifier
     * @param documentation       Documentation
     * @param organization        Identifier of an organization
     * @return Documentation      created documentation Documentation
     * @throws APIManagementException if failed to add documentation
     */
    Documentation addDocumentation(String uuid, Documentation documentation, String organization) throws APIManagementException;

    /**
     * Adds Document content to an API/Product
     *
     * @param uuid          API/Product Identifier
     * @param content       Documentation content
     * @param docId         doc uuid
     * @param organization  Identifier of an organization
     * @throws APIManagementException if failed to add documentation
     */
    void addDocumentationContent(String uuid, String docId, String organization, DocumentationContent content)
            throws APIManagementException;

    List<ResourcePath> getResourcePathsOfAPI(APIIdentifier apiId) throws APIManagementException;

    /**
     * Get an API Revisions Deployment mapping details by providing revision uuid
     *
     * @param revisionUUID Revision UUID
     * @return List<APIRevisionDeployment> Object
     * @throws APIManagementException if failed to get the related API revision Deployment Mapping details
     */
    List<APIRevisionDeployment> getAPIRevisionDeploymentList(String revisionUUID) throws APIManagementException;

    /**
     * Updates design and implementation of an existing API. This method must not be used to change API status. Implementations
     * should throw an exceptions when such attempts are made. All life cycle state changes
     * should be carried out using the changeAPIStatus method of this interface.
     *
     * @param api API
     * @param existingAPI existing api
     * @throws APIManagementException if failed to update API
     * @throws FaultGatewaysException on Gateway Failure
     * @return updated API
     */
    API updateAPI(API api, API existingAPI) throws APIManagementException, FaultGatewaysException;

    /**
     * Add or update thumbnail image of an api
     * @param apiId    ID of the API
     * @param resource Image resource
     * @param orgId    Identifier of an organization
     * @throws APIManagementException
     */
    void setThumbnailToAPI(String apiId, ResourceFile resource, String orgId) throws APIManagementException;

    /**
     * This method is to change registry lifecycle states for an API artifact
     *
     * @param orgId UUID of the organization
     * @param  apiTypeWrapper API Type Wrapper
     * @param  action  Action which need to execute from registry lifecycle
     * @param  checklist checklist items
     * @return APIStateChangeResponse API workflow state and WorkflowResponse
     * */
    APIStateChangeResponse changeLifeCycleStatus(String orgId, ApiTypeWrapper apiTypeWrapper, String action,
                                                 Map<String, Boolean> checklist) throws APIManagementException;

    /**
     * This method returns the lifecycle data for an API including current state,next states.
     *
     * @param apiId id of the api
     * @param orgId  Identifier of an organization
     * @return Map<String,Object> a map with lifecycle data
     */
    Map<String, Object> getAPILifeCycleData(String apiId, String orgId) throws APIManagementException;

    /**
     * Returns the details of all the life-cycle changes done per API or API Product
     *
     * @param uuid     Unique UUID of the API or API Product
     * @return List of life-cycle events per given API or API Product
     * @throws APIManagementException if failed to copy docs
     */
    List<LifeCycleEvent> getLifeCycleEvents(String uuid) throws APIManagementException;
}
