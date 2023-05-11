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

import ballerina/http;
import ballerina/mime;
import ballerina/io;

# Create http request including image form data
#
# + thumnailImageName - image name
# + return - http:Request | error
isolated function createRequestWithImageFormData(string thumnailImageName) returns http:Request|error {
    mime:Entity imageBodyPart = new;
    byte[] imageBytes = check io:fileReadBytes("./tests/resources/" + thumnailImageName);
    imageBodyPart.setByteArray(imageBytes);
    mime:InvalidContentTypeError? contentType = imageBodyPart.setContentType(mime:IMAGE_PNG);
    if contentType is mime:InvalidContentTypeError {
        return contentType;
    }
    mime:ContentDisposition contentDisposition = new;
    contentDisposition.disposition = "form-data";
    contentDisposition.name = "file";
    contentDisposition.fileName = thumnailImageName;
    imageBodyPart.setContentDisposition(contentDisposition);
    mime:Entity[] bodyParts = [imageBodyPart];
    http:Request request = new;
    request.setBodyParts(bodyParts);
    return request;
}
