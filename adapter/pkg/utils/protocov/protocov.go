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

package protocov

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	APIPrefix = "type.googleapis.com/"
)

var (
	marshalOpts = proto.MarshalOptions{}
)

func ToAnyWithError(msg proto.Message) (*anypb.Any, error) {
	b, err := marshalOpts.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &anypb.Any{
		TypeUrl: APIPrefix + string(msg.ProtoReflect().Descriptor().FullName()),
		Value:   b,
	}, nil
}

func ToAny(msg proto.Message) *anypb.Any {
	res, err := ToAnyWithError(msg)
	if err != nil {
		return nil
	}
	return res
}
