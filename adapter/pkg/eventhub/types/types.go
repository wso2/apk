/*
 *  Copyright (c) 2022, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 */

package types

// API for struct Api
type API struct {
	APIID            int    `json:"apiId"`
	UUID             string `json:"uuid"`
	Provider         string `json:"provider"`
	Name             string `json:"name"`
	Version          string `json:"version"`
	BasePath         string `json:"basePath"`
	Policy           string `json:"policy"`
	APIType          string `json:"apiType"`
	IsDefaultVersion bool   `json:"isDefaultVersion"`
	APIStatus        string `json:"status"`
	TenantID         int32  `json:"tenanId,omitempty"`
	TenantDomain     string `json:"tenanDomain,omitempty"`
	TimeStamp        int64  `json:"timeStamp,omitempty"`
}

// APIList for struct ApiList
type APIList struct {
	List []API `json:"list"`
}

// ApplicationPolicy for struct ApplicationPolicy
type ApplicationPolicy struct {
	ID        int32  `json:"id"`
	TenantID  int32  `json:"tenantId"`
	Name      string `json:"name"`
	QuotaType string `json:"quotaType"`
}

// ApplicationPolicyList for struct list of ApplicationPolicy
type ApplicationPolicyList struct {
	List []ApplicationPolicy `json:"list"`
}

// SubscriptionPolicy for struct list of SubscriptionPolicy
type SubscriptionPolicy struct {
	ID                   int32  `json:"id" json:"policyId"`
	TenantID             int32  `json:"tenantId"`
	Name                 string `json:"name"`
	QuotaType            string `json:"quotaType"`
	GraphQLMaxComplexity int32  `json:"graphQLMaxComplexity"`
	GraphQLMaxDepth      int32  `json:"graphQLMaxDepth"`
	RateLimitCount       int32  `json:"rateLimitCount"`
	RateLimitTimeUnit    string `json:"rateLimitTimeUnit"`
	StopOnQuotaReach     bool   `json:"stopOnQuotaReach"`
	TenantDomain         string `json:"tenanDomain,omitempty"`
	TimeStamp            int64  `json:"timeStamp,omitempty"`
}

// SubscriptionPolicyList for struct list of SubscriptionPolicy
type SubscriptionPolicyList struct {
	List []SubscriptionPolicy `json:"list"`
}

// APIPolicy for struct policy Info events
type APIPolicy struct {
	PolicyID                 string `json:"policyId"`
	PolicyName               string `json:"policyName"`
	QuotaType                string `json:"quotaType"`
	PolicyType               string `json:"policyType"`
	AddedConditionGroupIDs   string `json:"addedConditionGroupIDs"`
	DeletedConditionGroupIDs string `json:"deletedConditionGroupIDs"`
	TimeStamp                int64  `json:"timeStamp,omitempty"`
}

// Scope for struct Scope
type Scope struct {
	Name            string `json:"name"`
	DisplayName     string `json:"displayName"`
	ApplicationName string `json:"description"`
}

// ScopeList for struct list of Scope
type ScopeList struct {
	List []Scope `json:"list"`
}

// KeyManagerList for struct list of KeyManager
type KeyManagerList struct {
	KeyManagers []KeyManager `json:"KeyManager"`
}

// KeyManager for struct
type KeyManager struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Enabled       bool   `json:"enabled"`
	TenantDomain  string `json:"tenantDomain,omitempty"`
	Configuration map[string]interface{}
	// Configuration KeyManagerConfig `json:"configuration"`
}

// KeyManagerConfig for struct Configuration map[string]interface{} `json:"value"`
type KeyManagerConfig struct {
	TokenFormatString          string   `json:"token_format_string"`
	ServerURL                  string   `json:"ServerURL"`
	ValidationEnable           bool     `json:"validation_enable"`
	ClaimMappings              []Claim  `json:"Claim"`
	GrantTypes                 []string `json:"grant_types"`
	EncryptPersistedTokens     bool     `json:"OAuthConfigurations.EncryptPersistedTokens"`
	EnableOauthAppCreation     bool     `json:"enable_oauth_app_creation"`
	ValidityPeriod             string   `json:"VALIDITY_PERIOD"`
	EnableTokenGeneration      bool     `json:"enable_token_generation"`
	Issuer                     string   `json:"issuer"`
	EnableMapOauthConsumerApps bool     `json:"enable_map_oauth_consumer_apps"`
	EnableTokenHash            bool     `json:"enable_token_hash"`
	SelfValidateJwt            bool     `json:"self_validate_jwt"`
	RevokeEndpoint             string   `json:"revoke_endpoint"`
	EnableTokenEncryption      bool     `json:"enable_token_encryption"`
	RevokeURL                  string   `json:"RevokeURL"`
	TokenURL                   string   `json:"TokenURL,token_endpoint"`
	CertificateType            string   `json:"certificate_type"`
	CertificateValue           string   `json:"certificate_value"`
}

// Claim for struct
type Claim struct {
	RemoteClaim string `json:"remoteClaim"`
	LocalClaim  string `json:"localClaim"`
}
