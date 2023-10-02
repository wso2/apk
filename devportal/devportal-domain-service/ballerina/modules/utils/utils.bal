import ballerina/io;
import ballerina/file;
import wso2/apk_common_lib as commons;

public isolated function storeAndRetrieveCertificate(string fileName, string base64EncodedContent) returns string|error {
    string tempDir = check file:createTempDir();
    string filePath = tempDir + file:pathSeparator + fileName;
    byte[] base64DecodedContent = check commons:EncoderUtil_decodeBase64(base64EncodedContent.toBytes());
    _ = check io:fileWriteBytes(filePath, base64DecodedContent);
    return filePath;
}
