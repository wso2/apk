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

package stringutils

import (
	"fmt"
	"strings"
)

// DefaultMaskLength is the mask length of masked string
const DefaultMaskLength = 10

// MaskString takes string as argument and keep a visibleLength number of charators and mask
// maskCharacter is used to generate the masked part of the string
// It can mask the charators from either left or right hand side of the string depending on maskRight boolean
// It returns the masked string
func MaskString(str string, visibleLength int, maskCharacter string, maskRight bool) string {
	stringLength := len(str)

	if 0 <= visibleLength && visibleLength < stringLength {
		maskLength := DefaultMaskLength
		mask := strings.Repeat(maskCharacter, maskLength)
		if maskRight {
			return str[:visibleLength] + mask
		}
		return mask + str[stringLength-visibleLength:]
	}

	return strings.Repeat(maskCharacter, stringLength)
}

// MaskToken masks the token with "*" character and keeps last four charaters visible
func MaskToken(token string) string {
	return MaskString(token, 4, "*", false)
}

// StringInSlice checks whether a given string is included in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ToStringArray returns the vHosts related to an API.
func ToStringArray[T any](array []T) []string {
	var vHosts []string
	for _, stringVal := range array {
		vHosts = append(vHosts, fmt.Sprintf("%v", stringVal))
	}
	return vHosts
}

// RemoveString returns a newly created []string that contains all items from slice that
// are not equal to s.
func RemoveString(slice []string, s string) []string {
	var newSlice []string
	for _, item := range slice {
		if item == s {
			continue
		}
		newSlice = append(newSlice, item)
	}
	return newSlice
}
