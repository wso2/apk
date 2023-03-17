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
	Error1100 = 1100
	Error1101 = 1101
	Error1102 = 1102
	Error1103 = 1103
	Error1104 = 1104
)

// Error Log Internal API(1200-1299) Constants
const (
	Error1200 = 1200
)

// Error Log Internal discovery(1400-1499) Config Constants
const (
	Error1400 = 1400
	Error1401 = 1401
	Error1402 = 1402
	Error1403 = 1403
	Error1410 = 1410
	Error1411 = 1411
	Error1413 = 1413
	Error1414 = 1414
)

// Error Log Internal XDS(1700-1799) Config Constants
const (
	Error1700 = 1700
	Error1701 = 1701
	Error1702 = 1702
	Error1703 = 1703
	Error1704 = 1704
	Error1705 = 1705
	Error1706 = 1706
	Error1707 = 1707
	Error1709 = 1709
	Error1710 = 1710
	Error1711 = 1711
)

// Error Log Internal intercepter(1800-1899) Config Constants
const (
	Error1800 = 1800
	Error1801 = 1801
)

// Error Log Internal OASParser(2200-2299) Config Constants
const (
	Error2200 = 2200
	Error2201 = 2201
	Error2204 = 2204
	Error2205 = 2205
	Error2206 = 2206
	Error2207 = 2207
	Error2208 = 2208
	Error2209 = 2209
	Error2210 = 2210
	Error2211 = 2211
	Error2212 = 2212
	Error2231 = 2231
	Error2234 = 2234
	Error2235 = 2235
	Error2236 = 2236
	Error2237 = 2237
	Error2238 = 2238
	Error2239 = 2239
)

// Error Log Internal GRPC(2700-2799) Config Constants
const (
	Error2700 = 2700
)

