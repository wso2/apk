//  Copyright (c) 2021, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
//
//  WSO2 LLC. licenses this file to you under the Apache License,
//  Version 2.0 (the "License"); you may not use this file except
//  in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an
//  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
//  KIND, either express or implied.  See the License for the
//  specific language governing permissions and limitations
//  under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.13.0
// source: wso2/discovery/api/api.proto

package api

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// API config model
type Api struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                     string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title                  string      `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Version                string      `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	ApiType                string      `protobuf:"bytes,4,opt,name=apiType,proto3" json:"apiType,omitempty"`
	DisableAuthentications bool        `protobuf:"varint,5,opt,name=disableAuthentications,proto3" json:"disableAuthentications,omitempty"`
	DisableScopes          bool        `protobuf:"varint,6,opt,name=disableScopes,proto3" json:"disableScopes,omitempty"`
	EnvType                string      `protobuf:"bytes,7,opt,name=envType,proto3" json:"envType,omitempty"`
	Resources              []*Resource `protobuf:"bytes,8,rep,name=resources,proto3" json:"resources,omitempty"`
	BasePath               string      `protobuf:"bytes,9,opt,name=basePath,proto3" json:"basePath,omitempty"`
	Tier                   string      `protobuf:"bytes,10,opt,name=tier,proto3" json:"tier,omitempty"`
	ApiLifeCycleState      string      `protobuf:"bytes,11,opt,name=apiLifeCycleState,proto3" json:"apiLifeCycleState,omitempty"`
	Vhost                  string      `protobuf:"bytes,12,opt,name=vhost,proto3" json:"vhost,omitempty"`
	OrganizationId         string      `protobuf:"bytes,13,opt,name=organizationId,proto3" json:"organizationId,omitempty"`
	// bool isMockedApi = 18;
	ClientCertificates  []*Certificate  `protobuf:"bytes,14,rep,name=clientCertificates,proto3" json:"clientCertificates,omitempty"`
	MutualSSL           string          `protobuf:"bytes,15,opt,name=mutualSSL,proto3" json:"mutualSSL,omitempty"`
	ApplicationSecurity map[string]bool `protobuf:"bytes,16,rep,name=applicationSecurity,proto3" json:"applicationSecurity,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	TransportSecurity   bool            `protobuf:"varint,17,opt,name=transportSecurity,proto3" json:"transportSecurity,omitempty"`
	HttpRouteIDs        []string        `protobuf:"bytes,19,rep,name=httpRouteIDs,proto3" json:"httpRouteIDs,omitempty"`
	/// string graphQLSchema = 22;
	GraphqlComplexityInfo  []*GraphqlComplexity    `protobuf:"bytes,23,rep,name=graphqlComplexityInfo,proto3" json:"graphqlComplexityInfo,omitempty"`
	SystemAPI              bool                    `protobuf:"varint,24,opt,name=systemAPI,proto3" json:"systemAPI,omitempty"`
	BackendJWTTokenInfo    *BackendJWTTokenInfo    `protobuf:"bytes,25,opt,name=backendJWTTokenInfo,proto3" json:"backendJWTTokenInfo,omitempty"`
	ApiDefinitionFile      []byte                  `protobuf:"bytes,26,opt,name=apiDefinitionFile,proto3" json:"apiDefinitionFile,omitempty"`
	Environment            string                  `protobuf:"bytes,27,opt,name=environment,proto3" json:"environment,omitempty"`
	SubscriptionValidation bool                    `protobuf:"varint,28,opt,name=subscriptionValidation,proto3" json:"subscriptionValidation,omitempty"`
	Endpoints              *EndpointCluster        `protobuf:"bytes,29,opt,name=endpoints,proto3" json:"endpoints,omitempty"`
	EndpointSecurity       []*SecurityInfo         `protobuf:"bytes,30,rep,name=endpointSecurity,proto3" json:"endpointSecurity,omitempty"`
	Aiprovider             *AIProvider             `protobuf:"bytes,31,opt,name=aiprovider,proto3" json:"aiprovider,omitempty"`
	AiModelBasedRoundRobin *AIModelBasedRoundRobin `protobuf:"bytes,32,opt,name=aiModelBasedRoundRobin,proto3" json:"aiModelBasedRoundRobin,omitempty"`
	ApiDefinitionPath      string                  `protobuf:"bytes,33,opt,name=apiDefinitionPath,proto3" json:"apiDefinitionPath,omitempty"`
}

func (x *Api) Reset() {
	*x = Api{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wso2_discovery_api_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Api) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Api) ProtoMessage() {}

func (x *Api) ProtoReflect() protoreflect.Message {
	mi := &file_wso2_discovery_api_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Api.ProtoReflect.Descriptor instead.
func (*Api) Descriptor() ([]byte, []int) {
	return file_wso2_discovery_api_api_proto_rawDescGZIP(), []int{0}
}

func (x *Api) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Api) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Api) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Api) GetApiType() string {
	if x != nil {
		return x.ApiType
	}
	return ""
}

func (x *Api) GetDisableAuthentications() bool {
	if x != nil {
		return x.DisableAuthentications
	}
	return false
}

func (x *Api) GetDisableScopes() bool {
	if x != nil {
		return x.DisableScopes
	}
	return false
}

func (x *Api) GetEnvType() string {
	if x != nil {
		return x.EnvType
	}
	return ""
}

func (x *Api) GetResources() []*Resource {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *Api) GetBasePath() string {
	if x != nil {
		return x.BasePath
	}
	return ""
}

func (x *Api) GetTier() string {
	if x != nil {
		return x.Tier
	}
	return ""
}

func (x *Api) GetApiLifeCycleState() string {
	if x != nil {
		return x.ApiLifeCycleState
	}
	return ""
}

func (x *Api) GetVhost() string {
	if x != nil {
		return x.Vhost
	}
	return ""
}

func (x *Api) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *Api) GetClientCertificates() []*Certificate {
	if x != nil {
		return x.ClientCertificates
	}
	return nil
}

func (x *Api) GetMutualSSL() string {
	if x != nil {
		return x.MutualSSL
	}
	return ""
}

func (x *Api) GetApplicationSecurity() map[string]bool {
	if x != nil {
		return x.ApplicationSecurity
	}
	return nil
}

func (x *Api) GetTransportSecurity() bool {
	if x != nil {
		return x.TransportSecurity
	}
	return false
}

func (x *Api) GetHttpRouteIDs() []string {
	if x != nil {
		return x.HttpRouteIDs
	}
	return nil
}

func (x *Api) GetGraphqlComplexityInfo() []*GraphqlComplexity {
	if x != nil {
		return x.GraphqlComplexityInfo
	}
	return nil
}

func (x *Api) GetSystemAPI() bool {
	if x != nil {
		return x.SystemAPI
	}
	return false
}

func (x *Api) GetBackendJWTTokenInfo() *BackendJWTTokenInfo {
	if x != nil {
		return x.BackendJWTTokenInfo
	}
	return nil
}

func (x *Api) GetApiDefinitionFile() []byte {
	if x != nil {
		return x.ApiDefinitionFile
	}
	return nil
}

func (x *Api) GetEnvironment() string {
	if x != nil {
		return x.Environment
	}
	return ""
}

func (x *Api) GetSubscriptionValidation() bool {
	if x != nil {
		return x.SubscriptionValidation
	}
	return false
}

func (x *Api) GetEndpoints() *EndpointCluster {
	if x != nil {
		return x.Endpoints
	}
	return nil
}

func (x *Api) GetEndpointSecurity() []*SecurityInfo {
	if x != nil {
		return x.EndpointSecurity
	}
	return nil
}

func (x *Api) GetAiprovider() *AIProvider {
	if x != nil {
		return x.Aiprovider
	}
	return nil
}

func (x *Api) GetAiModelBasedRoundRobin() *AIModelBasedRoundRobin {
	if x != nil {
		return x.AiModelBasedRoundRobin
	}
	return nil
}

func (x *Api) GetApiDefinitionPath() string {
	if x != nil {
		return x.ApiDefinitionPath
	}
	return ""
}

var File_wso2_discovery_api_api_proto protoreflect.FileDescriptor

var file_wso2_discovery_api_api_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12,
	0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61,
	0x70, 0x69, 0x1a, 0x21, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2c, 0x77, 0x73, 0x6f,
	0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x29, 0x77, 0x73, 0x6f, 0x32, 0x2f,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x6e,
	0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x77, 0x73,
	0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24,
	0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x61, 0x69, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x33, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x69, 0x5f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x5f, 0x62, 0x61, 0x73, 0x65, 0x64, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x5f, 0x72, 0x6f,
	0x62, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdd, 0x0b, 0x0a, 0x03, 0x41, 0x70,
	0x69, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x70, 0x69, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x61, 0x70, 0x69, 0x54, 0x79, 0x70, 0x65, 0x12, 0x36, 0x0a, 0x16, 0x64,
	0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x64, 0x69, 0x73,
	0x61, 0x62, 0x6c, 0x65, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x63,
	0x6f, 0x70, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x64, 0x69, 0x73, 0x61,
	0x62, 0x6c, 0x65, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x76,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x76, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x3a, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73,
	0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69,
	0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12,
	0x1a, 0x0a, 0x08, 0x62, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x62, 0x61, 0x73, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x69, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x65, 0x72, 0x12,
	0x2c, 0x0a, 0x11, 0x61, 0x70, 0x69, 0x4c, 0x69, 0x66, 0x65, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x61, 0x70, 0x69, 0x4c,
	0x69, 0x66, 0x65, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x68,
	0x6f, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x4f, 0x0a, 0x12, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x73, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x65, 0x72,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x52, 0x12, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x6d, 0x75, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x53, 0x4c, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x6d, 0x75, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x53, 0x4c, 0x12, 0x62, 0x0a, 0x13, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x10, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x69,
	0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x75,
	0x72, 0x69, 0x74, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x13, 0x61, 0x70, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12, 0x2c,
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x63, 0x75, 0x72,
	0x69, 0x74, 0x79, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73,
	0x70, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12, 0x22, 0x0a, 0x0c,
	0x68, 0x74, 0x74, 0x70, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x49, 0x44, 0x73, 0x18, 0x13, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0c, 0x68, 0x74, 0x74, 0x70, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x49, 0x44, 0x73,
	0x12, 0x5b, 0x0a, 0x15, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x43, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x78, 0x69, 0x74, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x17, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x43, 0x6f, 0x6d, 0x70,
	0x6c, 0x65, 0x78, 0x69, 0x74, 0x79, 0x52, 0x15, 0x67, 0x72, 0x61, 0x70, 0x68, 0x71, 0x6c, 0x43,
	0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x78, 0x69, 0x74, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x41, 0x50, 0x49, 0x18, 0x18, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x41, 0x50, 0x49, 0x12, 0x59, 0x0a, 0x13, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x18, 0x19, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x42, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x13, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4a, 0x57, 0x54, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2c, 0x0a, 0x11, 0x61, 0x70, 0x69, 0x44, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x18, 0x1a, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x11, 0x61, 0x70, 0x69, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x69, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x36, 0x0a, 0x16, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x1c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x41,
	0x0a, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x1d, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x43,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x73, 0x12, 0x4c, 0x0a, 0x10, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x53, 0x65, 0x63,
	0x75, 0x72, 0x69, 0x74, 0x79, 0x18, 0x1e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x77, 0x73,
	0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x10, 0x65,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x12,
	0x3e, 0x0a, 0x0a, 0x61, 0x69, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x1f, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x49, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x52, 0x0a, 0x61, 0x69, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12,
	0x62, 0x0a, 0x16, 0x61, 0x69, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x61, 0x73, 0x65, 0x64, 0x52,
	0x6f, 0x75, 0x6e, 0x64, 0x52, 0x6f, 0x62, 0x69, 0x6e, 0x18, 0x20, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x49, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x61, 0x73, 0x65,
	0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x6f, 0x62, 0x69, 0x6e, 0x52, 0x16, 0x61, 0x69, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x42, 0x61, 0x73, 0x65, 0x64, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x6f,
	0x62, 0x69, 0x6e, 0x12, 0x2c, 0x0a, 0x11, 0x61, 0x70, 0x69, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x74, 0x68, 0x18, 0x21, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11,
	0x61, 0x70, 0x69, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x61, 0x74,
	0x68, 0x1a, 0x46, 0x0a, 0x18, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x70, 0x0a, 0x23, 0x6f, 0x72, 0x67,
	0x2e, 0x77, 0x73, 0x6f, 0x32, 0x2e, 0x61, 0x70, 0x6b, 0x2e, 0x65, 0x6e, 0x66, 0x6f, 0x72, 0x63,
	0x65, 0x72, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x2e, 0x61, 0x70, 0x69,
	0x42, 0x08, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3d, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72,
	0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x77, 0x73, 0x6f, 0x32, 0x2f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76,
	0x65, 0x72, 0x79, 0x2f, 0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_wso2_discovery_api_api_proto_rawDescOnce sync.Once
	file_wso2_discovery_api_api_proto_rawDescData = file_wso2_discovery_api_api_proto_rawDesc
)

func file_wso2_discovery_api_api_proto_rawDescGZIP() []byte {
	file_wso2_discovery_api_api_proto_rawDescOnce.Do(func() {
		file_wso2_discovery_api_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_wso2_discovery_api_api_proto_rawDescData)
	})
	return file_wso2_discovery_api_api_proto_rawDescData
}

var file_wso2_discovery_api_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_wso2_discovery_api_api_proto_goTypes = []interface{}{
	(*Api)(nil),                    // 0: wso2.discovery.api.Api
	nil,                            // 1: wso2.discovery.api.Api.ApplicationSecurityEntry
	(*Resource)(nil),               // 2: wso2.discovery.api.Resource
	(*Certificate)(nil),            // 3: wso2.discovery.api.Certificate
	(*GraphqlComplexity)(nil),      // 4: wso2.discovery.api.GraphqlComplexity
	(*BackendJWTTokenInfo)(nil),    // 5: wso2.discovery.api.BackendJWTTokenInfo
	(*EndpointCluster)(nil),        // 6: wso2.discovery.api.EndpointCluster
	(*SecurityInfo)(nil),           // 7: wso2.discovery.api.SecurityInfo
	(*AIProvider)(nil),             // 8: wso2.discovery.api.AIProvider
	(*AIModelBasedRoundRobin)(nil), // 9: wso2.discovery.api.AIModelBasedRoundRobin
}
var file_wso2_discovery_api_api_proto_depIdxs = []int32{
	2, // 0: wso2.discovery.api.Api.resources:type_name -> wso2.discovery.api.Resource
	3, // 1: wso2.discovery.api.Api.clientCertificates:type_name -> wso2.discovery.api.Certificate
	1, // 2: wso2.discovery.api.Api.applicationSecurity:type_name -> wso2.discovery.api.Api.ApplicationSecurityEntry
	4, // 3: wso2.discovery.api.Api.graphqlComplexityInfo:type_name -> wso2.discovery.api.GraphqlComplexity
	5, // 4: wso2.discovery.api.Api.backendJWTTokenInfo:type_name -> wso2.discovery.api.BackendJWTTokenInfo
	6, // 5: wso2.discovery.api.Api.endpoints:type_name -> wso2.discovery.api.EndpointCluster
	7, // 6: wso2.discovery.api.Api.endpointSecurity:type_name -> wso2.discovery.api.SecurityInfo
	8, // 7: wso2.discovery.api.Api.aiprovider:type_name -> wso2.discovery.api.AIProvider
	9, // 8: wso2.discovery.api.Api.aiModelBasedRoundRobin:type_name -> wso2.discovery.api.AIModelBasedRoundRobin
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_wso2_discovery_api_api_proto_init() }
func file_wso2_discovery_api_api_proto_init() {
	if File_wso2_discovery_api_api_proto != nil {
		return
	}
	file_wso2_discovery_api_Resource_proto_init()
	file_wso2_discovery_api_Certificate_proto_init()
	file_wso2_discovery_api_BackendJWTTokenInfo_proto_init()
	file_wso2_discovery_api_endpoint_cluster_proto_init()
	file_wso2_discovery_api_security_info_proto_init()
	file_wso2_discovery_api_graphql_proto_init()
	file_wso2_discovery_api_ai_provider_proto_init()
	file_wso2_discovery_api_ai_model_based_round_robin_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_wso2_discovery_api_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Api); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wso2_discovery_api_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wso2_discovery_api_api_proto_goTypes,
		DependencyIndexes: file_wso2_discovery_api_api_proto_depIdxs,
		MessageInfos:      file_wso2_discovery_api_api_proto_msgTypes,
	}.Build()
	File_wso2_discovery_api_api_proto = out.File
	file_wso2_discovery_api_api_proto_rawDesc = nil
	file_wso2_discovery_api_api_proto_goTypes = nil
	file_wso2_discovery_api_api_proto_depIdxs = nil
}
