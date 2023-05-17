import devportal_service.types;
import ballerina/cache;
import devportal_service.nonprodidp;
import wso2/apk_common_lib as commons;
import devportal_service.kmclient;

final cache:Cache kmClientCache = new (capacity = 50, evictionFactor = 0.2);

public isolated function getKmClient(types:KeyManager keyManagerConfig) returns kmclient:KeyManagerClient|commons:APKError {
    do {
        if (kmClientCache.hasKey(<string>keyManagerConfig.id)) {
            lock {
                return <kmclient:KeyManagerClient>check kmClientCache.get(<string>keyManagerConfig.id);
            }
        } else {
            if (kmClientCache.hasKey(<string>keyManagerConfig.id)) {
                lock {
                    return <kmclient:KeyManagerClient>check kmClientCache.get(<string>keyManagerConfig.id);
                }
            }
            lock {
                if keyManagerConfig.'type == "nonProdIdp" {
                    nonprodidp:NonProdIdpKeyManagerClient nonProdIdpClient = check new (keyManagerConfig);
                    _ = check kmClientCache.put(<string>keyManagerConfig.id, nonProdIdpClient);
                    return nonProdIdpClient;
                } else {
                    return error("Unsupported key manager type", code = 900959, description = "Unsupported key manager type", statusCode = 400, message = "Unsupported key manager type");
                }
            }
        }
    } on fail var e {
        return error("Internal Server Error", code = 900500, description = e.message(), statusCode = 500, message = e.message());
    }
}
