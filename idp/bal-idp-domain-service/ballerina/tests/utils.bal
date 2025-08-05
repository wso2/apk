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
import ballerina/url;

# Get Encoded URI for a given value.
#
# + value - Value to be encoded
# + return - Encoded string
isolated function getEncodedUri(anydata value) returns string {
    string|error encoded = url:encode(value.toString(), "UTF8");
    if (encoded is string) {
        return encoded;
    } else {
        return value.toString();
    }
}

type SimpleBasicType string|boolean|int|float|decimal;

# Represents encoding mechanism details.
type Encoding record {
    # Defines how multiple values are delimited
    string style = FORM;
    # Specifies whether arrays and objects should generate as separate fields
    boolean explode = true;
    # Specifies the custom content type
    string contentType?;
    # Specifies the custom headers
    map<any> headers?;
};

enum EncodingStyle {
    DEEPOBJECT, FORM, SPACEDELIMITED, PIPEDELIMITED
}

final Encoding & readonly defaultEncoding = {};

# Generate client request when the media type is given as application/x-www-form-urlencoded.
#
# + encodingMap - Includes the information about the encoding mechanism
# + anyRecord - Record to be serialized
# + return - Serialized request body or query parameter as a string
isolated function createFormURLEncodedRequestBody(record {|anydata...;|} anyRecord, map<Encoding> encodingMap = {}) returns string {
    string[] payload = [];
    foreach [string, anydata] [key, value] in anyRecord.entries() {
        Encoding encodingData = encodingMap.hasKey(key) ? encodingMap.get(key) : defaultEncoding;
        if value is SimpleBasicType {
            payload.push(key, "=", getEncodedUri(value.toString()));
        } else if value is SimpleBasicType[] {
            payload.push(getSerializedArray(key, value, encodingData.style, encodingData.explode));
        } else if (value is record {}) {
            if encodingData.style == DEEPOBJECT {
                payload.push(getDeepObjectStyleRequest(key, value));
            } else {
                payload.push(getFormStyleRequest(key, value));
            }
        } else if (value is record {}[]) {
            payload.push(getSerializedRecordArray(key, value, encodingData.style, encodingData.explode));
        }
        payload.push("&");
    }
    _ = payload.pop();
    return string:'join("", ...payload);
}

# Serialize the record according to the deepObject style.
#
# + parent - Parent record name
# + anyRecord - Record to be serialized
# + return - Serialized record as a string
isolated function getDeepObjectStyleRequest(string parent, record {} anyRecord) returns string {
    string[] recordArray = [];
    foreach [string, anydata] [key, value] in anyRecord.entries() {
        if value is SimpleBasicType {
            recordArray.push(parent + "[" + key + "]" + "=" + getEncodedUri(value.toString()));
        } else if value is SimpleBasicType[] {
            recordArray.push(getSerializedArray(parent + "[" + key + "]" + "[]", value, DEEPOBJECT, true));
        } else if value is record {} {
            string nextParent = parent + "[" + key + "]";
            recordArray.push(getDeepObjectStyleRequest(nextParent, value));
        } else if value is record {}[] {
            string nextParent = parent + "[" + key + "]";
            recordArray.push(getSerializedRecordArray(nextParent, value, DEEPOBJECT));
        }
        recordArray.push("&");
    }
    _ = recordArray.pop();
    return string:'join("", ...recordArray);
}

# Serialize the record according to the form style.
#
# + parent - Parent record name
# + anyRecord - Record to be serialized
# + explode - Specifies whether arrays and objects should generate separate parameters
# + return - Serialized record as a string
isolated function getFormStyleRequest(string parent, record {} anyRecord, boolean explode = true) returns string {
    string[] recordArray = [];
    if explode {
        foreach [string, anydata] [key, value] in anyRecord.entries() {
            if (value is SimpleBasicType) {
                recordArray.push(key, "=", getEncodedUri(value.toString()));
            } else if (value is SimpleBasicType[]) {
                recordArray.push(getSerializedArray(key, value, explode = explode));
            } else if (value is record {}) {
                recordArray.push(getFormStyleRequest(parent, value, explode));
            }
            recordArray.push("&");
        }
        _ = recordArray.pop();
    } else {
        foreach [string, anydata] [key, value] in anyRecord.entries() {
            if (value is SimpleBasicType) {
                recordArray.push(key, ",", getEncodedUri(value.toString()));
            } else if (value is SimpleBasicType[]) {
                recordArray.push(getSerializedArray(key, value, explode = false));
            } else if (value is record {}) {
                recordArray.push(getFormStyleRequest(parent, value, explode));
            }
            recordArray.push(",");
        }
        _ = recordArray.pop();
    }
    return string:'join("", ...recordArray);
}

