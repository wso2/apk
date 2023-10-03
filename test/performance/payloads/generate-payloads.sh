#!/bin/bash -e
# Copyright 2017 WSO2 Inc. (http://wso2.org)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# ----------------------------------------------------------------------------
# Generate JSON Payloads
# ----------------------------------------------------------------------------

script_dir=$(dirname "$0")
payload_type=""
declare -a payloads

function usage() {
    echo ""
    echo "Usage: "
    echo "$0 [-p <payload_type>] [-s <payload_size>]"
    echo ""
    echo "-p: The Payload Type."
    echo "-s: The Payload Size. You can give multiple payload sizes."
    echo "-h: Display this help and exit."
    echo ""
}

while getopts "p:s:h" opts; do
    case $opts in
    p)
        payload_type=${OPTARG}
        ;;
    s)
        payloads+=("${OPTARG}")
        ;;
    h)
        usage
        exit 0
        ;;
    \?)
        usage
        exit 1
        ;;
    esac
done

if [[ -z $payload_type ]]; then
    payload_type="simple"
fi

if [[ ${#payloads[@]} -eq 0 ]]; then
    payloads=("50 1024 10240 102400")
fi

for s in ${payloads[*]}; do
    echo "Generating ${s}B file"
    java -jar $script_dir/payload-generator-0.4.6-SNAPSHOT.jar --size $s --payload-type ${payload_type}
done
