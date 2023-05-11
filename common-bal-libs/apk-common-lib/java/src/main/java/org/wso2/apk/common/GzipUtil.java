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

package org.wso2.apk.common;

import java.io.*;
import java.util.zip.GZIPInputStream;
import java.util.zip.GZIPOutputStream;

/**
 * This class used to gzip related utilities.
 */
public class GzipUtil {

    /**
     * This method used to compress a file using gzip.
     *
     * @param filePath       file path
     * @param outputFilePath output file path
     * @throws IOException when error occurred while compressing the file
     */
    public static void compressGzipFile(String filePath, String outputFilePath) throws IOException {
        try (FileInputStream fileInputStream = new FileInputStream(filePath)) {
            try (GZIPOutputStream gzipOutputStream = new GZIPOutputStream(new FileOutputStream(outputFilePath))) {
                byte[] buffer = new byte[1024];
                int len;
                while ((len = fileInputStream.read(buffer)) != -1) {
                    gzipOutputStream.write(buffer, 0, len);
                }
            }
        }
    }

    /**
     * This method used to decompress a gzip file.
     *
     * @param gzipFilePath   gzip file path
     * @param outputFilePath output file path
     * @throws IOException when error occurred while decompressing the file
     */
    public static void decompressGzipFile(String gzipFilePath, String outputFilePath) throws IOException {
        try (GZIPInputStream gzipInputStream = new GZIPInputStream(new FileInputStream(gzipFilePath))) {
            try (FileOutputStream fileOutputStream = new FileOutputStream(outputFilePath)) {
                byte[] buffer = new byte[1024];
                int len;
                while ((len = gzipInputStream.read(buffer)) != -1) {
                    fileOutputStream.write(buffer, 0, len);
                }
            }
        }
    }

    /**
     * This method used to compress a file using gzip.
     *
     * @param data file as a byte array
     * @return compressed file as a byte array
     * @throws IOException when error occurred while compressing the file
     */
    public static byte[] compressGzipFile(byte[] data) throws IOException {
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            try (GZIPOutputStream gzipOutputStream = new GZIPOutputStream(byteArrayOutputStream)) {
                gzipOutputStream.write(data);
            }
            return byteArrayOutputStream.toByteArray();
        }
    }

    /**
     * This method used to decompress a gzip file.
     * @param data compressed file as a byte array
     * @return decompressed file as a byte array
     * @throws IOException when error occurred while decompressing the file
     */
    public static byte[] decompressGzipFile(byte[] data) throws IOException{
        try (ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream()) {
            try (ByteArrayInputStream byteArrayInputStream = new ByteArrayInputStream(data)) {
                GZIPInputStream gzipInputStream = new GZIPInputStream(byteArrayInputStream);
                byte[] buffer = new byte[1024];
                int len;
                while ((len = gzipInputStream.read(buffer)) != -1) {
                    byteArrayOutputStream.write(buffer, 0, len);
                }
            }
            return byteArrayOutputStream.toByteArray();
        }
    }
}
