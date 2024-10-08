// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: wso2/discovery/api/Resource.proto

package org.wso2.apk.enforcer.discovery.api;

public final class ResourceProto {
  private ResourceProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_api_Resource_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_api_Resource_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_api_Operation_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_api_Operation_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_api_OperationPolicies_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_api_OperationPolicies_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_api_Policy_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_api_Policy_fieldAccessorTable;
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_api_Policy_ParametersEntry_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_api_Policy_ParametersEntry_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n!wso2/discovery/api/Resource.proto\022\022wso" +
      "2.discovery.api\032)wso2/discovery/api/endp" +
      "oint_cluster.proto\032+wso2/discovery/api/a" +
      "pi_authentication.proto\032&wso2/discovery/" +
      "api/security_info.proto\"\310\001\n\010Resource\022\n\n\002" +
      "id\030\001 \001(\t\022\014\n\004path\030\002 \001(\t\022.\n\007methods\030\003 \003(\0132" +
      "\035.wso2.discovery.api.Operation\0226\n\tendpoi" +
      "nts\030\004 \001(\0132#.wso2.discovery.api.EndpointC" +
      "luster\022:\n\020endpointSecurity\030\005 \003(\0132 .wso2." +
      "discovery.api.SecurityInfo\"\305\001\n\tOperation" +
      "\022\016\n\006method\030\001 \001(\t\022@\n\021apiAuthentication\030\002 " +
      "\001(\0132%.wso2.discovery.api.APIAuthenticati" +
      "on\022\014\n\004tier\030\003 \001(\t\0227\n\010policies\030\004 \001(\0132%.wso" +
      "2.discovery.api.OperationPolicies\022\016\n\006sco" +
      "pes\030\005 \003(\t\022\017\n\007matchID\030\006 \001(\t\"\231\001\n\021Operation" +
      "Policies\022+\n\007request\030\001 \003(\0132\032.wso2.discove" +
      "ry.api.Policy\022,\n\010response\030\002 \003(\0132\032.wso2.d" +
      "iscovery.api.Policy\022)\n\005fault\030\003 \003(\0132\032.wso" +
      "2.discovery.api.Policy\"\213\001\n\006Policy\022\016\n\006act" +
      "ion\030\001 \001(\t\022>\n\nparameters\030\002 \003(\0132*.wso2.dis" +
      "covery.api.Policy.ParametersEntry\0321\n\017Par" +
      "ametersEntry\022\013\n\003key\030\001 \001(\t\022\r\n\005value\030\002 \001(\t" +
      ":\0028\001Bu\n#org.wso2.apk.enforcer.discovery." +
      "apiB\rResourceProtoP\001Z=github.com/envoypr" +
      "oxy/go-control-plane/wso2/discovery/api;" +
      "apib\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          org.wso2.apk.enforcer.discovery.api.EndpointClusterProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.api.APIAuthenticationProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.api.SecurityInfoProto.getDescriptor(),
        });
    internal_static_wso2_discovery_api_Resource_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_wso2_discovery_api_Resource_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_api_Resource_descriptor,
        new java.lang.String[] { "Id", "Path", "Methods", "Endpoints", "EndpointSecurity", });
    internal_static_wso2_discovery_api_Operation_descriptor =
      getDescriptor().getMessageTypes().get(1);
    internal_static_wso2_discovery_api_Operation_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_api_Operation_descriptor,
        new java.lang.String[] { "Method", "ApiAuthentication", "Tier", "Policies", "Scopes", "MatchID", });
    internal_static_wso2_discovery_api_OperationPolicies_descriptor =
      getDescriptor().getMessageTypes().get(2);
    internal_static_wso2_discovery_api_OperationPolicies_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_api_OperationPolicies_descriptor,
        new java.lang.String[] { "Request", "Response", "Fault", });
    internal_static_wso2_discovery_api_Policy_descriptor =
      getDescriptor().getMessageTypes().get(3);
    internal_static_wso2_discovery_api_Policy_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_api_Policy_descriptor,
        new java.lang.String[] { "Action", "Parameters", });
    internal_static_wso2_discovery_api_Policy_ParametersEntry_descriptor =
      internal_static_wso2_discovery_api_Policy_descriptor.getNestedTypes().get(0);
    internal_static_wso2_discovery_api_Policy_ParametersEntry_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_api_Policy_ParametersEntry_descriptor,
        new java.lang.String[] { "Key", "Value", });
    org.wso2.apk.enforcer.discovery.api.EndpointClusterProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.api.APIAuthenticationProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.api.SecurityInfoProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
