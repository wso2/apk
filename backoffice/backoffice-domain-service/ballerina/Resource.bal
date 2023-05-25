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

type Resource record {|
    string resourceUUID;
    string apiUuid;
    int resourceCategoryId;
    string dataType;
    string resourceContent;
    byte[] resourceBinaryValue?;
    json respurceJsonValue?;
|};

type Thumbnail record {|
    string imageType;
    byte[] imageContent;
|};

type DocumentMetaData record {|
    string documentId?;
    string resourceId?;
    string name;
    string documentType;
    string summary?;
    string sourceType;
    string sourceUrl?;
    string fileName?;
    string inlineContent?;
    string otherTypeName?;
    string visibility;
|};
