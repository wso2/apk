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

import ballerina/test;


@test:Config {dependsOn: [createAPITest]}
function changeLCStateTest() {
    LifecycleState|error lcState1 = changeLifeCyleState("PUBLISHED", "01ed75e2-b30b-18c8-wwf2-25da7edd2231", "carbon.super");
    if lcState1 is LifecycleState {
        test:assertTrue(true, "Successfully change the LC state");
    } else {
        test:assertFail("Error occured while changing LC state");
    }
}


@test:Config {dependsOn: [createAPITest, changeLCStateTest]}
function getLcEventHistoryTest() {
    LifecycleHistory|error? lcState = getLcEventHistory("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if lcState is LifecycleHistory {
        test:assertTrue(true, "Successfully retrive the LC events");
    } else {
        test:assertFail("Error occured while retrive LC events");
    }
}

@test:Config {}
function getLifeCyleStateTest() {
    LifecycleState|error lcState1 = getLifeCyleState("01ed75e2-b30b-18c8-wwf2-25da7edd2231", "carbon.super");
    if lcState1 is LifecycleState {
        test:assertTrue(true, "Successfully getting the LC state");
    } else {
        test:assertFail("Error occured while getting LC state");
    }
}
