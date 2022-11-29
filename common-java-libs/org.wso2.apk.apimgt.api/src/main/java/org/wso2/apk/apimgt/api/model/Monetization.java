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

package org.wso2.apk.apimgt.api.model;

import org.wso2.apk.apimgt.api.model.policy.SubscriptionPolicy;
import org.wso2.apk.apimgt.api.MonetizationException;

/**
 * Interface for monetization
 */

public interface Monetization {

    /**
     * Create billing plan for a policy
     *
     * @param subPolicy subscription policy
     * @return true if successful, false otherwise
     * @throws MonetizationException if the action failed
     */
    boolean createBillingPlan(SubscriptionPolicy subPolicy) throws MonetizationException;

    /**
     * Update billing plan of a policy
     *
     * @param subPolicy subscription policy
     * @return true if successful, false otherwise
     * @throws MonetizationException if the action failed
     */
    boolean updateBillingPlan(SubscriptionPolicy subPolicy) throws MonetizationException;

    /**
     * Delete a billing plan of a policy
     *
     * @param subPolicy subscription policy
     * @return true if successful, false otherwise
     * @throws MonetizationException if the action failed
     */
    boolean deleteBillingPlan(SubscriptionPolicy subPolicy) throws MonetizationException;

    /**
     * Publish the usage for a subscription to the billing engine
     *
     * @return true if the job is successfull, and false otherwise
     * @throws MonetizationException if failed to get current usage for a subscription
     */
    boolean publishMonetizationUsageRecords(MonetizationUsagePublishInfo monetizationUsagePublishInfo)
            throws MonetizationException;

}
