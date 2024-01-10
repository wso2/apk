/*
 * Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.enforcer.security.jwt;

import io.opentelemetry.context.Scope;
import org.apache.logging.log4j.ThreadContext;
import org.wso2.apk.enforcer.commons.exception.APISecurityException;
import org.wso2.apk.enforcer.commons.model.AuthenticationContext;
import org.wso2.apk.enforcer.commons.model.RequestContext;
import org.wso2.apk.enforcer.commons.model.ResourceConfig;
import org.wso2.apk.enforcer.constants.APIConstants;
import org.wso2.apk.enforcer.constants.GeneralErrorCodeConstants;
import org.wso2.apk.enforcer.models.API;
import org.wso2.apk.enforcer.security.Authenticator;
import org.wso2.apk.enforcer.subscription.SubscriptionDataHolder;
import org.wso2.apk.enforcer.subscription.SubscriptionDataStore;
import org.wso2.apk.enforcer.tracing.TracingConstants;
import org.wso2.apk.enforcer.tracing.TracingSpan;
import org.wso2.apk.enforcer.tracing.TracingTracer;
import org.wso2.apk.enforcer.tracing.Utils;
import org.wso2.apk.enforcer.util.FilterUtils;

/**
 * Implements the authenticator interface to authenticate non-secured APIs.
 */

public class UnsecuredAPIAuthenticator implements Authenticator {

    @Override
    public boolean canAuthenticate(RequestContext requestContext) {
        // Retrieve the disable security value. If security is disabled for all matching resources,
        // then you can proceed directly with the authentication.
        for (ResourceConfig resourceConfig : requestContext.getMatchedResourcePaths()) {
            if (!resourceConfig.getAuthenticationConfig().isDisabled() || requestContext.getMatchedAPI().isTransportSecurity()) {
                return false;
            }
        }
        return true;
    }

    @Override
    public AuthenticationContext authenticate(RequestContext requestContext) throws APISecurityException {

        TracingSpan unsecuredApiAuthenticatorSpan = null;
        Scope unsecuredApiAuthenticatorSpanScope = null;
        try {
            if (Utils.tracingEnabled()) {
                TracingTracer tracer = Utils.getGlobalTracer();
                unsecuredApiAuthenticatorSpan = Utils
                        .startSpan(TracingConstants.UNSECURED_API_AUTHENTICATOR_SPAN, tracer);
                unsecuredApiAuthenticatorSpanScope = unsecuredApiAuthenticatorSpan.getSpan().makeCurrent();
                Utils.setTag(unsecuredApiAuthenticatorSpan, APIConstants.LOG_TRACE_ID,
                        ThreadContext.get(APIConstants.LOG_TRACE_ID));
            }
            String organization = requestContext.getMatchedAPI().getOrganizationId();
            return FilterUtils.generateAuthenticationContextForUnsecured(requestContext);
        } finally {
            if (Utils.tracingEnabled()) {
                unsecuredApiAuthenticatorSpanScope.close();
                Utils.finishSpan(unsecuredApiAuthenticatorSpan);
            }
        }
    }

    @Override
    public String getChallengeString() {

        return "";
    }

    @Override
    public String getName() {

        return "Unsecured";
    }

    @Override
    public int getPriority() {

        return -20;
    }
}
