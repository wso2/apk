//
// Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
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
import ballerina/http;
import wso2/apk_common_lib as commons;

@test:Config {dependsOn: [createAPITest]}
function addThumbnailTest() {
    http:Request|error request = createRequestWithImageFormData("thumbnail.png", RESOURCE_DATA_TYPE_PNG_IMAGE);
    if request is http:Request {
        FileInfo|NotFoundError|PreconditionFailedError|commons:APKError|error thumbnail = updateThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231", request);
        if thumbnail is FileInfo {
            test:assertTrue(true, "Successfully added the thumbnail");
        } else {
            test:assertFail("Error occured while adding the thumbnail");
        }
    }
}

@test:Config {dependsOn: [createAPITest, addThumbnailTest]}
function updateThumbnailTest() {
     http:Request|error request = createRequestWithImageFormData("thumbnail.png", RESOURCE_DATA_TYPE_PNG_IMAGE);
    if request is http:Request {
        FileInfo|NotFoundError|PreconditionFailedError|commons:APKError|error thumbnail = updateThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231", request);
        if thumbnail is FileInfo {
            test:assertTrue(true, "Successfully updated the thumbnail");
        } else {
            test:assertFail("Error occured while updating the thumbnail");
        }
    }
}

@test:Config {dependsOn: [createAPITest, addThumbnailTest]}
function addThumbnaiGreaterThan1MB() {
     http:Request|error request = createRequestWithImageFormData("largeThumbnail.jpg", RESOURCE_DATA_TYPE_JPG_IMAGE);
    if request is http:Request {
        FileInfo|NotFoundError|PreconditionFailedError|commons:APKError|error thumbnail = updateThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231", request);
        if thumbnail is PreconditionFailedError {
            test:assertEquals(thumbnail.body.message, "Thumbnail size should be less than 1MB");
        } else {
            test:assertFail("Thumbnail size which is greater than 1MB is added");
        }
    }
}

@test:Config {dependsOn: [createAPITest, addThumbnailTest]}
function addInvalidThumbnaiFileFormat() {
     http:Request|error request = createRequestWithImageFormData("invalidThumbnail.pdf", "application/pdf");
    if request is http:Request {
        FileInfo|NotFoundError|PreconditionFailedError|commons:APKError|error thumbnail = updateThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231", request);
        if thumbnail is PreconditionFailedError {
            test:assertEquals(thumbnail.body.message, "Thumbnail file extension is not allowed. Supported extensions are .jpg, .png, .jpeg .svg and .gif");
        } else {
            test:assertFail("Thumbnail which has invalid file format is added");
        }
    }
}

@test:Config {}
function gethumbnailTest() {
    http:Response|NotFoundError|commons:APKError thumbnail = getThumbnail("01ed75e2-b30b-18c8-wwf2-25da7edd2231");
    if thumbnail is http:Response {
        test:assertTrue(true, "Successfully getting the thumbnail");
    } else {
        test:assertFail("Error occured while getting the thumbnail");
    }
}
