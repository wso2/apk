import wso2/apk_common_lib as commons;
import config_deployer_service.model;

public type PartitionResolver isolated object {
    isolated function getAvailablePartitionForAPI(string id, string organization) returns model:Partition|commons:APKError?;
    isolated function getDeployablePartition() returns model:Partition|commons:APKError;
};
