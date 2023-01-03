set -e
#This is a sample script to build APK. When you build project please run this script in project root level.
#All relative paths etc designed from root directory. Users can customize this as per demand. Ex: If you wish to 
#build and run runtime domain service then can build it alone and do deployment.
current_dir=$PWD;
cd $current_dir;
cd backoffice/backoffice-domain-service;gradle build;
cd $current_dir;
cd backoffice/backoffice-internal-domain-service;gradle build;
cd $current_dir;
cd runtime/runtime-domain-service;gradle build;
cd $current_dir;
cd admin/admin-domain-service;gradle build;
cd $current_dir;
cd devportal/devportal-domain-service;gradle build;
cd $current_dir;
cd adapter;gradle build;
cd $current_dir;
cd management-server;gradle build;