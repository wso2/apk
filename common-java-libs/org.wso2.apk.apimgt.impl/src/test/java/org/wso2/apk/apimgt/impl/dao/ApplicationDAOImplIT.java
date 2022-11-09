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

import org.testng.Assert;
import org.testng.annotations.BeforeClass;
import org.testng.annotations.Test;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.Application;
import org.wso2.apk.apimgt.api.model.Subscriber;
import org.wso2.apk.apimgt.api.model.policy.PolicyConstants;
import org.wso2.apk.apimgt.impl.APIConstants;
import org.wso2.apk.apimgt.impl.dao.constants.SQLConstants;
import org.wso2.apk.apimgt.impl.dao.impl.ApplicationDAOImpl;
import org.wso2.apk.apimgt.impl.utils.APIMgtDBUtil;
import org.wso2.apk.apimgt.impl.utils.APIUtil;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.Timestamp;
import java.time.Instant;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Date;

import java.util.UUID;

public class ApplicationDAOImplIT extends DAOIntegrationTestBase {
    private static final String ORGANIZATION = "carbon.super";
    private static final String USERNAME = "developer@carbon.super";
    private static String applicationUUID;
    private static List<Application> applicationList = new ArrayList<>();

    @BeforeClass
    public void init() throws Exception {
        Subscriber subscriber = getSubscriber();
        addSubscriber(subscriber);
    }

   /*-----------------------------------------------------------------------------------------------------------------
   |  Test Scenario:  Add get update and delete an application
   |  Traceability: https://github.com/wso2/apk/issues/63
   |  Pre-condition: N/A
   |  Post-condition: N/A
   |  Dependencies: N/A
   |  Assertions:
   |         Verify applicationId of the fetched application equals to the ID returned in add application call
   |         Verify application name of the Application name fetched by ID is equal to that of the original application
   |         Verify application fetched by ID is not null
   |         Verify application name of the Application fetched by ID is equal to that of the original application
   |         Verify application fetched by name is not null
   |         Verify application name of the Application fetched by name is equal to that of the original application
   |         Verify application fetched by uuid is not null
   |         Verify application name of the Application fetched by uuid is equal to that of the original application
   |         Verify application tier of the Application fetched by uuid is equal to that of the original application
   |         Verify application tier is updated upon application update
   |         Verify application fetched by uuid is null after deleting the application
   ------------------------------------------------------------------------------------------------------------------*/
    @Test(description = "Add get update and delete an application")
    public void testAddGetUpdateDeleteApplication() throws Exception {
        ApplicationDAO applicationDAO = ApplicationDAOImpl.getInstance();
        Application application = getApplication();
        int applicationId = applicationDAO.addApplication(application, USERNAME, ORGANIZATION);

        int fetchedAppId = applicationDAO.getApplicationId(application.getName(), USERNAME);
        Assert.assertEquals(fetchedAppId, applicationId);

        String applicationName = applicationDAO.getApplicationNameFromId(applicationId);
        Assert.assertEquals(applicationName, application.getName());

        Application fetchedApplication = applicationDAO.getApplicationById(applicationId);
        Assert.assertNotNull(fetchedApplication);
        Assert.assertEquals(fetchedApplication.getName(), application.getName());

        fetchedApplication = applicationDAO.getApplicationByName(application.getName(), USERNAME, null);
        Assert.assertNotNull(fetchedApplication);
        Assert.assertEquals(fetchedApplication.getName(), application.getName());

        fetchedApplication = applicationDAO.getApplicationByUUID(applicationUUID);
        Assert.assertNotNull(fetchedApplication);
        Assert.assertEquals(fetchedApplication.getName(), application.getName());
        Assert.assertEquals(fetchedApplication.getTier(), APIConstants.UNLIMITED_TIER);
        fetchedApplication.setTier(APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN);
        applicationDAO.updateApplication(fetchedApplication);

        fetchedApplication = applicationDAO.getApplicationByUUID(applicationUUID);
        Assert.assertEquals(fetchedApplication.getTier(), APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN);

        applicationDAO.deleteApplication(fetchedApplication);
        fetchedApplication = applicationDAO.getApplicationByUUID(applicationUUID);
        Assert.assertNull(fetchedApplication);
    }

