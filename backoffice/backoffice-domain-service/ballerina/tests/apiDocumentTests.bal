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

string documentId = "";
@test:Config {dependsOn: [createAPITest]}
function addDocumentMetaDataTest() {
    Document documentBody = {
        name:"NewDoc",
        documentType: "HOWTO",
        summary: "Doc summary",
        sourceType: "INLINE",
        visibility: "API_LEVEL",
        sourceUrl: "",
        otherTypeName: "sdds",
        inlineContent: ""
    };

    Document|commons:APKError createdDocument = createDocument("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentBody);
    if createdDocument is Document {
        documentId = <string>createdDocument.documentId;
        test:assertTrue(true, "Successfully added the thumbnail");
    } else {
        test:assertFail("Error occured while adding the Document");
    }
}

@test:Config {dependsOn: [addDocumentMetaDataTest]}
function updateDocumentMetaDataTest() {
    Document updatedDocumentBody = {
        name:"NewDoc",
        documentType: "HOWTO",
        summary: "Doc summary updated",
        sourceType: "INLINE",
        visibility: "API_LEVEL",
        sourceUrl: "",
        otherTypeName: "sdds",
        inlineContent: ""
    };

    Document|NotFoundError|commons:APKError|error updatedDocument = UpdateDocumentMetaData("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentId, updatedDocumentBody);
    if  updatedDocument is Document {
        test:assertTrue(true, "Successfully added the thumbnail");
        test:assertEquals(updatedDocumentBody.summary,  updatedDocument.summary, "Updated value should be equal");
    } else {
        test:assertFail("Error occured while adding the Document");
    }
}

@test:Config {dependsOn: [createAPITest, addDocumentMetaDataTest]}
function addDocumentContentTest() {
    http:Request|error request = createRequestWithImageFormData("invalidThumbnail.pdf", "application/pdf");
    if request is http:Request {
        Document|NotFoundError|commons:APKError|error addedContent = addDocumentContent("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentId, request);
        if addedContent is Document {
            test:assertTrue(true, "Successfully added the Document content");
        } else {
            test:assertFail("Error occured while adding Document content");
        }
    }
}

@test:Config {dependsOn: [addDocumentMetaDataTest]}
function getDocumentMetaDataTest() {
    Document|NotFoundError|commons:APKError docMetaData = getDocumentMetaData("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentId);
    if docMetaData is Document {
        test:assertTrue(true, "Successfully getting the Document meta data");
    } else {
        test:assertFail("Error occured while getting the Document meta data");
    }
}

@test:Config {dependsOn: [addDocumentContentTest]}
function getDocumentContentTest() {
    http:Response|NotFoundError|commons:APKError docContent = getDocumentContent("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentId);
    if docContent is http:Response {
        test:assertTrue(true, "Successfully getting the Document content");
    } else {
        test:assertFail("Error occured while getting the Document content");
    }
}

@test:Config {dependsOn: [addDocumentMetaDataTest]}
function getDocumentListTest() {
    DocumentList|commons:APKError documentList = getDocumentList("01ed75e2-b30b-18c8-wwf2-25da7edd2231", 25, 0);
    if documentList is DocumentList {
        test:assertTrue(true, "Successfully getting the Document List");
        test:assertEquals(documentList.count, 1, "Document count should be equal to 1");
    } else {
        test:assertFail("Error occured while getting the Document List");
    }
}

@test:Config {}
function deleteDocumentTest() {
    http:Ok|NotFoundError|commons:APKError deletdeDoc = deleteDocument("01ed75e2-b30b-18c8-wwf2-25da7edd2231", documentId);
    if deletdeDoc is http:Ok {
        test:assertTrue(true, "Successfully deleted the Document");
    } else {
        test:assertFail("Error occured while deleting the Document");
    }
}
