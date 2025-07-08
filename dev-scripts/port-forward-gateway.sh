#!/bin/bash

#
# Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com). All Rights Reserved.
#
# This software is the property of WSO2 LLC. and its suppliers, if any.
# Dissemination of any information or reproduction of any material contained
# herein in any form is strictly forbidden, unless permitted by WSO2 expressly.
# You may not alter or remove any copyright or other notice from copies of this content.
#

kubectl port-forward -n apk-egress-gateway deploy/apk-eg-wso2-apk-gateway-runtime-deployment 9000:9000 &
GW_ADMIN_PID=$!

kubectl port-forward -n apk-egress-gateway deploy/apk-eg-wso2-apk-gateway-runtime-deployment 9095:9095 &
GW_HTTPS_LISTENER_PID=$!

kubectl port-forward -n apk-egress-gateway deploy/apk-eg-wso2-apk-gateway-runtime-deployment 9080:9080 &
GW_HTTP_LISTENER_PID=$!

cleanup() {
    kill -9 $GW_ADMIN_PID $GW_HTTPS_LISTENER_PID $GW_HTTP_LISTENER_PID
}

trap cleanup SIGINT

wait