    /*-----------------------------------------------------------------------------------------------------------------
   |  Test Scenario:  Get Applications count
   |  Traceability: https://github.com/wso2/apk/issues/63
   |  Pre-condition: N/A
   |  Post-condition: N/A
   |  Dependencies: N/A
   |  Assertions:
   |         Verify application count in DB
   ------------------------------------------------------------------------------------------------------------------*/
    @Test(description = "Get applications")
    public void testGetApplicationsCount() throws Exception {
        ApplicationDAO applicationDAO = ApplicationDAOImpl.getInstance();
        int initialCount = applicationDAO.getApplicationsCount(ORGANIZATION, USERNAME, "App");
        List<Application> applications = getApplications();
        for (Application application : applications) {
            applicationDAO.addApplication(application, USERNAME, ORGANIZATION);
            applicationList.add(application);
        }

        int count = applicationDAO.getApplicationsCount(ORGANIZATION, USERNAME, "App");
        Assert.assertEquals(count, initialCount + 10);
    }

   /*-----------------------------------------------------------------------------------------------------------------
   |  Test Scenario:  Update application owner
   |  Traceability: https://github.com/wso2/apk/issues/63
   |  Pre-condition: N/A
   |  Post-condition: N/A
   |  Dependencies: N/A
   |  Assertions:
   |         Verify application Id returned from add application call is not null
   |         Verify application owner of the fetched application is equal to the original owner before update
   |         Verify application owner of the fetched application is equal to the new owner after update
   ------------------------------------------------------------------------------------------------------------------*/
    @Test(description = "Update application owner")
    public void testUpdateApplicationOwner() throws Exception {
        ApplicationDAO applicationDAO = ApplicationDAOImpl.getInstance();
        String newOwner = "deveoper2@carbon.super";

        Subscriber subscriber = getSubscriber(newOwner, "Subscriber User 2", "developer@gmail.com");
        addSubscriber(subscriber);

        Application application = getApplication("ChangeOwnerApp", "Test Change Owner App", UUID.randomUUID().toString(), APIConstants.UNLIMITED_TIER);
        int appId = applicationDAO.addApplication(application, getSubscriber().getName(), ORGANIZATION);
        Assert.assertNotEquals(appId, 0);

        Application fetchedApplication = applicationDAO.getApplicationById(appId);
        Assert.assertEquals(fetchedApplication.getOwner(), USERNAME);
        applicationDAO.updateApplicationOwner(newOwner, application);
        fetchedApplication = applicationDAO.getApplicationById(appId);
        Assert.assertEquals(fetchedApplication.getOwner(), newOwner);
    }

    /*-----------------------------------------------------------------------------------------------------------------
    |  Test Scenario:  Update application owner
    |  Traceability: https://github.com/wso2/apk/issues/63
    |  Pre-condition: N/A
    |  Post-condition: N/A
    |  Dependencies: N/A
    |  Assertions:
    |         Verify application Id returned from add application call is not null
    |         Verify application owner of the fetched application is equal to the original owner before update
    |         Verify application owner of the fetched application is equal to the new owner after update
    ------------------------------------------------------------------------------------------------------------------*/
    @Test(description = "Add and delete application attributes")
    public void testAddDeleteApplicationAttributes() throws Exception {
        ApplicationDAO applicationDAO = ApplicationDAOImpl.getInstance();
        String uuid = UUID.randomUUID().toString();
        Application application = getApplication("AttributesApp", "Sample Application to test Application Attributes",
                uuid, APIConstants.UNLIMITED_TIER);
        int appId = applicationDAO.addApplication(application, getSubscriber().getName(), ORGANIZATION);
        applicationList.add(application);
        Assert.assertNotEquals(appId, 0);

        Map<String, String> applicationAttributes = new HashMap<>();
        applicationAttributes.put("appAttribute-1", "value-1");
        applicationAttributes.put("appAttribute-2", "value-2");
        applicationDAO.addApplicationAttributes(applicationAttributes, appId, ORGANIZATION);

        Application fetchedApplication = applicationDAO.getApplicationById(appId);
        Map<String, String> attributesFromApp = fetchedApplication.getApplicationAttributes();
        Assert.assertNotNull(attributesFromApp);
        Assert.assertEquals(attributesFromApp.size(), 2);

        applicationDAO.deleteApplicationAttributes("appAttribute-1", appId);
        applicationDAO.deleteApplicationAttributes("appAttribute-2", appId);
        fetchedApplication = applicationDAO.getApplicationById(appId);
        attributesFromApp = fetchedApplication.getApplicationAttributes();
        Assert.assertNotNull(attributesFromApp);
        Assert.assertEquals(attributesFromApp.size(), 0);
    }



