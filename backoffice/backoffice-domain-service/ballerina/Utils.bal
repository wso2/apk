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
import ballerina/log;
import ballerina/mime;

isolated function isFileSizeGreaterThan1MB(byte[] data) returns boolean {
    int fileSizeInBytes = data.length();
    int fileSizeInMB = fileSizeInBytes / (1024 * 1024);
    if fileSizeInMB > 1 {
        log:printDebug("File is greater than 1MB");
        return true;
    }
    return false;
}

isolated function isThumbnailHasValidFileExtention(string contentType) returns boolean {
    if (contentType == RESOURCE_DATA_TYPE_JPG_IMAGE || contentType == RESOURCE_DATA_TYPE_PNG_IMAGE ||
    contentType == RESOURCE_DATA_TYPE_GIF_IMAGE || contentType == RESOURCE_DATA_TYPE_SVG_IMAGE) {
        return true;
    } else {
        return false;
    }
}

isolated function getContentBaseType(string contentType) returns string {
    var result = mime:getMediaType(contentType);
    if result is mime:MediaType {
        return result.getBaseType();
    }
    panic result;
}