// Mapper used to keep error details for error logs
var Mapper = map[int]logging.ErrorDetails{
	Error1100: {
		ErrorCode: Error1100,
		Message:   "Failed to listen on port: %v, error: %v",
		Severity:  BLOCKER,
	},
	Error1101: {
		ErrorCode: Error1101,
		Message:   "Failed to start XDS GRPS server %s",
		Severity:  BLOCKER,
	},
	Error1102: {
		ErrorCode: Error1102,
		Message:   "Error reading the log configs. %v",
		Severity:  CRITICAL,
	},
	Error1103: {
		ErrorCode: Error1103,
		Message:   "Error while initializing authorization component, when intializing adapter REST API",
		Severity:  BLOCKER,
	},
	Error1104: {
		ErrorCode: Error1104,
		Message:   "Readiness probe is not set as local api artifacts processing has failed.",
		Severity:  CRITICAL,
	},
	Error1200: {
		ErrorCode: Error1200,
		Message:   "The provided port value for the REST Api Server :%v is not an integer. %v",
		Severity:  BLOCKER,
	},
	Error1400: {
		ErrorCode: Error1400,
		Message:   "Stream request for type %s on stream id: %d Error: %s",
		Severity:  CRITICAL,
	},
	Error1401: {
		ErrorCode: Error1401,
		Message:   "Stream request for type %s on stream id: %d, from node: %s, Error: %s",
		Severity:  CRITICAL,
	},
	Error1402: {
		ErrorCode: Error1402,
		Message:   "Consul syntax parse error %v",
		Severity:  CRITICAL,
	},
	Error1403: {
		ErrorCode: Error1403,
		Message:   "Internal Error while marshalling the upstream TLS Context. %v",
		Severity:  CRITICAL,
	},
	Error1410: {
		ErrorCode: Error1410,
		Message:   "Error undeploying API %v of Organization %v from environments %v",
		Severity:  MAJOR,
	},
	Error1411: {
		ErrorCode: Error1411,
		Message:   "Error extracting vhost from API identifier: %v for Organization %v. Ignore deploying the API",
		Severity:  MAJOR,
	},
	Error1413: {
		ErrorCode: Error1413,
		Message:   "Error creating new snapshot : %v",
		Severity:  MAJOR,
	},
	Error1414: {
		ErrorCode: Error1414,
		Message:   "Error while setting the snapshot : %v",
		Severity:  MAJOR,
	},
	Error1700: {
		ErrorCode: Error1700,
		Message:   "Error while connecting to the APK Management Server. %v",
		Severity:  BLOCKER,
	},
	Error1701: {
		ErrorCode: Error1701,
		Message:   "Error while starting APK Management application stream. %v",
		Severity:  BLOCKER,
	},
	Error1702: {
		ErrorCode: Error1702,
		Message:   "EOF is received from the APK Management Server application stream. %v",
		Severity:  CRITICAL,
	},
	Error1703: {
		ErrorCode: Error1703,
		Message:   "Failed to receive the discovery response from the APK Management Server application stream. %v",
		Severity:  CRITICAL,
	},
	Error1704: {
		ErrorCode: Error1704,
		Message:   "The APK Management Server application stream connection stopped: %v",
		Severity:  MINOR,
	},
	Error1705: {
		ErrorCode: Error1705,
		Message:   "Error while starting the APK Management Server: %v",
		Severity:  BLOCKER,
	},
	Error1706: {
		ErrorCode: Error1706,
		Message:   "Error while unmarshalling APK Management Server Application discovery response: %v",
		Severity:  MINOR,
	},
	Error1707: {
		ErrorCode: Error1707,
		Message:   "Error creating application: %v",
		Severity:  CRITICAL,
	},
	Error1709: {
		ErrorCode: Error1709,
		Message:   "Error updating application: %v",
		Severity:  CRITICAL,
	},
	Error1710: {
		ErrorCode: Error1710,
		Message:   "Error deleting application: %v",
		Severity:  CRITICAL,
	},
	Error1711: {
		ErrorCode: Error1711,
		Message:   "Error retrieving application: %v",
		Severity:  CRITICAL,
	},
	Error1800: {
		ErrorCode: Error1800,
		Message:   "error while parsing the interceptor template: %v",
		Severity:  CRITICAL,
	},
	Error1801: {
		ErrorCode: Error1801,
		Message:   "executing request interceptor template: %v",
		Severity:  CRITICAL,
	},
	Error2200: {
		ErrorCode: Error2200,
		Message:   "Error marsheling access log configs. %v",
		Severity:  CRITICAL,
	},
	Error2201: {
		ErrorCode: Error2201,
		Message:   "Error marshalling gRPC access log configs. %v",
		Severity:  CRITICAL,
	},
	Error2204: {
		ErrorCode: Error2204,
		Message:   "Operation policy validation failed for API %q in org %q:, policy %q: %v",
		Severity:  MINOR,
	},
	Error2205: {
		ErrorCode: Error2205,
		Message:   "Error parsing the operation policy definition %q into go template of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	Error2206: {
		ErrorCode: Error2206,
		Message:   "Error parsing operation policy definition %q of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	Error2207: {
		ErrorCode: Error2207,
		Message:   "Error parsing formalized operation policy definition %q into yaml of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	Error2208: {
		ErrorCode: Error2208,
		Message:   "API policy validation failed, policy: %q of the API %q in org %q: %v",
		Severity:  MINOR,
	},
	Error2209: {
		ErrorCode: Error2209,
		Message:   "Error while JSON unmarshalling to find the API definition version. %s",
		Severity:  MINOR,
	},
	Error2210: {
		ErrorCode: Error2210,
		Message:   "AsyncAPI version %s is not supported.",
		Severity:  MINOR,
	},
	Error2211: {
		ErrorCode: Error2211,
		Message:   "API definition version is not defined.",
		Severity:  MINOR,
	},
	Error2212: {
		ErrorCode: Error2212,
		Message:   "Error adding request policy %s to operation %s of resource %s. %v",
		Severity:  MINOR,
	},
	Error2231: {
		ErrorCode: Error2231,
		Message:   "Error while creating routes for API %s %s for path: %s Error: %s",
		Severity:  MAJOR,
	},
	Error2234: {
		ErrorCode: Error2234,
		Message:   "Error occurred while creating the compression filter: %v",
		Severity:  MINOR,
	},
	Error2235: {
		ErrorCode: Error2235,
		Message:   "Error while parsing the gzip configuration value for the memory level: %v",
		Severity:  MINOR,
	},
	Error2236: {
		ErrorCode: Error2236,
		Message:   "Error while parsing the gzip configuration value for the window bits: %v",
		Severity:  MINOR,
	},
	Error2237: {
		ErrorCode: Error2237,
		Message:   "Error while parsing the gzip configuration value for the compression level: %v",
		Severity:  MINOR,
	},
	Error2238: {
		ErrorCode: Error2238,
		Message:   "Error while parsing the gzip configuration value for the chunk size: %v",
		Severity:  MINOR,
	},
	Error2239: {
		ErrorCode: Error2239,
		Message:   "Error while adding resource level endpoints for %s:%v-%v. %v",
		Severity:  MAJOR,
	},
	Error2700: {
		ErrorCode: Error2700,
		Message:   "Error while processing the private-public key pair : %v",
		Severity:  BLOCKER,
	},
}
