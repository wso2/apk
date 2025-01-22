/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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
 
package authenticator

import "fmt"

// JWTAuthenticator implements Authenticator for JWT tokens.
type JWTAuthenticator struct{}

// CanAuthenticate checks if the data contains a JWT token.
func (j JWTAuthenticator) CanAuthenticate(data map[string]string) bool {
	_, exists := data["jwtToken"]
	return exists
}

// Authenticate validates the JWT token.
func (j JWTAuthenticator) Authenticate(data map[string]string) (bool, error) {
	token, exists := data["jwtToken"]
	if !exists {
		return false, fmt.Errorf("no JWT token found")
	}
	// Add actual JWT validation logic here.
	fmt.Println("Validating JWT:", token)
	return true, nil
}
