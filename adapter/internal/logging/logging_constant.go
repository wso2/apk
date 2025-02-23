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

import (
	logging "github.com/wso2/apk/adapter/pkg/logging"
)

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
// - LoggerAPK
const (
	Error1100 = 1100
	Error1101 = 1101
	Error1102 = 1102
	Error1103 = 1103
	Error1104 = 1104
	Error1105 = 1105
)

// Error Log Internal discovery(1400-1499) Config Constants
// - LoggerXds
const (
	Error1400 = 1400
	Error1401 = 1401
	Error1403 = 1403
	Error1410 = 1410
	Error1411 = 1411
	Error1413 = 1413
	Error1414 = 1414
	Error1415 = 1415
)

// Error Log Internal XDS(1700-1799) Config Constants
// - LoggerXds
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
	Error1712 = 1712
	Error1713 = 1713
	Error1714 = 1714
	Error1715 = 1715
	Error1716 = 1716
	Error1717 = 1717
	Error1718 = 1718
	Error1719 = 1719
	Error1720 = 1720
	Error1721 = 1721
	Error1722 = 1722
	Error1723 = 1723
	Error1724 = 1724
)

// Error Log Internal intercepter(1800-1899) Config Constants
// - LoggerInterceptor
const (
	Error1800 = 1800
	Error1801 = 1801
)

// Error Log Internal OASParser(2200-2299) Config Constants
// - LoggerOasparser
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
	Error2240 = 2240
	Error2241 = 2241
	Error2242 = 2242
	Error2243 = 2243
	Error2244 = 2244
	Error2245 = 2245
	Error2246 = 2246
	Error2247 = 2247
	Error2248 = 2248
	Error2249 = 2249
	Error2250 = 2250
	Error2301 = 2301
	Error2302 = 2302
)

// Error Log RateLimiter callbacks(2300-2399) Config Constants
// - LoggerEnforcerXdsCallbacks
const (
	Error2300 = 2300
)

// Error Log Internal GRPC(2700-2799) Config Constants
const (
	Error2700 = 2700
)

