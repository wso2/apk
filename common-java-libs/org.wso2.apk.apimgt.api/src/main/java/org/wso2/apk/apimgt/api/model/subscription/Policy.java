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

package org.wso2.apk.apimgt.api.model.subscription;

import org.wso2.apk.apimgt.api.model.policy.PolicyConstants;
import org.wso2.apk.apimgt.api.model.policy.QuotaPolicy;

/**
 * Top level entity for representing a Throttling Policy.
 */
public class Policy implements CacheableEntity<String> {

    public enum POLICY_TYPE {
        SUBSCRIPTION,
        APPLICATION
    }
    private int id = -1;
    private int tenantId = -1;
    private String name = null;
    private String quotaType = null;
    private QuotaPolicy quotaPolicy;
    private String tenantDomain;

    public int getId() {

        return id;
    }

    public void setId(int id) {

        this.id = id;
    }

    public String getQuotaType() {

        return quotaType;
    }

    public void setQuotaType(String quotaType) {

        this.quotaType = quotaType;
    }

    public boolean isContentAware() {

        return PolicyConstants.BANDWIDTH_TYPE.equals(quotaType);
    }

    public int getTenantId() {

        return tenantId;
    }

    public void setTenantId(int tenantId) {

        this.tenantId = tenantId;
    }

    public String getName() {

        return name;
    }

    public void setName(String name) {

        this.name = name;
    }

    @Override
    public String getCacheKey() {

        return getPolicyCacheKey(getName(), getTenantId());
    }

    public static String getPolicyCacheKey(String tierName, int tenantId) {

        return tierName + DELEM_PERIOD + tenantId;
    }

    public QuotaPolicy getQuotaPolicy() {
        return quotaPolicy;
    }

    public void setQuotaPolicy(QuotaPolicy quotaPolicy) {
        this.quotaPolicy = quotaPolicy;
    }

    public String getTenantDomain() {
        return tenantDomain;
    }

    public void setTenantDomain(String tenantDomain) {
        this.tenantDomain = tenantDomain;
    }
}
