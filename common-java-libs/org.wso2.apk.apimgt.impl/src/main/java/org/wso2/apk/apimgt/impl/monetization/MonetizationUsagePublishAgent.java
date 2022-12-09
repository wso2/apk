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

package org.wso2.apk.apimgt.impl.monetization;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.wso2.apk.apimgt.api.APIAdmin;
import org.wso2.apk.apimgt.api.APIManagementException;
import org.wso2.apk.apimgt.api.model.Monetization;
import org.wso2.apk.apimgt.api.model.MonetizationUsagePublishInfo;
import org.wso2.apk.apimgt.impl.APIConstants;

import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.TimeZone;

public class MonetizationUsagePublishAgent implements Runnable {

    private static final Log log = LogFactory.getLog(MonetizationUsagePublishAgent.class);
    MonetizationUsagePublishInfo monetizationUsagePublishInfo;

    public MonetizationUsagePublishAgent(MonetizationUsagePublishInfo publishInfo) {
        monetizationUsagePublishInfo = publishInfo;
    }

    @Override
    public void run() {
    }
}
