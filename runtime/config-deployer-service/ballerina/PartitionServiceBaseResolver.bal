import wso2/apk_common_lib as commons;
import config_deployer_service.partitionClient;
import ballerina/http;
import config_deployer_service.model;

public isolated class PartitionServiceBaseResolver {
    *PartitionResolver;
    final partitionClient:Client partitionClient;
    isolated function init(PartitionServiceConfiguration partitionServiceConfiguration) returns error? {
        partitionClient:ConnectionConfig connectionConfig = {};
        string url = <string>partitionServiceConfiguration.url;
        boolean httpsProtocol = url.startsWith("https://") ? true : false;
        if httpsProtocol {
            if partitionServiceConfiguration.tlsCertificatePath is string {
                connectionConfig.secureSocket = {
                    verifyHostName: partitionServiceConfiguration.hostnameVerificationEnable,
                    cert: partitionServiceConfiguration.tlsCertificatePath
                };
            } else {
                connectionConfig.secureSocket = {
                    verifyHostName: partitionServiceConfiguration.hostnameVerificationEnable
                };

            }
        }
        self.partitionClient = check new (<string>partitionServiceConfiguration.url, connectionConfig);
    }
    isolated function getAvailablePartitionForAPI(string id, string organization) returns model:Partition|commons:APKError? {
        partitionClient:Partition|error apiPartition = self.partitionClient->/api\-deployment/[id];
        if apiPartition is partitionClient:Partition {
            return {name: <string>apiPartition.name, namespace: <string>apiPartition.namespace, apiCount: apiPartition.apiCount};
        } else if apiPartition is http:ApplicationResponseError {
            if apiPartition.detail().statusCode == 404 {
                return ();
            } else {
                return e909022("Internal error occured ", apiPartition);
            }
        } else {
            return e909022("Internal error occured ", apiPartition);
        }
    }
    isolated function getDeployablePartition() returns model:Partition|commons:APKError {
        partitionClient:Partition|error deployablePartition = self.partitionClient->/deployable\-partition(());
        if deployablePartition is partitionClient:Partition {
            return {name: <string>deployablePartition.name, namespace: <string>deployablePartition.namespace, apiCount: deployablePartition.apiCount};
        } else {
            return e909022("Internal error occured ", deployablePartition);
        }
    }
}

