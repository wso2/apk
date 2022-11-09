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

import org.wso2.apk.apimgt.api.model.policy.APIPolicy;
import org.wso2.apk.apimgt.api.model.policy.ApplicationPolicy;
import org.wso2.apk.apimgt.api.model.policy.PolicyConstants;
import org.wso2.apk.apimgt.api.model.policy.QuotaPolicy;
import org.wso2.apk.apimgt.api.model.policy.RequestCountLimit;
import org.wso2.apk.apimgt.api.model.policy.SubscriptionPolicy;

import org.wso2.apk.apimgt.impl.APIConstants;

import java.util.UUID;

public class TestObjectCreator {
    private static final String APPLICATION_POLICY_NAME = "100PerMin";
    private static final String SUBSCRIPTION_POLICY_NAME = "Platinum";
    private static final String API_POLICY_NAME = "70PerMin";
    private static final String ORGANIZATION = "carbon.super";

    public static ApplicationPolicy createDefaultApplicationPolicy() {
        ApplicationPolicy applicationPolicy = new ApplicationPolicy(APPLICATION_POLICY_NAME);
        applicationPolicy.setDisplayName("100 Per Min");
        applicationPolicy.setUUID(UUID.randomUUID().toString());
        applicationPolicy.setDescription("Custom Application Policy");
        applicationPolicy.setTenantDomain(ORGANIZATION);
        applicationPolicy.setDeployed(true);
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(100);
        quotaPolicy.setLimit(limit);
        applicationPolicy.setDefaultQuotaPolicy(quotaPolicy);
        return applicationPolicy;
    }

    public static ApplicationPolicy createUpdatedApplicationPolicy() {
        ApplicationPolicy policyToUpdate = createDefaultApplicationPolicy();
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(200);
        quotaPolicy.setLimit(limit);
        policyToUpdate.setDefaultQuotaPolicy(quotaPolicy);
        policyToUpdate.setDescription("Updated Custom Application Policy");
        return policyToUpdate;
    }

    public static SubscriptionPolicy createDefaultSubscriptionPolicy() {
        SubscriptionPolicy subscriptionPolicy = new SubscriptionPolicy(SUBSCRIPTION_POLICY_NAME);
        subscriptionPolicy.setDisplayName("Platinum Policy");
        subscriptionPolicy.setUUID(UUID.randomUUID().toString());
        subscriptionPolicy.setDescription("Custom Subscription Policy");
        subscriptionPolicy.setTenantDomain(ORGANIZATION);
        subscriptionPolicy.setDeployed(true);
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(1000);
        quotaPolicy.setLimit(limit);
        subscriptionPolicy.setDefaultQuotaPolicy(quotaPolicy);
        subscriptionPolicy.setBillingPlan(APIConstants.BILLING_PLAN_FREE);
        return subscriptionPolicy;
    }

    public static SubscriptionPolicy createUpdatedSubscriptionPolicy() {
        SubscriptionPolicy policyToUpdate = createDefaultSubscriptionPolicy();
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(200);
        quotaPolicy.setLimit(limit);
        policyToUpdate.setDefaultQuotaPolicy(quotaPolicy);
        policyToUpdate.setDescription("Updated Custom Subscription Policy");
        return policyToUpdate;
    }

    public static APIPolicy createDefaultAPIPolicy() {
        APIPolicy apiPolicy = new APIPolicy(API_POLICY_NAME);
        apiPolicy.setDisplayName("70 Per Min");
        apiPolicy.setUUID(UUID.randomUUID().toString());
        apiPolicy.setDescription("Custom API Policy");
        apiPolicy.setTenantDomain(ORGANIZATION);
        apiPolicy.setDeployed(true);
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(100);
        quotaPolicy.setLimit(limit);
        apiPolicy.setDefaultQuotaPolicy(quotaPolicy);
        apiPolicy.setUserLevel(PolicyConstants.ACROSS_ALL);
        return apiPolicy;
    }

    public static APIPolicy createUpdatedAPIPolicy() {
        APIPolicy policyToUpdate = createDefaultAPIPolicy();
        QuotaPolicy quotaPolicy = new QuotaPolicy();
        quotaPolicy.setType(PolicyConstants.REQUEST_COUNT_TYPE);
        RequestCountLimit limit = new RequestCountLimit();
        limit.setTimeUnit("min");
        limit.setUnitTime(1);
        limit.setRequestCount(200);
        quotaPolicy.setLimit(limit);
        policyToUpdate.setDefaultQuotaPolicy(quotaPolicy);
        policyToUpdate.setDescription("Updated Custom API Policy");
        return policyToUpdate;
    }
}
