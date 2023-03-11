set -e
#This is a sample script to build APK. When you build project please run this script in project root level.
#All relative paths etc designed from root directory. Users can customize this as per demand. Ex: If you wish to 
#build and run runtime domain service then can build it alone and do deployment.
current_dir=$PWD;
cd $current_dir;
cd common-bal-libs/apk-common-lib;./gradlew build;
cd $current_dir;
cd backoffice/backoffice-domain-service;./gradlew build;
cd $current_dir;
cd runtime/runtime-domain-service;./gradlew build;
cd $current_dir;
cd admin/admin-domain-service;./gradlew build;
cd $current_dir;
cd devportal/devportal-domain-service;./gradlew build;
cd $current_dir;
cd adapter;./gradlew build;
cd $current_dir;
cd management-server;./gradlew build;
cd $current_dir;
cd gateway/router;./gradlew build;
cd $current_dir;
cd gateway/enforcer;./gradlew build;
cd $current_dir;
cd idp/idp-domain-service;./gradlew build;
cd $current_dir;
cd idp/idp-ui;./gradlew build;
cd $current_dir;
cd ratelimiter;./gradlew build;
