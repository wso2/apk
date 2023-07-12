#!/usr/bin/env bash

# Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.

set -e
#This is a sample script to build APK. When you build project please run this script in project root level.
#All relative paths etc designed from root directory. Users can customize this as per demand. Ex: If you wish to 
#build and run runtime domain service then can build it alone and do deployment.

# mkdir -p ballerina-dist
# wget 'https://dist.ballerina.io/downloads/2201.5.0/ballerina-2201.5.0-swan-lake-linux-x64.deb' -P ballerina-dist
# sudo dpkg -i ballerina-dist/ballerina-2201.5.0-swan-lake-linux-x64.deb

current_dir=$PWD;
# cd $current_dir;
# cd ../../common-bal-libs/apk-common-lib;./gradlew build;
cd $current_dir;
cd ../../adapter;./gradlew build -Pversion='test';
cd $current_dir;
cd ../../gateway/enforcer;./gradlew build -Pversion='test';
