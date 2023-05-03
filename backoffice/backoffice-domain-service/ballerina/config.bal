//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
import wso2/apk_common_lib as commons;


public type KeyStores record{|
        commons:KeyStore tls;
|};

public type K8sConfiguration record {|
   string host = "kubernetes.default";
   string serviceAccountPath = "/var/run/secrets/kubernetes.io/serviceaccount";
   decimal readTimeout = 5;
|};

public type ManagementServerConfiguration record {|
    string serviceName = "apk-test-wso2-apk-management-server";
    string namespace = "apk";
    string certPath = "/home/wso2apk/devportal/security/truststore/management-server.pem";
|};