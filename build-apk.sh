#This is a sample script to build APK. When you build project please run this script in project root level.
#All relative paths etc designed from root directory. Users can customize this as per demand. Ex: If you wish to 
#build and run runtime domain service then can build it alone and do deployment.
declare -i x=0;
git pull
current_dir=$PWD;
#Component build sample. You can build any component with this
cd apkbase;gradle build;
cd common-java-libs;gradle build;
cd $current_dir;
cd common-go-libs/;gradle build;
cd $current_dir;
cd backoffice/backoffice-domain-service;gradle build;
cd $current_dir;
cd runtime/runtime-domain-service;gradle build;
cd $current_dir;
cd helm-charts; 
#Undeploy already deployed setup and deploy again
kubectl get pods -n apk
helm uninstall apk-test -n apk
#If you are use minikube implementation then you need to load images with below command.
minikube image load backoffice_service:0.1.0-SNAPSHOT  
#minikube image load adapter:0.0.1-SNAPSHOT   
#minikube image load management-server:0.0.1-SNAPSHOT  
minikube image load runtime-api:0.0.1-SNAPSHOT   
#minikube image load devportal_service:0.1.0-SNAPSHOT  
#minikube image load admin-service:0.1.0-SNAPSHOT  
minikube image load apkbase:0.0.1-SNAPSHOT
#Waiting for deployment to terminate.
while test $x -eq 0; do
getpods_output=$(kubectl get pods -n apk 2>&1)
if [[ $getpods_output ==  *"No resources found in apk namespace."* ]]; then
   x=1;
   echo "Deployment terminated successfully..........."
else
   echo "Terminating Cluster..............."
   sleep 5
fi
done
#Helm install APK with all components
helm dependency build
helm install apk-test -n apk .
#List pods created for deployment.
kubectl get pods -n apk