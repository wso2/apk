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

package org.wso2.apk.apimgt.impl.dao;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.ErrorHandler;
import org.wso2.apk.apimgt.api.ExceptionCodes;
import org.wso2.apk.apimgt.api.model.API;
import org.wso2.apk.apimgt.api.model.APICategory;
import org.wso2.apk.apimgt.api.model.APIIdentifier;
import org.wso2.apk.apimgt.api.model.APIInfo;
import org.wso2.apk.apimgt.api.model.APIRevision;
import org.wso2.apk.apimgt.impl.APIConstants;
import org.wso2.apk.apimgt.impl.dao.constants.SQLConstants;
import org.wso2.apk.apimgt.impl.internal.ServiceReferenceHolder;
import org.wso2.apk.apimgt.impl.utils.APIMgtDBUtil;
import org.wso2.apk.apimgt.impl.utils.APIUtil;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * This class represent the ApiMgtDAO.
 */
public class ApiMgtDAO {

    private static final Log log = LogFactory.getLog(ApiMgtDAO.class);
    private static ApiMgtDAO INSTANCE = null;
    private final Object scopeMutex = new Object();
    private boolean forceCaseInsensitiveComparisons = false;
    private boolean multiGroupAppSharingEnabled = false;

    private ApiMgtDAO() {

        String caseSensitiveComparison = ServiceReferenceHolder.getInstance().
                getAPIManagerConfigurationService().getAPIManagerConfiguration()
                .getFirstProperty(APIConstants.API_STORE_FORCE_CI_COMPARISIONS);

        forceCaseInsensitiveComparisons = Boolean.parseBoolean(caseSensitiveComparison);
//        multiGroupAppSharingEnabled = APIUtil.isMultiGroupAppSharingEnabled();
    }

    /**
     * Method to get the instance of the ApiMgtDAO.
     *
     * @return {@link ApiMgtDAO} instance
     */
    public static ApiMgtDAO getInstance() {

        if (INSTANCE == null) {
            INSTANCE = new ApiMgtDAO();
        }

        return INSTANCE;
    }

    /**
     * Get API Identifier by the the API's UUID.
     *
     * @param uuid uuid of the API
     * @return API Identifier
     * @throws APIManagementException if an error occurs
     */
    public APIIdentifier getAPIIdentifierFromUUID(String uuid) throws APIManagementException {

        APIIdentifier identifier = null;
        String sql = SQLConstants.GET_API_IDENTIFIER_BY_UUID_SQL;
        try (Connection connection = APIMgtDBUtil.getConnection()) {
            PreparedStatement prepStmt = connection.prepareStatement(sql);
            prepStmt.setString(1, uuid);
            try (ResultSet resultSet = prepStmt.executeQuery()) {
                while (resultSet.next()) {
                    String provider = resultSet.getString(1);
                    String name = resultSet.getString(2);
                    String version = resultSet.getString(3);
                    identifier = new APIIdentifier(APIUtil.replaceEmailDomain(provider), name, version, uuid);
                }
            }
        } catch (SQLException e) {
            handleExceptionWithCode("Failed to retrieve the API Identifier details for UUID : " + uuid, e,
                    ExceptionCodes.APIMGT_DAO_EXCEPTION);
        }
        return identifier;
    }

    private void handleExceptionWithCode(String msg, Throwable t, ErrorHandler code) throws APIManagementException {
        log.error(msg, t);
        throw new APIManagementException(msg, code);
    }

    /**
     * Get all available API categories of the organization
     * @param organization
     * @return
     * @throws APIManagementException
     */
    public List<APICategory> getAllCategories(String organization) throws APIManagementException {

        List<APICategory> categoriesList = new ArrayList<>();
        try (Connection connection = APIMgtDBUtil.getConnection();
             PreparedStatement statement =
                     connection.prepareStatement(SQLConstants.GET_CATEGORIES_BY_ORGANIZATION_SQL)) {
            statement.setString(1, organization);
            try (ResultSet rs = statement.executeQuery()) {
                while (rs.next()) {
                    String id = rs.getString("UUID");
                    String name = rs.getString("NAME");
                    String description = rs.getString("DESCRIPTION");

                    APICategory category = new APICategory();
                    category.setId(id);
                    category.setName(name);
                    category.setDescription(description);
                    category.setOrganization(organization);

                    categoriesList.add(category);
                }
            }
        } catch (SQLException e) {
            handleExceptionWithCode("Failed to retrieve API categories for organization " + organization,
                    e, ExceptionCodes.APIMGT_DAO_EXCEPTION);
        }
        return categoriesList;
    }


