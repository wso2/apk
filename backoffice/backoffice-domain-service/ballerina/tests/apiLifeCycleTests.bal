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

@test:Mock { functionName: "db_changeLCState" }
test:MockFunction db_changeLCStateMock = new();

@test:Mock { functionName: "db_getCurrentLCStatus" }
test:MockFunction db_getCurrentLCStatusMock = new();

@test:Mock { functionName: "db_AddLCEvent" }
test:MockFunction db_AddLCEventMock = new();

@test:Mock { functionName: "db_getLCEventHistory" }
test:MockFunction db_getLCEventHistoryMock = new();



@test:Config {}
function changeLCStateTest() {
    test:when(db_getCurrentLCStatusMock).thenReturn("CREATED");
    test:when(db_changeLCStateMock).thenReturn("PUBLISHED");
    test:when(db_AddLCEventMock).thenReturn("PUBLISHED");
    LifecycleState|error lcState1 = changeLifeCyleState("PUBLISHED", "ap01ed7552-b30b-18c8-wwf2-25da7a7c46ceiId", "carbon.super");
    if lcState1 is LifecycleState {
        test:assertTrue(true, "Successfully change the LC state");
    } else {
        test:assertFail("Error occured while changing LC state");
    }
}


@test:Config {}
function getLcEventHistoryTest() {
    LifecycleHistoryItem[] | error lc = 
        [{

            previousState: "Created",
            postState: "Published",
            user: "admin",
            updatedTime: "2019-02-31T23:59:60Z"
        }];
    test:when(db_getLCEventHistoryMock).thenReturn(lc);
    LifecycleHistory|error? lcState = getLcEventHistory("ap01ed7552-b30b-18c8-wwf2-25da7a7c46ceiId");
    if lcState is LifecycleHistory {
        test:assertTrue(true, "Successfully retrive the LC events");
    } else {
        test:assertFail("Error occured while retrive LC events");
    }
}

@test:Config {}
function getLifeCyleStateTest() {
    test:when(db_getCurrentLCStatusMock).thenReturn("PUBLISHED");
    LifecycleState|error lcState1 = getLifeCyleState("ap01ed7552-b30b-18c8-wwf2-25da7a7c46ceiId", "carbon.super");
    if lcState1 is LifecycleState {
        test:assertTrue(true, "Successfully getting the LC state");
    } else {
        test:assertFail("Error occured while getting LC state");
    }
}
