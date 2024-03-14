#!/bin/bash
# Copyright 2019 WSO2 Inc. (http://wso2.org)
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
# Generate JWT Tokens
# ----------------------------------------------------------------------------

script_dir=$(dirname "$0")
tokens_count=""
consumer_key=""

function usage() {
    echo ""
    echo "Usage: "
    echo "$0 -t <tokens_count> -c <consumer_key>"
    echo ""
    echo "-t: Count of the tokens."
    echo "-c: Consumer key relevant to the tokens."
    echo "-h: Display this help and exit."
    echo ""
}

while getopts "t:c:h" opt; do
    case "${opt}" in
    t)
        tokens_count=${OPTARG}
        ;;
    c)
        consumer_key=${OPTARG}
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
shift "$((OPTIND - 1))"

if [[ -z $tokens_count ]]; then
    echo "Please provide the Token count"
    exit 1
fi

if [[ -z $consumer_key ]]; then
    echo "Please provide the consumer key"
    exit 1
fi

mkdir -p "$script_dir/target"
tokens_file="$script_dir/target/jwt-tokens-${tokens_count}.csv"

if [[ -f $tokens_file ]]; then
    mv $tokens_file /tmp
fi

generate_tokens_command="java -Xms128m -Xmx128m -jar $script_dir/jwt-generator-0.1.1-SNAPSHOT.jar \
        --key-store-file $script_dir/keystore.jks --consumer-key $consumer_key\
        --tokens-count $tokens_count --output-file $tokens_file --issuer "https://idp1.com""
echo "Generating Tokens: $generate_tokens_command"
$generate_tokens_command
