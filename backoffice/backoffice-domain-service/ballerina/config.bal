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

type DatasourceConfiguration record {
    string name = "jdbc/apkdb";
    string description;
    string url;
    string username;
    string password;
    string host;
    int port;
    string databaseName;
    int maxPoolSize = 50;
    int minIdle = 20;
    int maxLifeTime = 60000;
    int validationTimeout;
    boolean autoCommit = true;
    string testQuery;
    string driver;
};

type APKConfiguration record {
    DatasourceConfiguration datasourceConfiguration;
};
public type KeyStores record{|
        commons:KeyStore tls;
|};