    private Subscriber getSubscriber(String username, String description, String email) {
        Subscriber subscriber = new Subscriber(username);
        subscriber.setDescription(description);
        subscriber.setOrganization(ORGANIZATION);
        subscriber.setEmail(email);
        subscriber.setSubscribedDate(Date.from(Instant.now()));
        return subscriber;
    }

    private Application getApplication() {
        applicationUUID = UUID.randomUUID().toString();
        return getApplication("TestApp", "Test Application", applicationUUID, APIConstants.UNLIMITED_TIER);
    }

    private Subscriber getSubscriber() {
        return getSubscriber(USERNAME, "Subscriber User", "developer@gmail.com");
    }

    private Application getApplication(String name, String description, String uuid, String tier) {
        Subscriber subscriber = getSubscriber();
        Application application = new Application(name, subscriber);
        application.setDescription(description);
        application.setUUID(uuid);
        application.setOrganization(ORGANIZATION);
        application.setOwner(USERNAME);
        application.setTier(tier);
        return application;
    }

    private List<Application> getApplications() {
        List<Application> applications = new ArrayList<>();
        applications.add(getApplication("App1", "Test Application 1", UUID.randomUUID().toString(),
                APIConstants.UNLIMITED_TIER));
        applications.add(getApplication("App2", "Test Application 2", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN));
        applications.add(getApplication("App3", "Test Application 3", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_TEN_REQ_PER_MIN));
        applications.add(getApplication("App4", "Test Application 4", UUID.randomUUID().toString(),
                APIConstants.UNLIMITED_TIER));
        applications.add(getApplication("App5", "Test Application 5", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN));
        applications.add(getApplication("App6", "Test Application 6", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_TEN_REQ_PER_MIN));
        applications.add(getApplication("App7", "Test Application 7", UUID.randomUUID().toString(),
                APIConstants.UNLIMITED_TIER));
        applications.add(getApplication("App8", "Test Application 8", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN));
        applications.add(getApplication("App9", "Test Application 9", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_FIFTY_REQ_PER_MIN));
        applications.add(getApplication("App10", "Test Application 10", UUID.randomUUID().toString(),
                APIConstants.DEFAULT_APP_POLICY_TEN_REQ_PER_MIN));
        return applications;
    }

    private void addSubscriber(Subscriber subscriber) throws APIManagementException {
        try (Connection conn = APIMgtDBUtil.getConnection()){
            conn.setAutoCommit(false);

            String query = SQLConstants.ADD_SUBSCRIBER_SQL;
            try (PreparedStatement ps = conn.prepareStatement(query)) {
                ps.setString(1, subscriber.getName());
                ps.setString(2, subscriber.getOrganization());
                ps.setString(3, subscriber.getEmail());

                Timestamp timestamp = new Timestamp(subscriber.getSubscribedDate().getTime());
                ps.setTimestamp(4, timestamp);
                ps.setString(5, subscriber.getName());
                ps.setTimestamp(6, timestamp);
                ps.setTimestamp(7, timestamp);
                ps.executeUpdate();
                conn.commit();
            }
        } catch (SQLException e) {
            APIUtil.handleException("Error in adding new subscriber: " + e.getMessage(), e);
        }
    }
}