// Mapper used to keep error details for error logs
var Mapper = map[int]logging.ErrorDetails{
	Error1100: {
		ErrorCode: Error1100,
		Message:   "Failed to listen on port.",
	},
	Error1101: {
		ErrorCode: Error1101,
		Message:   "Failed to start XDS GRPS server.",
	},
	Error1102: {
		ErrorCode: Error1102,
		Message:   "Error reading the log configs.",
	},
	Error1103: {
		ErrorCode: Error1103,
		Message:   "Error while initializing authorization component, when intializing adapter REST API",
	},
	Error1104: {
		ErrorCode: Error1104,
		Message:   "Readiness probe is not set as local api artifacts processing has failed.",
	},
	Error1105: {
		ErrorCode: Error1105,
		Message:   "Error serving Rate Limiter xDS gRPC server.",
	},
	Error1400: {
		ErrorCode: Error1400,
		Message:   "Error in Stream request type.",
	},
	Error1401: {
		ErrorCode: Error1401,
		Message:   "Error in Stream request type.",
	},
	Error1403: {
		ErrorCode: Error1403,
		Message:   "Internal Error while marshalling the upstream TLS Context.",
	},
	Error1410: {
		ErrorCode: Error1410,
		Message:   "Error undeploying API.",
	},
	Error1411: {
		ErrorCode: Error1411,
		Message:   "Error extracting vhost from API identifier. Ignore deploying the API",
	},
	Error1413: {
		ErrorCode: Error1413,
		Message:   "Error creating new snapshot.",
	},
	Error1414: {
		ErrorCode: Error1414,
		Message:   "Error while setting the snapshot.",
	},
	Error1700: {
		ErrorCode: Error1700,
		Message:   "Error while connecting to the APK Management Server.",
	},
	Error1701: {
		ErrorCode: Error1701,
		Message:   "Error while starting APK Management application stream.",
	},
	Error1702: {
		ErrorCode: Error1702,
		Message:   "EOF is received from the APK Management Server application stream.",
	},
	Error1703: {
		ErrorCode: Error1703,
		Message:   "Failed to receive the discovery response from the APK Management Server application stream.",
	},
	Error1704: {
		ErrorCode: Error1704,
		Message:   "The APK Management Server application stream connection stopped.",
	},
	Error1705: {
		ErrorCode: Error1705,
		Message:   "Error while starting the APK Management Server.",
	},
	Error1706: {
		ErrorCode: Error1706,
		Message:   "Error while unmarshalling APK Management Server Application discovery response.",
	},
	Error1707: {
		ErrorCode: Error1707,
		Message:   "Error creating application.",
	},
	Error1709: {
		ErrorCode: Error1709,
		Message:   "Error updating application.",
	},
	Error1710: {
		ErrorCode: Error1710,
		Message:   "Error deleting application.",
	},
	Error1711: {
		ErrorCode: Error1711,
		Message:   "Error retrieving application.",
	},
	Error1712: {
		ErrorCode: Error1712,
		Message:   "Unknown rate limit unit. Defaulting to UNKNOWN",
	},
	Error1713: {
		ErrorCode: Error1713,
		Message:   "Error extracting vhost from apiIdentifier. Continue cleaning other maps.",
	},
	Error1714: {
		ErrorCode: Error1714,
		Message:   "Error while creating the rate limit snapshot.",
	},
	Error1715: {
		ErrorCode: Error1715,
		Message:   "Inconsistent rate limiter snapshot.",
	},
	Error1716: {
		ErrorCode: Error1716,
		Message:   "Error while updating the rate limit snapshot.",
	},
	Error1717: {
		ErrorCode: Error1717,
		Message:   "EOF is received from the APK Management Server subscription stream.",
	},
	Error1718: {
		ErrorCode: Error1718,
		Message:   "Failed to receive the discovery response from the APK Management Server subscription stream.",
	},
	Error1719: {
		ErrorCode: Error1719,
		Message:   "The APK Management Server subscription stream connection stopped.",
	},
	Error1720: {
		ErrorCode: Error1720,
		Message:   "Error while unmarshalling APK Management Server Subscription discovery response.",
	},
	Error1721: {
		ErrorCode: Error1721,
		Message:   "Error creating subscription.",
	},
	Error1722: {
		ErrorCode: Error1722,
		Message:   "Error updating subscription.",
	},
	Error1723: {
		ErrorCode: Error1723,
		Message:   "Error deleting subscription.",
	},
	Error1724: {
		ErrorCode: Error1724,
		Message:   "Error retrieving subscription.",
	},
	Error1800: {
		ErrorCode: Error1800,
		Message:   "Error while parsing the interceptor template.",
	},
	Error1801: {
		ErrorCode: Error1801,
		Message:   "Executing request interceptor template.",
	},
	Error2200: {
		ErrorCode: Error2200,
		Message:   "Error marsheling access log configs.",
	},
	Error2201: {
		ErrorCode: Error2201,
		Message:   "Error marshalling gRPC access log configs.",
	},
	Error2204: {
		ErrorCode: Error2204,
		Message:   "Operation policy validation failed for API.",
	},
	Error2205: {
		ErrorCode: Error2205,
		Message:   "Error parsing the operation policy definition into go template of the API.",
	},
	Error2206: {
		ErrorCode: Error2206,
		Message:   "Error parsing operation policy definition of the API.",
	},
	Error2207: {
		ErrorCode: Error2207,
		Message:   "Error parsing formalized operation policy definition into yaml of the API.",
	},
	Error2208: {
		ErrorCode: Error2208,
		Message:   "API policy validation failed.",
	},
	Error2209: {
		ErrorCode: Error2209,
		Message:   "Error while JSON unmarshalling to find the API definition version.",
	},
	Error2210: {
		ErrorCode: Error2210,
		Message:   "AsyncAPI version is not supported.",
	},
	Error2211: {
		ErrorCode: Error2211,
		Message:   "API definition version is not defined.",
	},
	Error2212: {
		ErrorCode: Error2212,
		Message:   "Error adding request policy to operation.",
	},
	Error2231: {
		ErrorCode: Error2231,
		Message:   "Error while creating routes for API.",
	},
	Error2234: {
		ErrorCode: Error2234,
		Message:   "Error occurred while creating the compression filter.",
	},
	Error2235: {
		ErrorCode: Error2235,
		Message:   "Error while parsing the gzip configuration value for the memory level.",
	},
	Error2236: {
		ErrorCode: Error2236,
		Message:   "Error while parsing the gzip configuration value for the window bits.",
	},
	Error2237: {
		ErrorCode: Error2237,
		Message:   "Error while parsing the gzip configuration value for the compression level.",
	},
	Error2238: {
		ErrorCode: Error2238,
		Message:   "Error while parsing the gzip configuration value for the chunk size.",
	},
	Error2239: {
		ErrorCode: Error2239,
		Message:   "Error while adding resource level endpoints.",
	},
	Error2240: {
		ErrorCode: Error2240,
		Message:   "Invalid XRatelimitHeaders type, continue with default type.",
	},
	Error2241: {
		ErrorCode: Error2241,
		Message:   "Error occurred while parsing ratelimit filter config.",
	},
	Error2242: {
		ErrorCode: Error2242,
		Message:   "Error while adding api level request intercepter external cluster.",
	},
	Error2243: {
		ErrorCode: Error2243,
		Message:   "Error while adding api level response intercepter external cluster.",
	},
	Error2244: {
		ErrorCode: Error2244,
		Message:   "Error while adding resource level request intercept external cluster.",
	},
	Error2245: {
		ErrorCode: Error2245,
		Message:   "Error while adding operational level request intercept external cluster.",
	},
	Error2246: {
		ErrorCode: Error2246,
		Message:   "Error while adding resource level response intercept external cluster.",
	},
	Error2247: {
		ErrorCode: Error2247,
		Message:   "Error while adding operational level response intercept external cluster.",
	},
	Error2248: {
		ErrorCode: Error2248,
		Message:   "Failed to initialize ratelimit cluster. Hence terminating the adapter.",
	},
	Error2249: {
		ErrorCode: Error2249,
		Message:   "Failed to initialize tracer's cluster. Router tracing will be disabled.",
	},
	Error2250: {
		ErrorCode: Error2250,
		Message:   "Failed to parse JWTAuthentication",
	},
	Error2700: {
		ErrorCode: Error2700,
		Message:   "Error while processing the private-public key pair.",
	},
	Error2300: {
		ErrorCode: Error2300,
		Message:   "Error in Stream request.",
	},
	Error2301: {
		ErrorCode: Error2301,
		Message:   "Error while generating JWTProviders",
	},
	Error2302: {
		ErrorCode: Error2302,
		Message:   "Error while generating APIKey JWTProvider",
	},
}
