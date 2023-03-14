/*
 *  Copyright (c) 2023, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
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

package logging

import logging "github.com/wso2/apk/adapter/pkg/logging"

// Log (Error) severity level constants
const (
	BLOCKER  = "Blocker"
	CRITICAL = "Critical"
	MAJOR    = "Major"
	MINOR    = "Minor"
	TRIVIAL  = "Trivial"
	DEFAULT  = "Default"
)

// Error Log Internal Adapter(1100-1199) Constants
const (
	error1100 = 1100
	error1101 = 1101
	error1102 = 1102
	error1103 = 1103
	error1104 = 1104
)

// Error Log Internal API(1200-1299) Constants
const (
	error1200 = 1200
)

// Error Log Internal discovery(1400-1499) Config Constants
const (
	error1400 = 1400
	error1401 = 1401
	error1402 = 1402
	error1403 = 1403
	error1410 = 1410
	error1411 = 1411
	error1413 = 1413
	error1414 = 1414
)

// Error Log Internal XDS(1700-1799) Config Constants
const (
	error1700 = 1700
	error1701 = 1701
	error1702 = 1702
	error1703 = 1703
	error1704 = 1704
	error1705 = 1705
	error1706 = 1706
	error1707 = 1707
	error1709 = 1709
	error1710 = 1710
	error1711 = 1711
)

// Error Log Internal intercepter(1800-1899) Config Constants
const (
	error1800 = 1800
	error1801 = 1801
)

// Error Log Internal OASParser(2200-2299) Config Constants
const (
	error2200 = 2200
	error2201 = 2201
	error2204 = 2204
	error2205 = 2205
	error2206 = 2206
	error2207 = 2207
	error2208 = 2208
	error2209 = 2209
	error2210 = 2210
	error2211 = 2211
	error2212 = 2212
	error2231 = 2231
	error2234 = 2234
	error2235 = 2235
	error2236 = 2236
	error2237 = 2237
	error2238 = 2238
	error2239 = 2239
)

// Error Log Internal GRPC(2700-2799) Config Constants
const (
	error2700 = 2700
)

// Mapper used to keep error details for error logs
var Mapper = map[int]logging.ErrorDetails{
	error1100: {
		ErrorCode: error1100,
		Message:   "Failed to listen on port: %v, error: %v",
		Severity:  BLOCKER,
	},
	error1101: {
		ErrorCode: error1101,
		Message:   "Failed to start XDS GRPS server %s",
		Severity:  BLOCKER,
	},
	error1102: {
		ErrorCode: error1102,
		Message:   "Error reading the log configs. %v",
		Severity:  CRITICAL,
	},
	error1103: {
		ErrorCode: error1103,
		Message:   "Error while initializing authorization component, when intializing adapter REST API",
		Severity:  BLOCKER,
	},
	error1104: {
		ErrorCode: error1104,
		Message:   "Readiness probe is not set as local api artifacts processing has failed.",
		Severity:  CRITICAL,
	},
	error1200: {
		ErrorCode: error1200,
		Message:   "The provided port value for the REST Api Server :%v is not an integer. %v",
		Severity:  BLOCKER,
	},
	error1400: {
		ErrorCode: error1400,
		Message:   "Stream request for type %s on stream id: %d Error: %s",
		Severity:  CRITICAL,
	},
	error1401: {
		ErrorCode: error1401,
		Message:   "Stream request for type %s on stream id: %d, from node: %s, Error: %s",
		Severity:  CRITICAL,
	},
	error1402: {
		ErrorCode: error1402,
		Message:   "Consul syntax parse error %v",
		Severity:  CRITICAL,
	},
	error1403: {
		ErrorCode: error1403,
		Message:   "Internal Error while marshalling the upstream TLS Context. %v",
		Severity:  CRITICAL,
	},
	error1410: {
		ErrorCode: error1410,
		Message:   "Error undeploying API %v of Organization %v from environments %v",
		Severity:  MAJOR,
	},
	error1411: {
		ErrorCode: error1411,
		Message:   "Error extracting vhost from API identifier: %v for Organization %v. Ignore deploying the API",
		Severity:  MAJOR,
	},
	error1413: {
		ErrorCode: error1413,
		Message:   "Error creating new snapshot : %v",
		Severity:  MAJOR,
	},
	error1414: {
		ErrorCode: error1414,
		Message:   "Error while setting the snapshot : %v",
		Severity:  MAJOR,
	},
	error1700: {
		ErrorCode: error1700,
		Message:   "Error while connecting to the APK Management Server. %v",
		Severity:  BLOCKER,
	},
	error1701: {
		ErrorCode: error1701,
		Message:   "Error while starting APK Management application stream. %v",
		Severity:  BLOCKER,
	},
	error1702: {
		ErrorCode: error1702,
		Message:   "EOF is received from the APK Management Server application stream. %v",
		Severity:  CRITICAL,
	},
	error1703: {
		ErrorCode: error1703,
		Message:   "Failed to receive the discovery response from the APK Management Server application stream. %v",
		Severity:  CRITICAL,
	},
	error1704: {
		ErrorCode: error1704,
		Message:   "The APK Management Server application stream connection stopped: %v",
		Severity:  MINOR,
	},
	error1705: {
		ErrorCode: error1705,
		Message:   "Error while starting the APK Management Server: %v",
		Severity:  BLOCKER,
	},
	error1706: {
		ErrorCode: error1706,
		Message:   "Error while unmarshalling APK Management Server Application discovery response: %v",
		Severity:  MINOR,
	},
	error1707: {
		ErrorCode: error1707,
		Message:   "Error creating application: %v",
		Severity:  CRITICAL,
	},
	error1709: {
		ErrorCode: error1709,
		Message:   "Error updating application: %v",
		Severity:  CRITICAL,
	},
	error1710: {
		ErrorCode: error1710,
		Message:   "Error deleting application: %v",
		Severity:  CRITICAL,
	},
	error1711: {
		ErrorCode: error1711,
		Message:   "Error retrieving application: %v",
		Severity:  CRITICAL,
	},
	error1800: {
		ErrorCode: error1800,
		Message:   "error while parsing the interceptor template: %v",
		Severity:  CRITICAL,
	},
	error1801: {
		ErrorCode: error1801,
		Message:   "executing request interceptor template: %v",
		Severity:  CRITICAL,
	},
	error2200: {
		ErrorCode: error2200,
		Message:   "Error marsheling access log configs. %v",
		Severity:  CRITICAL,
	},
	error2201: {
		ErrorCode: error2201,
		Message:   "Error marshalling gRPC access log configs. %v",
		Severity:  CRITICAL,
	},
	error2204: {
		ErrorCode: error2204,
		Message:   "Operation policy validation failed for API %q in org %q:, policy %q: %v",
		Severity:  MINOR,
	},
	error2205: {
		ErrorCode: error2205,
		Message:   "Error parsing the operation policy definition %q into go template of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	error2206: {
		ErrorCode: error2206,
		Message:   "Error parsing operation policy definition %q of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	error2207: {
		ErrorCode: error2207,
		Message:   "Error parsing formalized operation policy definition %q into yaml of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	error2208: {
		ErrorCode: error2208,
		Message:   "API policy validation failed, policy: %q of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	error2209: {
		ErrorCode: error2209,
		Message:   "Error while JSON unmarshalling to find the API definition version. %s",
		Severity:  MINOR,
	},
	error2210: {
		ErrorCode: error2210,
		Message:   "AsyncAPI version %s is not supported.",
		Severity:  MINOR,
	},
	error2211: {
		ErrorCode: error2211,
		Message:   "API definition version is not defined.",
		Severity:  MINOR,
	},
	error2212: {
		ErrorCode: error2212,
		Message:   "Error adding request policy %s to operation %s of resource %s. %v",
		Severity:  MINOR,
	},
	error2231: {
		ErrorCode: error2231,
		Message:   "Error while creating routes for API %s %s for path: %s Error: %s",
		Severity:  MAJOR,
	},
	error2234: {
		ErrorCode: error2234,
		Message:   "Error occurred while creating the compression filter: %v",
		Severity:  MINOR,
	},
	error2235: {
		ErrorCode: error2235,
		Message:   "Error while parsing the gzip configuration value for the memory level: %v",
		Severity:  MINOR,
	},
	error2236: {
		ErrorCode: error2236,
		Message:   "Error while parsing the gzip configuration value for the window bits: %v",
		Severity:  MINOR,
	},
	error2237: {
		ErrorCode: error2237,
		Message:   "Error while parsing the gzip configuration value for the compression level: %v",
		Severity:  MINOR,
	},
	error2238: {
		ErrorCode: error2238,
		Message:   "Error while parsing the gzip configuration value for the chunk size: %v",
		Severity:  MINOR,
	},
	error2239: {
		ErrorCode: error2239,
		Message: "Error while adding resource level endpoints for %s:%v-%v. %v",
		Severity: MAJOR,
	},
	error2700: {
		ErrorCode: error2700,
		Message:   "Error while processing the private-public key pair : %v",
		Severity:  BLOCKER,
	},
}
