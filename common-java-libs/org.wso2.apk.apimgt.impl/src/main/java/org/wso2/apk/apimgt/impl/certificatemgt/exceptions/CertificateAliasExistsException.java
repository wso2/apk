/*
 * Copyright (c) 2017, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * you may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.wso2.apk.apimgt.impl.certificatemgt.exceptions;

/**
 * This represents custom exception class for certificate management in api manager, which will be thrown in
 * scenarios that certificate alias exists in the data base.
 */
public class CertificateAliasExistsException extends Throwable {

    public CertificateAliasExistsException(String message) {
        super(message);
    }

    public CertificateAliasExistsException(String message, Throwable e) {
        super(message, e);
    }

    public CertificateAliasExistsException(Throwable e) {
        super(e);
    }
}
