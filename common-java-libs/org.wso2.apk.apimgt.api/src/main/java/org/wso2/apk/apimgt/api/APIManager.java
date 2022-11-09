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

import org.wso2.apk.apimgt.api.model.*;
import org.wso2.apk.apimgt.api.model.graphql.queryanalysis.GraphqlComplexityInfo;
import org.wso2.apk.apimgt.api.model.policy.Policy;

import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * Core API management interface which provides functionality related to APIs, API metadata
 * and API subscribers (consumers).
 */
public interface APIManager {

    /**
     * This method is to add a comment.
     *
     * @param uuid Api uuid
     * @param comment    comment object
     * @param user       Username of the comment author
     * @throws APIManagementException if failed to add comment for API
     */
    String addComment(String uuid, Comment comment, String user) throws APIManagementException;

    /**
     * This method is to get a comment of an API.
     *
     * @param apiTypeWrapper Api Type Wrapper
     * @param commentId      Comment ID
     * @param replyLimit
     * @param replyOffset
     * @return Comment
     * @throws APIManagementException if failed to get comments for identifier
     */
    Comment getComment(ApiTypeWrapper apiTypeWrapper, String commentId, Integer replyLimit, Integer replyOffset) throws
            APIManagementException;

    /**
     * @param apiTypeWrapper  Api type wrapper
     * @param parentCommentID
     * @param replyLimit
     * @param replyOffset
     * @return Comments
     * @throws APIManagementException if failed to get comments for identifier
     */
    CommentList getComments(ApiTypeWrapper apiTypeWrapper, String parentCommentID, Integer replyLimit, Integer replyOffset) throws APIManagementException;

    /**
     * @param apiTypeWrapper Api Type Wrapper
     * @param commentId      comment ID
     * @param comment        Comment object
     * @return Comments
     * @throws APIManagementException if failed to get comments for identifier
     */
    boolean editComment(ApiTypeWrapper apiTypeWrapper, String commentId, Comment comment) throws APIManagementException;

    /**
     * This method is to delete a comment.
     *
     * @param apiTypeWrapper API Type Wrapper
     * @param commentId      Comment ID
     * @return boolean
     * @throws APIManagementException if failed to delete comment for identifier
     */
    boolean deleteComment(ApiTypeWrapper apiTypeWrapper, String commentId) throws APIManagementException;

    /**
     * Returns the minimalistic information about the API given the UUID. This will only query from AM database AM_API
     * table.
     *
     * @param id UUID of the API
     * @return basic information about the API
     * @throws APIManagementException error while getting the API information from AM_API
     */
    APIInfo getAPIInfoByUUID(String id) throws APIManagementException;

    /**
     * Get API or APIProduct by registry artifact id
     *
     * @param uuid   Registry artifact id
     * @param organization  Organization
     * @return ApiTypeWrapper wrapping the API or APIProduct of the provided artifact id
     * @throws APIManagementException
     */
    ApiTypeWrapper getAPIorAPIProductByUUID(String uuid, String organization) throws APIManagementException;

    /**
     * @param searchQuery search query. ex : provider:admin
     * @param organization Identifier of an organization
     * @param start starting number
     * @param end ending number
     * @return
     * @throws APIManagementException
     */
    Map<String, Object> searchPaginatedAPIs(String searchQuery, String organization, int start, int end,
                                            String sortBy, String sortOrder) throws APIManagementException;


    /**
     * Returns a list of documentation attached to a particular API
     *
     * @param uuid id of the api
     * @param organization  Identifier of an organization
     * @return List<Documentation>
     * @throws APIManagementException if failed to get Documentations
     */
    List<Documentation> getAllDocumentation(String uuid, String organization) throws APIManagementException;

    /**
     * Get a documentation by artifact Id
     *
     * @param apiId         apiId
     * @param docId         DocumentID
     * @param organization  Identifier of the organization
     * @return Documentation
     * @throws APIManagementException if failed to get Documentation
     */
    Documentation getDocumentation(String apiId, String docId, String organization)
            throws APIManagementException;

    /**
     * Get a documentation Content by apiid and doc id
     *
     * @param apiId         ID of the API
     * @param docId         DocumentID
     * @param organization  Identifier of an organization
     * @return DocumentationContent
     * @throws APIManagementException if failed to get Documentation
     */
    DocumentationContent getDocumentationContent(String apiId, String docId, String organization)
            throws APIManagementException;

    /**
     * Retrieves the icon image associated with a particular API as a stream.
     *
     * @param apiId ID representing the API
     * @param orgId  Identifier of an organization
     * @return an Icon containing image content and content type information
     * @throws APIManagementException if an error occurs while retrieving the image
     */
    ResourceFile getIcon(String apiId, String orgId) throws APIManagementException;

    /**
     * Get minimal details of API by registry artifact id
     *
     * @param uuid          Registry artifact id
     * @param organization  Identifier of an organization
     * @return API of the provided artifact id
     * @throws APIManagementException
     */
    API getLightweightAPIByUUID(String uuid, String organization) throws APIManagementException;

    /**
     * Returns the OpenAPI definition as a string
     *
     * @param apiId        ID of the API
     * @param organization Identifier of an organization
     * @return swagger string
     * @throws APIManagementException
     */
    String getOpenAPIDefinition(String apiId, String organization) throws APIManagementException;

    /**
     * Returns a set of API versions for the given provider and API name
     *
     * @param providerName name of the provider (common)
     * @param apiName      name of the api
     * @param organization organization
     * @return Set of version strings (possibly empty)
     * @throws APIManagementException if failed to get version for api
     */
    Set<String> getAPIVersions(String providerName, String apiName, String organization) throws APIManagementException;
}
