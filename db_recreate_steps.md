# Steps to recreate the DB

This step is necessary when there are changes to the DB Schema in [initdb-conf.yaml](https://github.com/wso2/apk/blob/main/helm-charts/templates/control-plane/postgres/initdb-conf.yaml) 
file and need to recreate the DB by removing all previous data and schema.

You can follow the following steps to achieve this requirement.
 
1. Helm Uninstall

> helm uninstall <HELM_RELEASE> -n apk

2. Delete Existing PVC and PV.

> kubectl get pvc -n apk

> kubectl delete pvc <PVC_NAME> -n apk

> kubectl get pv -n apk

> kubectl delete pv <PV_NAME> -n apk

3. Delete Standard Storage Claim

> kubectl get sc

> kubectl delete sc standard

4. Stop K8s Cluster.

> minikube stop

5. Start K8s Cluster.

> minikube start --container-runtime=docker

6. Helm install

> helm install <HELM_RELEASE> . -n apk

By following above steps, a new DB will be created using the new schema provided through helm.