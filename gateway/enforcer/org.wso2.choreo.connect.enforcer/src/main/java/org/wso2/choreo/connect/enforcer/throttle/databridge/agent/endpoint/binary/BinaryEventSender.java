/*
 * Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.wso2.choreo.connect.enforcer.throttle.databridge.agent.endpoint.binary;

import org.wso2.carbon.databridge.commons.Event;
import org.wso2.carbon.databridge.commons.binary.BinaryMessageConstants;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import static org.wso2.carbon.databridge.commons.binary.BinaryMessageConverterUtil.assignData;
import static org.wso2.carbon.databridge.commons.binary.BinaryMessageConverterUtil.getSize;
import static org.wso2.carbon.databridge.commons.binary.BinaryMessageConverterUtil.loadData;

/**
 * This is a Util class which does the Binary message transformation for publish, login, logout operations.
 */
public class BinaryEventSender {
    public static void sendBinaryLoginMessage(Socket socket, String userName, String password) throws IOException {
        ByteBuffer buf = ByteBuffer.allocate(13 + userName.length() + password.length());
        buf.put((byte) 0);
        buf.putInt(8 + userName.length() + password.length());
        buf.putInt(userName.length());
        buf.putInt(password.length());
        buf.put(userName.getBytes(BinaryMessageConstants.DEFAULT_CHARSET));
        buf.put(password.getBytes(BinaryMessageConstants.DEFAULT_CHARSET));

        OutputStream outputStream = new BufferedOutputStream(socket.getOutputStream());
        outputStream.write(buf.array());
        outputStream.flush();
    }

    public static void sendBinaryLogoutMessage(Socket socket, String sessionId) throws IOException {
        ByteBuffer buf = ByteBuffer.allocate(9 + sessionId.length());
        buf.put((byte) 1);
        buf.putInt(4 + sessionId.length());
        buf.putInt(sessionId.length());
        buf.put(sessionId.getBytes(BinaryMessageConstants.DEFAULT_CHARSET));

        OutputStream outputStream = new BufferedOutputStream(socket.getOutputStream());
        outputStream.write(buf.array());
        outputStream.flush();
    }

    public static void sendBinaryPublishMessage(Socket socket, List<Event> events, String sessionId)
            throws IOException {
        int messageSize = 8 + sessionId.length();
        List<byte[]> bytes = new ArrayList<>();

        for (Event event : events) {

            int eventSize = getEventSize(event);
            messageSize += eventSize + 4;
            ByteBuffer eventDataBuffer = ByteBuffer.allocate(4 + eventSize);
            eventDataBuffer.putInt(eventSize);
            eventDataBuffer.putLong(event.getTimeStamp());
            eventDataBuffer.putInt(event.getStreamId().length());
            eventDataBuffer.put(event.getStreamId().getBytes(BinaryMessageConstants.DEFAULT_CHARSET));

            if (event.getMetaData() != null && event.getMetaData().length != 0) {
                for (Object aMetaData : event.getMetaData()) {
                    assignData(aMetaData, eventDataBuffer);
                }
            }
            if (event.getCorrelationData() != null && event.getCorrelationData().length != 0) {
                for (Object aCorrelationData : event.getCorrelationData()) {
                    assignData(aCorrelationData, eventDataBuffer);
                }
            }
            if (event.getPayloadData() != null && event.getPayloadData().length != 0) {
                for (Object aPayloadData : event.getPayloadData()) {
                    assignData(aPayloadData, eventDataBuffer);
                }
            }
            if (event.getArbitraryDataMap() != null && event.getArbitraryDataMap().size() != 0) {
                for (Map.Entry<String, String> aArbitraryData : event.getArbitraryDataMap().entrySet()) {
                    assignData(aArbitraryData.getKey(), eventDataBuffer);
                    assignData(aArbitraryData.getValue(), eventDataBuffer);
                }
            }
            bytes.add(eventDataBuffer.array());
        }

        ByteBuffer buf = ByteBuffer.allocate(sessionId.length() + 13);
        buf.put((byte) 2);  //1
        buf.putInt(messageSize); //4
        buf.putInt(sessionId.length()); //4
        buf.put(sessionId.getBytes(BinaryMessageConstants.DEFAULT_CHARSET));
        buf.putInt(events.size()); //4

        OutputStream outputstream = new BufferedOutputStream(socket.getOutputStream());
        outputstream.write(buf.array());
        for (byte[] byteArray : bytes) {
            outputstream.write(byteArray);
        }
        outputstream.flush();
    }

    private static int getEventSize(Event event) {
        int eventSize = 4 + event.getStreamId().length() + 8;
        Object[] data = event.getMetaData();
        if (data != null) {
            for (Object aData : data) {
                eventSize += getSize(aData);
            }
        }
        data = event.getCorrelationData();
        if (data != null) {
            for (Object aData : data) {
                eventSize += getSize(aData);
            }
        }
        data = event.getPayloadData();
        if (data != null) {
            for (Object aData : data) {
                eventSize += getSize(aData);
            }
        }
        if (event.getArbitraryDataMap() != null && event.getArbitraryDataMap().size() != 0) {
            for (Map.Entry<String, String> aArbitraryData : event.getArbitraryDataMap().entrySet()) {
                eventSize += 8 + aArbitraryData.getKey().length() + aArbitraryData.getValue().length();
            }
        }
        return eventSize;
    }

    public static String processResponse(Socket socket) throws Exception {

        InputStream inputStream = socket.getInputStream();
        BufferedInputStream bufferedInputStream = new BufferedInputStream(inputStream);
        int messageType = bufferedInputStream.read();
        ByteBuffer bbuf;
        switch (messageType) {
            case 0:
                //OK message
                break;
            case 1:
                //Error Message
                bbuf = ByteBuffer.wrap(loadData(bufferedInputStream, new byte[8]));
                int errorClassNameLength = bbuf.getInt();
                int errorMsgLength = bbuf.getInt();

                String className = new String(ByteBuffer.wrap(loadData(bufferedInputStream,
                        new byte[errorClassNameLength])).array());
                String errorMsg = new String(ByteBuffer.wrap(loadData(bufferedInputStream,
                        new byte[errorMsgLength])).array());

                throw (Exception) (BinaryDataEndpoint.class.getClassLoader().
                        loadClass(className).getConstructor(String.class).newInstance(errorMsg));
            case 2:
                //Logging OK response
                bbuf = ByteBuffer.wrap(loadData(bufferedInputStream, new byte[4]));
                int sessionIdLength = bbuf.getInt();
                return new String(ByteBuffer.wrap(loadData(bufferedInputStream, new byte[sessionIdLength])).array());
        }
        return null;
    }
}
