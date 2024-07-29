/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  This file contains code derived from Envoy Gateway,
 *  https://github.com/envoyproxy/gateway
 *  and is provided here subject to the following:
 *  Copyright Project Envoy Gateway Authors
 *
 */

package utils

import (
	"bytes"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

func MarshalResourcesToJSON(resources []types.Resource) ([]byte, error) {
	msgs := make([]proto.Message, 0)
	for _, resource := range resources {
		msgs = append(msgs, resource.(proto.Message))
	}
	var buffer bytes.Buffer
	buffer.WriteByte('[')
	for idx, msg := range msgs {
		if idx != 0 {
			buffer.WriteByte(',')
		}
		b, err := protojson.Marshal(msg)
		if err != nil {
			return nil, err
		}
		buffer.Write(b)
	}
	buffer.WriteByte(']')
	return buffer.Bytes(), nil
}

// ResourcesToYAMLString converts xDS Resource types into YAML string
func ResourcesToYAMLString(resources []types.Resource) (string, error) {
	jsonBytes, err := MarshalResourcesToJSON(resources)
	if err != nil {
		return "", err
	}
	data, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
