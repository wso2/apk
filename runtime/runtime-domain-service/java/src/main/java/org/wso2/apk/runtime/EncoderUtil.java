/*
 * Copyright (c) 2023, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.apk.runtime;

/**
 * This class is used to encode and decode the base64 strings.
 */
public class EncoderUtil {
    /**
     * This method is used to encode the base64 string.
     *
     * @param data data to be encoded
     * @return encoded data
     */
    public static byte[] encodeBase64(byte[] data) {
        return java.util.Base64.getEncoder().encode(data);
    }

    /**
     * This method is used to encode the base64 string.
     *
     * @param data data to be encoded
     * @return encoded data
     */
    public static byte[] encodeBase64(String data) {
        return encodeBase64(data.getBytes());
    }

    /**
     * This method is used to decode the base64 string.
     *
     * @param data data to be decoded
     * @return decoded data
     */
    public static byte[] decodeBase64(byte[] data) {
        return java.util.Base64.getDecoder().decode(data);
    }
}
