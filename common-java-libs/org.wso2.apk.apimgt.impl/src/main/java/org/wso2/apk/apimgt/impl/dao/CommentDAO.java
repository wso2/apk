/*
 * Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.apimgt.impl.dao;

import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.APIIdentifier;
import org.wso2.apk.apimgt.api.model.ApiTypeWrapper;
import org.wso2.apk.apimgt.api.model.Comment;
import org.wso2.apk.apimgt.api.model.CommentList;

public interface CommentDAO {

    /**
     * Adds a comment for an API
     *
     * @param uuid API uuid
     * @param comment Commented Text
     * @param user User who did the comment
     * @return Comment ID
     */
    String addComment(String uuid, Comment comment, String user) throws APIManagementException;

    /**
     * Returns a specific comment of an API
     *
     * @param commentId  Comment ID
     * @param apiTypeWrapper Api Type Wrapper
     * @return Comment Array
     * @throws APIManagementException
     */
    Comment getComment(ApiTypeWrapper apiTypeWrapper, String commentId, Integer replyLimit,
                       Integer replyOffset) throws APIManagementException;

    /**
     * Returns all the Comments on an API
     *
     * @param apiTypeWrapper API type Wrapper
     * @param parentCommentID Parent Comment ID
     * @return Comment Array
     * @throws APIManagementException
     */
    CommentList getComments(ApiTypeWrapper apiTypeWrapper, String parentCommentID, Integer limit,
                            Integer offset) throws APIManagementException;

    /**
     * Edit a comment
     *
     * @param apiTypeWrapper API Type Wrapper
     * @param commentId Comment ID
     * @param comment Comment object
     * @throws APIManagementException
     */
    boolean editComment(ApiTypeWrapper apiTypeWrapper, String commentId, Comment comment) throws
            APIManagementException;

    /**
     * Delete a comment
     *
     * @param apiTypeWrapper API Type Wrapper
     * @param commentId Comment ID
     * @throws APIManagementException
     */
    boolean deleteComment(ApiTypeWrapper apiTypeWrapper, String commentId) throws APIManagementException;

    /**
     * Adds a comment for an API
     *
     * @param identifier  API Identifier
     * @param commentText Commented Text
     * @param user        User who did the comment
     * @return Comment ID
     * @deprecated This method needs to be removed once the Jaggery web apps are removed.
     */
    int addComment(APIIdentifier identifier, String commentText, String user)  throws APIManagementException;

    /**
     * Returns all the Comments on an API
     *
     * @param uuid      API uuid
     * @param parentCommentID Parent Comment ID
     * @return Comment Array
     * @throws APIManagementException
     */
    Comment[] getComments(String uuid, String parentCommentID) throws APIManagementException;

    /**
     * Delete a comment
     *
     * @param uuid API uuid
     * @param commentId  Comment ID
     * @throws APIManagementException
     */
    void deleteComment(String uuid, String commentId) throws APIManagementException;
}