    /**
     * Retrieve basic information about the given API by the UUID quering only from AM_API
     *
     * @param apiId UUID of the API
     * @return basic information about the API
     * @throws APIManagementException error while getting the API information from AM_API
     */
    public APIInfo getAPIInfoByUUID(String apiId) throws APIManagementException {

        try (Connection connection = APIMgtDBUtil.getConnection()) {
            APIRevision apiRevision = getRevisionByRevisionUUID(connection, apiId);
            String sql = SQLConstants.RETRIEVE_API_INFO_FROM_UUID;
            try (PreparedStatement preparedStatement = connection.prepareStatement(sql)) {
                if (apiRevision != null) {
                    preparedStatement.setString(1, apiRevision.getApiUUID());
                } else {
                    preparedStatement.setString(1, apiId);
                }
                try (ResultSet resultSet = preparedStatement.executeQuery()) {
                    if (resultSet.next()) {
                        APIInfo.Builder apiInfoBuilder = new APIInfo.Builder();
                        apiInfoBuilder = apiInfoBuilder.id(resultSet.getString("API_UUID"))
                                .name(resultSet.getString("API_NAME"))
                                .version(resultSet.getString("API_VERSION"))
                                .provider(resultSet.getString("API_PROVIDER"))
                                .context(resultSet.getString("CONTEXT"))
                                .contextTemplate(resultSet.getString("CONTEXT_TEMPLATE"))
                                .status(APIUtil.getApiStatus(resultSet.getString("STATUS")))
                                .apiType(resultSet.getString("API_TYPE"))
                                .createdBy(resultSet.getString("CREATED_BY"))
                                .createdTime(resultSet.getString("CREATED_TIME"))
                                .updatedBy(resultSet.getString("UPDATED_BY"))
                                .updatedTime(resultSet.getString("UPDATED_TIME"))
                                .revisionsCreated(resultSet.getInt("REVISIONS_CREATED"))
                                .organization(resultSet.getString("ORGANIZATION"))
                                .isRevision(apiRevision != null).organization(resultSet.getString("ORGANIZATION"));
                        if (apiRevision != null) {
                            apiInfoBuilder = apiInfoBuilder.apiTier(getAPILevelTier(connection,
                                    apiRevision.getApiUUID(), apiId));
                        } else {
                            apiInfoBuilder = apiInfoBuilder.apiTier(resultSet.getString("API_TIER"));
                        }
                        return apiInfoBuilder.build();
                    }
                }
            }
        } catch (SQLException e) {
            throw new APIManagementException("Error while retrieving apimgt connection", e,
                    ExceptionCodes.INTERNAL_ERROR);
        }
        return null;
    }

    private APIRevision getRevisionByRevisionUUID(Connection connection, String revisionUUID) throws SQLException {

        try (PreparedStatement statement = connection
                .prepareStatement(SQLConstants.APIRevisionSqlConstants.GET_REVISION_BY_REVISION_UUID)) {
            statement.setString(1, revisionUUID);
            try (ResultSet rs = statement.executeQuery()) {
                while (rs.next()) {
                    APIRevision apiRevision = new APIRevision();
                    apiRevision.setId(rs.getInt("ID"));
                    apiRevision.setApiUUID(rs.getString("API_UUID"));
                    apiRevision.setRevisionUUID(rs.getString("REVISION_UUID"));
                    apiRevision.setDescription(rs.getString("DESCRIPTION"));
                    apiRevision.setCreatedTime(rs.getString("CREATED_TIME"));
                    apiRevision.setCreatedBy(rs.getString("CREATED_BY"));
                    return apiRevision;
                }
            }
        }
        return null;
    }

    private String getAPILevelTier(Connection connection, String apiUUID, String revisionUUID) throws SQLException {

        try (PreparedStatement preparedStatement =
                     connection.prepareStatement(SQLConstants.GET_REVISIONED_API_TIER_SQL)) {
            preparedStatement.setString(1, apiUUID);
            preparedStatement.setString(2, revisionUUID);
            try (ResultSet resultSet = preparedStatement.executeQuery()) {
                if (resultSet.next()) {
                    return resultSet.getString("API_TIER");
                }
            }
        }
        return null;
    }

    public void setDefaultVersion(API api) throws APIManagementException {

        APIIdentifier apiId = api.getId();
        try (Connection connection = APIMgtDBUtil.getConnection()) {
            try (PreparedStatement preparedStatement =
                         connection.prepareStatement(SQLConstants.RETRIEVE_DEFAULT_VERSION)) {
                preparedStatement.setString(1, apiId.getApiName());
                preparedStatement.setString(2, APIUtil.replaceEmailDomainBack(apiId.getProviderName()));
                try (ResultSet resultSet = preparedStatement.executeQuery()) {
                    if (resultSet.next()) {
                        api.setDefaultVersion(apiId.getVersion().equals(resultSet.getString("DEFAULT_API_VERSION")));
                        api.setAsPublishedDefaultVersion(apiId.getVersion().equals(resultSet.getString(
                                "PUBLISHED_DEFAULT_API_VERSION")));
                    }
                }
            }

        } catch (SQLException e) {
            throw new APIManagementException("Error while retrieving apimgt connection", e,
                    ExceptionCodes.INTERNAL_ERROR);
        }
    }

    /**
     * Return the existing versions for the given api name for the provider
     *
     * @param apiName     api name
     * @param apiProvider provider
     * @param organization identifier of the organization
     * @return set version
     * @throws APIManagementException
     */
    public Set<String> getAPIVersions(String apiName, String apiProvider, String organization) throws APIManagementException {

        Set<String> versions = new HashSet<String>();

        try (Connection connection = APIMgtDBUtil.getConnection();
             PreparedStatement statement = connection.prepareStatement(SQLConstants.GET_API_VERSIONS)) {
            statement.setString(1, APIUtil.replaceEmailDomainBack(apiProvider));
            statement.setString(2, apiName);
            statement.setString(3, organization);
            ResultSet resultSet = statement.executeQuery();
            while (resultSet.next()) {
                versions.add(resultSet.getString("API_VERSION"));
            }
        } catch (SQLException e) {
            handleExceptionWithCode("Error while retrieving versions for api " + apiName + " for the provider "
                    + apiProvider, e, ExceptionCodes.APIMGT_DAO_EXCEPTION);
        }
        return versions;
    }
}
