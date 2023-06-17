import wso2/apk_common_lib as commons;
import config_deployer_service.model;

public isolated class SinglePartitionResolver {
    *PartitionResolver;

    isolated function getAvailablePartitionForAPI(string id, string organization) returns model:Partition|commons:APKError? {
        model:API? k8sAPIByNameAndNamespace = check getK8sAPIByNameAndNamespace(id, currentNameSpace);
        if k8sAPIByNameAndNamespace is model:API {
            model:Partition partition = {name: DEFAULT_PARTITION, namespace: currentNameSpace};
            return partition;
        } else {
            return ();
        }
    }
    isolated function getDeployablePartition() returns model:Partition|commons:APKError {
        return {name: DEFAULT_PARTITION, namespace: currentNameSpace};
    }
}
