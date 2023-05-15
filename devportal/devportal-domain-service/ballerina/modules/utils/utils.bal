import ballerina/io;
import wso2/apk_common_lib as commons;

public isolated function storeCertificate(string fliePath, string base64EncodedContent) returns error? {
    byte[] base64DecodedContent = check commons:EncoderUtil_decodeBase64(base64EncodedContent.toBytes());
    _ = check io:fileWriteBytes(fliePath, base64DecodedContent);
}
