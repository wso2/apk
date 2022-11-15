# APK Admin Domain Service

This contains the Ballerina implementation of the Admin Domain Service.

## OpenAPI generation

The REST API skeleton for Admin Service is generated using the [Ballerina OpenAPI tool](https://lib.ballerina.io/ballerina/openapi/1.2.1). The OpenAPI definition for the Admin Service is available in [admin-api.yaml](/ballerina/modules/admin/resources/admin-api.yaml)

```
bal openapi -i admin-api.yaml --mode service
```

The above command will generate the service bal file and types.bal file.

## Bridge code generation

[Ballerina bindgen tool](https://ballerina.io/learn/java-interoperability-guide/the-bindgen-tool/) is used to generate the bridge code for the business logic available in Java code. The generated bridge code will be available under modules.

```
bal bindgen -mvn org.wso2.apk:org.wso2.apk.apimgt.rest.api.admin.v1.common:0.1.0-SNAPSHOT org.wso2.apk.apimgt.rest.api.admin.v1.common.impl.ThrottlingCommonImp
```

```
bal bindgen -mvn org.wso2.apk:org.wso2.apk.apimgt.init:0.1.0-SNAPSHOT org.wso2.apk.apimgt.init.APKComponent
```

The bridge code generated using the above sample commands can be directly used in the bal files. This allows us to directly call the business logic in Java.

## Configuration model

- The configurations should be specified in the Config.toml
- A seperate record type for each configuration should be created in [config.bal](/ballerina/modules/admin/config.bal). This will create the mappings between the values in the toml files through configurable variables.
- Any new configuration should be added to `APKConfig` record. This config will be passed as a json string to the `APKComponent` through the generated bridge code.
- The necessary property classes should be created and linked to the [Configuration Holder](../../../common-java-libs/org.wso2.apk.apimgt.impl/src/main/java/org/wso2/apk/apimgt/impl/ConfigurationHolder.java).
- The configurations will be mapped during the runtime and will be available through e Reference Holder in Java code.



![config-model](/ballerina/modules/admin/resources/apkconf.png)
