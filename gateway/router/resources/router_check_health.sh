#!/bin/bash
# --------------------------------------------------------------------
# Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com) All Rights Reserved.
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.
# -----------------------------------------------------------------------

status_code=$(curl --write-out %{http_code} --silent --output /dev/null https://localhost:9095/health -H 'Authorization: Basic YWRtaW46YWRtaW4=' -k -v)

if [[ "$status_code" -ne 200 ]] ; then
  echo "Health check status changed to $status_code"
else
  echo "Health check status changed to $status_code"
  exit 0
fi
