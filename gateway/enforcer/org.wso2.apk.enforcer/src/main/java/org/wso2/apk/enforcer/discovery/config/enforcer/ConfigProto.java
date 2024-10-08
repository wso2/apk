// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: wso2/discovery/config/enforcer/config.proto

package org.wso2.apk.enforcer.discovery.config.enforcer;

public final class ConfigProto {
  private ConfigProto() {}
  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistryLite registry) {
  }

  public static void registerAllExtensions(
      com.google.protobuf.ExtensionRegistry registry) {
    registerAllExtensions(
        (com.google.protobuf.ExtensionRegistryLite) registry);
  }
  static final com.google.protobuf.Descriptors.Descriptor
    internal_static_wso2_discovery_config_enforcer_Config_descriptor;
  static final 
    com.google.protobuf.GeneratedMessageV3.FieldAccessorTable
      internal_static_wso2_discovery_config_enforcer_Config_fieldAccessorTable;

  public static com.google.protobuf.Descriptors.FileDescriptor
      getDescriptor() {
    return descriptor;
  }
  private static  com.google.protobuf.Descriptors.FileDescriptor
      descriptor;
  static {
    java.lang.String[] descriptorData = {
      "\n+wso2/discovery/config/enforcer/config." +
      "proto\022\036wso2.discovery.config.enforcer\032)w" +
      "so2/discovery/config/enforcer/cert.proto" +
      "\032,wso2/discovery/config/enforcer/service" +
      ".proto\0322wso2/discovery/config/enforcer/j" +
      "wt_generator.proto\032*wso2/discovery/confi" +
      "g/enforcer/cache.proto\032.wso2/discovery/c" +
      "onfig/enforcer/analytics.proto\032-wso2/dis" +
      "covery/config/enforcer/security.proto\032/w" +
      "so2/discovery/config/enforcer/management" +
      ".proto\032+wso2/discovery/config/enforcer/f" +
      "ilter.proto\032,wso2/discovery/config/enfor" +
      "cer/tracing.proto\032,wso2/discovery/config" +
      "/enforcer/metrics.proto\032)wso2/discovery/" +
      "config/enforcer/soap.proto\032+wso2/discove" +
      "ry/config/enforcer/client.proto\"\212\007\n\006Conf" +
      "ig\022:\n\010security\030\001 \001(\0132(.wso2.discovery.co" +
      "nfig.enforcer.Security\022;\n\010keystore\030\002 \001(\013" +
      "2).wso2.discovery.config.enforcer.CertSt" +
      "ore\022=\n\ntruststore\030\003 \001(\0132).wso2.discovery" +
      ".config.enforcer.CertStore\022<\n\013authServic" +
      "e\030\004 \001(\0132\'.wso2.discovery.config.enforcer" +
      ".Service\022B\n\014jwtGenerator\030\005 \001(\0132,.wso2.di" +
      "scovery.config.enforcer.JWTGenerator\0224\n\005" +
      "cache\030\006 \001(\0132%.wso2.discovery.config.enfo" +
      "rcer.Cache\022<\n\tanalytics\030\007 \001(\0132).wso2.dis" +
      "covery.config.enforcer.Analytics\022>\n\nmana" +
      "gement\030\010 \001(\0132*.wso2.discovery.config.enf" +
      "orcer.Management\0228\n\007tracing\030\t \001(\0132\'.wso2" +
      ".discovery.config.enforcer.Tracing\0228\n\007me" +
      "trics\030\n \001(\0132\'.wso2.discovery.config.enfo" +
      "rcer.Metrics\0227\n\007filters\030\013 \003(\0132&.wso2.dis" +
      "covery.config.enforcer.Filter\0222\n\004soap\030\014 " +
      "\001(\0132$.wso2.discovery.config.enforcer.Soa" +
      "p\022%\n\035mandateSubscriptionValidation\030\r \001(\010" +
      "\022>\n\nhttpClient\030\016 \001(\0132*.wso2.discovery.co" +
      "nfig.enforcer.HttpClient\022$\n\034mandateInter" +
      "nalKeyValidation\030\017 \001(\010\022$\n\034enableGatewayC" +
      "lassController\030\020 \001(\010B\220\001\n/org.wso2.apk.en" +
      "forcer.discovery.config.enforcerB\013Config" +
      "ProtoP\001ZNgithub.com/envoyproxy/go-contro" +
      "l-plane/wso2/discovery/config/enforcer;e" +
      "nforcerb\006proto3"
    };
    descriptor = com.google.protobuf.Descriptors.FileDescriptor
      .internalBuildGeneratedFileFrom(descriptorData,
        new com.google.protobuf.Descriptors.FileDescriptor[] {
          org.wso2.apk.enforcer.discovery.config.enforcer.CertStoreProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.ServiceProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.JWTGeneratorProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.CacheProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.AnalyticsProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.SecurityProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.ManagementProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.FilterProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.TracingProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.MetricsProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.soapProto.getDescriptor(),
          org.wso2.apk.enforcer.discovery.config.enforcer.HttpClientProto.getDescriptor(),
        });
    internal_static_wso2_discovery_config_enforcer_Config_descriptor =
      getDescriptor().getMessageTypes().get(0);
    internal_static_wso2_discovery_config_enforcer_Config_fieldAccessorTable = new
      com.google.protobuf.GeneratedMessageV3.FieldAccessorTable(
        internal_static_wso2_discovery_config_enforcer_Config_descriptor,
        new java.lang.String[] { "Security", "Keystore", "Truststore", "AuthService", "JwtGenerator", "Cache", "Analytics", "Management", "Tracing", "Metrics", "Filters", "Soap", "MandateSubscriptionValidation", "HttpClient", "MandateInternalKeyValidation", "EnableGatewayClassController", });
    org.wso2.apk.enforcer.discovery.config.enforcer.CertStoreProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.ServiceProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.JWTGeneratorProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.CacheProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.AnalyticsProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.SecurityProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.ManagementProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.FilterProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.TracingProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.MetricsProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.soapProto.getDescriptor();
    org.wso2.apk.enforcer.discovery.config.enforcer.HttpClientProto.getDescriptor();
  }

  // @@protoc_insertion_point(outer_class_scope)
}