# Serialize arrays.
#
# + arrayName - Name of the field with arrays
# + anyArray - Array to be serialized
# + style - Defines how multiple values are delimited
# + explode - Specifies whether arrays and objects should generate separate parameters
# + return - Serialized array as a string
isolated function getSerializedArray(string arrayName, anydata[] anyArray, string style = "form", boolean explode = true) returns string {
    string key = arrayName;
    string[] arrayValues = [];
    if (anyArray.length() > 0) {
        if (style == FORM && !explode) {
            arrayValues.push(key, "=");
            foreach anydata i in anyArray {
                arrayValues.push(getEncodedUri(i.toString()), ",");
            }
        } else if (style == SPACEDELIMITED && !explode) {
            arrayValues.push(key, "=");
            foreach anydata i in anyArray {
                arrayValues.push(getEncodedUri(i.toString()), "%20");
            }
        } else if (style == PIPEDELIMITED && !explode) {
            arrayValues.push(key, "=");
            foreach anydata i in anyArray {
                arrayValues.push(getEncodedUri(i.toString()), "|");
            }
        } else if (style == DEEPOBJECT) {
            foreach anydata i in anyArray {
                arrayValues.push(key, "[]", "=", getEncodedUri(i.toString()), "&");
            }
        } else {
            foreach anydata i in anyArray {
                arrayValues.push(key, "=", getEncodedUri(i.toString()), "&");
            }
        }
        _ = arrayValues.pop();
    }
    return string:'join("", ...arrayValues);
}

# Serialize the array of records according to the form style.
#
# + parent - Parent record name
# + value - Array of records to be serialized
# + style - Defines how multiple values are delimited
# + explode - Specifies whether arrays and objects should generate separate parameters
# + return - Serialized record as a string
isolated function getSerializedRecordArray(string parent, record {}[] value, string style = FORM, boolean explode = true) returns string {
    string[] serializedArray = [];
    if style == DEEPOBJECT {
        int arayIndex = 0;
        foreach var recordItem in value {
            serializedArray.push(getDeepObjectStyleRequest(parent + "[" + arayIndex.toString() + "]", recordItem), "&");
            arayIndex = arayIndex + 1;
        }
    } else {
        if (!explode) {
            serializedArray.push(parent, "=");
        }
        foreach var recordItem in value {
            serializedArray.push(getFormStyleRequest(parent, recordItem, explode), ",");
        }
    }
    _ = serializedArray.pop();
    return string:'join("", ...serializedArray);
}

# Generate query path with query parameter.
#
# + queryParam - Query parameter map
# + encodingMap - Details on serialization mechanism
# + return - Returns generated Path or error at failure of client initialization
isolated function getPathForQueryParam(map<anydata> queryParam, map<Encoding> encodingMap = {}) returns string|error {
    string[] param = [];
    if (queryParam.length() > 0) {
        param.push("?");
        foreach var [key, value] in queryParam.entries() {
            if value is () {
                _ = queryParam.remove(key);
                continue;
            }
            Encoding encodingData = encodingMap.hasKey(key) ? encodingMap.get(key) : defaultEncoding;
            if (value is SimpleBasicType) {
                param.push(key, "=", getEncodedUri(value.toString()));
            } else if (value is SimpleBasicType[]) {
                param.push(getSerializedArray(key, value, encodingData.style, encodingData.explode));
            } else if (value is record {}) {
                if (encodingData.style == DEEPOBJECT) {
                    param.push(getDeepObjectStyleRequest(key, value));
                } else {
                    param.push(getFormStyleRequest(key, value, encodingData.explode));
                }
            } else {
                param.push(key, "=", value.toString());
            }
            param.push("&");
        }
        _ = param.pop();
    }
    string restOfPath = string:'join("", ...param);
    return restOfPath;
}

# Generate header map for given header values.
#
# + headerParam - Headers  map
# + return - Returns generated map or error at failure of client initialization
isolated function getMapForHeaders(map<any> headerParam) returns map<string|string[]> {
    map<string|string[]> headerMap = {};
    foreach var [key, value] in headerParam.entries() {
        if value is string || value is string[] {
            headerMap[key] = value;
        } else if value is int[] {
            string[] stringArray = [];
            foreach int intValue in value {
                stringArray.push(intValue.toString());
            }
            headerMap[key] = stringArray;
        } else if value is SimpleBasicType {
            headerMap[key] = value.toString();
        }
    }
    return headerMap;
}

