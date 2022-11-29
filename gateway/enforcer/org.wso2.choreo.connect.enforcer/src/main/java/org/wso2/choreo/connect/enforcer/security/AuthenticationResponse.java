/*
 * Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.wso2.choreo.connect.enforcer.security;

/**
 * AuthenticationResponse object contains the authentication status of the request when it passed
 * through an authenticator.
 */
public class AuthenticationResponse {
    private boolean authenticated;
    private boolean mandatoryAuthentication;
    private boolean continueToNextAuthenticator;

    public AuthenticationResponse(boolean authenticated, boolean mandatoryAuthentication,
                                  boolean continueToNextAuthenticator) {
        this.authenticated = authenticated;
        this.mandatoryAuthentication = mandatoryAuthentication;
        this.continueToNextAuthenticator = continueToNextAuthenticator;
    }

    public boolean isAuthenticated() {
        return authenticated;
    }

    public boolean isMandatoryAuthentication() {
        return mandatoryAuthentication;
    }

    public boolean isContinueToNextAuthenticator() {
        return continueToNextAuthenticator;
    }
}
