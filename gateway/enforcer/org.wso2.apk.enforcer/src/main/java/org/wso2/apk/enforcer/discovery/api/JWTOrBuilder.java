// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: wso2/discovery/api/api_authentication.proto

package org.wso2.apk.enforcer.discovery.api;

public interface JWTOrBuilder extends
    // @@protoc_insertion_point(interface_extends:wso2.discovery.api.JWT)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <pre>
   * name of the header containing the JWT
   * </pre>
   *
   * <code>string header = 1;</code>
   * @return The header.
   */
  java.lang.String getHeader();
  /**
   * <pre>
   * name of the header containing the JWT
   * </pre>
   *
   * <code>string header = 1;</code>
   * @return The bytes for header.
   */
  com.google.protobuf.ByteString
      getHeaderBytes();

  /**
   * <pre>
   * send the token to upstream
   * </pre>
   *
   * <code>bool sendTokenToUpstream = 2;</code>
   * @return The sendTokenToUpstream.
   */
  boolean getSendTokenToUpstream();

  /**
   * <code>repeated string audience = 3;</code>
   * @return A list containing the audience.
   */
  java.util.List<java.lang.String>
      getAudienceList();
  /**
   * <code>repeated string audience = 3;</code>
   * @return The count of audience.
   */
  int getAudienceCount();
  /**
   * <code>repeated string audience = 3;</code>
   * @param index The index of the element to return.
   * @return The audience at the given index.
   */
  java.lang.String getAudience(int index);
  /**
   * <code>repeated string audience = 3;</code>
   * @param index The index of the value to return.
   * @return The bytes of the audience at the given index.
   */
  com.google.protobuf.ByteString
      getAudienceBytes(int index);
}
