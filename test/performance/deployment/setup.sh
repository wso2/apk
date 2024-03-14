#bin/sh
wget "https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.22%2B7/OpenJDK11U-jdk_x64_linux_hotspot_11.0.22_7.tar.gz"
tar -xvzf OpenJDK11U-jdk_x64_linux_hotspot_11.0.22_7.tar.gz
export JAVA_HOME=$(pwd)/jdk-11.0.22+7>> ~/.bashrc
export PATH=$JAVA_HOME/bin:$PATH>> ~/.bashrc
source ~/.bashrc
wget https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-5.5.tgz
tar -xvzf apache-jmeter-5.5.tgz