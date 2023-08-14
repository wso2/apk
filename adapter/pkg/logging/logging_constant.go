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

// Log (Error) severity level constants
const (
	BLOCKER  = "Blocker"
	CRITICAL = "Critical"
	MAJOR    = "Major"
	MINOR    = "Minor"
	TRIVIAL  = "Trivial"
	DEFAULT  = "Default"
)

// Error Log Internal Configuration(1000-1099) Config Constants
// - loggerConfig
const (
	Error1000 = 1000
	Error1001 = 1001
	Error1002 = 1002
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

// Error Log RateLimiter callbacks(2300-2399) Config Constants
// - LoggerEnforcerXdsCallbacks
const (
	Error2300 = 2300
)

// Error Log Pkg operator(2600-2699) Config Constants
// - LoggerAPKOperator
const (
	Error2600 = 2600
	Error2601 = 2601
	Error2602 = 2602
	Error2603 = 2603
	Error2604 = 2604
	Error2605 = 2605
	Error2606 = 2606
	Error2607 = 2607
	Error2608 = 2608
	Error2609 = 2609
	Error2610 = 2610
	Error2611 = 2611
	Error2612 = 2612
	Error2613 = 2613
	Error2614 = 2614
	Error2615 = 2615
	Error2616 = 2616
	Error2617 = 2617
	Error2618 = 2618
	Error2619 = 2619
	Error2620 = 2620
	Error2621 = 2621
	Error2622 = 2622
	Error2623 = 2623
	Error2625 = 2625
	Error2626 = 2626
	Error2627 = 2627
	Error2628 = 2628
	Error2629 = 2629
	Error2630 = 2630
	Error2631 = 2631
	Error2632 = 2632
	Error2633 = 2633
	Error2634 = 2634
	Error2635 = 2635
	Error2636 = 2636
	Error2637 = 2637
	Error2638 = 2638
	Error2639 = 2639
	Error2640 = 2640
	Error2641 = 2641
	Error2642 = 2642
	Error2643 = 2643
	Error2644 = 2644
	Error2645 = 2645
	Error2646 = 2646
	Error2647 = 2647
	Error2648 = 2648
	Error2649 = 2649
	Error2650 = 2650
	Error2651 = 2651
	Error2652 = 2652
	Error2653 = 2653
	Error2654 = 2654
	Error2655 = 2655
	Error2656 = 2656
	Error2657 = 2657
	Error2658 = 2658
	Error2659 = 2659
	Error2660 = 2660
	Error2661 = 2661
	Error2662 = 2662
	Error2663 = 2663
)

// Error Log Pkg auth(3001-3099) Config Constants
const (
	Error3001 = 3001
	Error3002 = 3002
)

// Error codes gateway controller (3100-3199)
const (
	Error3100 = 3100
	Error3101 = 3101
	Error3102 = 3102
	Error3103 = 3103
	Error3104 = 3104
	Error3105 = 3105
	Error3106 = 3106
	Error3107 = 3107
	Error3108 = 3108
	Error3109 = 3109
	Error3110 = 3110
	Error3111 = 3111
	Error3112 = 3112
	Error3113 = 3113
	Error3114 = 3114
	Error3115 = 3115
	Error3116 = 3116
	Error3117 = 3117
	Error3118 = 3118
	Error3119 = 3119
	Error3120 = 3120
	Error3121 = 3121
	Error3122 = 3122
	Error3123 = 3123
	Error3124 = 3124
	Error3125 = 3125
	Error3126 = 3126
)

// Mapper used to keep error details for error logs
var Mapper = map[int]ErrorDetails{
	Error1000: {
		ErrorCode: Error1000,
		Message:   "Configuration file not found.",
	},
	Error1001: {
		ErrorCode: Error1001,
		Message:   "Error reading configurations.",
	},
	Error1002: {
		ErrorCode: Error1002,
		Message:   "Error parsing the configurations.",
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
	Error2300: {
		ErrorCode: Error2300,
		Message:   "Error in Stream request.",
	},
	Error2600: {
		ErrorCode: Error2600,
		Message:   "unable to start manager.",
	},
	Error2601: {
		ErrorCode: Error2601,
		Message:   "Unable to create webhook API.",
	},
	Error2602: {
		ErrorCode: Error2602,
		Message:   "unable to set up health check.",
	},
	Error2603: {
		ErrorCode: Error2603,
		Message:   "unable to set up ready check.",
	},
	Error2604: {
		ErrorCode: Error2604,
		Message:   "Problem running manager.",
	},
	Error2605: {
		ErrorCode: Error2605,
		Message:   "Unable to list APIs.",
	},
	Error2606: {
		ErrorCode: Error2606,
		Message:   "Error creating Application controller.",
	},
	Error2607: {
		ErrorCode: Error2607,
		Message:   "Error watching Application resources.",
	},
	Error2608: {
		ErrorCode: Error2608,
		Message:   "Error creating Subscription controller.",
	},
	Error2609: {
		ErrorCode: Error2609,
		Message:   "Error watching Subscription resources.",
	},
	Error2610: {
		ErrorCode: Error2610,
		Message:   "Error creating API controller .",
	},
	Error2611: {
		ErrorCode: Error2611,
		Message:   "Error watching API resources.",
	},
	Error2612: {
		ErrorCode: Error2612,
		Message:   "Error adding indexes.",
	},
	Error2613: {
		ErrorCode: Error2613,
		Message:   "Error watching HTTPRoute resources.",
	},
	Error2614: {
		ErrorCode: Error2614,
		Message:   "Error watching Service resources.",
	},
	Error2615: {
		ErrorCode: Error2615,
		Message:   "Error watching Backend resources.",
	},
	Error2616: {
		ErrorCode: Error2616,
		Message:   "Error watching Authentication resources.",
	},
	Error2617: {
		ErrorCode: Error2617,
		Message:   "Error watching APIPolicy resources.",
	},
	Error2618: {
		ErrorCode: Error2618,
		Message:   "Error watching scope resources.",
	},
	Error2619: {
		ErrorCode: Error2619,
		Message:   "Error applying startup APIs.",
	},
	Error2620: {
		ErrorCode: Error2620,
		Message:   "Error resolving Gateway State.",
	},
	Error2622: {
		ErrorCode: Error2622,
		Message:   "Unexpected object type, bypassing reconciliation.",
	},
	Error2623: {
		ErrorCode: Error2623,
		Message:   "Unable to find associated APIs.",
	},
	Error2625: {
		ErrorCode: Error2625,
		Message:   "Unable to find associated HTTPRoutes.",
	},
	Error2626: {
		ErrorCode: Error2626,
		Message:   "Unsupported object type.",
	},
	Error2628: {
		ErrorCode: Error2628,
		Message:   "API Event is nil.",
	},
	Error2629: {
		ErrorCode: Error2629,
		Message:   "API deployment failed for event.",
	},
	Error2630: {
		ErrorCode: Error2630,
		Message:   "Error undeploying prod httpRoute of API. Hence not checking on deleting the sand httpRoute of the API",
	},
	Error2631: {
		ErrorCode: Error2631,
		Message:   "Error setting HttpRoute CR info to adapterInternalAPI.",
	},
	Error2632: {
		ErrorCode: Error2632,
		Message:   "Error validating adapterInternalAPI intermediate representation.",
	},
	Error2633: {
		ErrorCode: Error2633,
		Message:   "Error updating the API.",
	},
	Error2634: {
		ErrorCode: Error2634,
		Message:   "Error creating connection.",
	},
	Error2635: {
		ErrorCode: Error2635,
		Message:   "Error sending API to APK management server.",
	},
	Error2636: {
		ErrorCode: Error2636,
		Message:   "Error while generating the soap fault message.",
	},
	Error2637: {
		ErrorCode: Error2637,
		Message:   "Unable to create webhook for Ratelimit.",
	},
	Error2638: {
		ErrorCode: Error2638,
		Message:   "Unable to create webhook for APIPolicy.",
	},
	Error2639: {
		ErrorCode: Error2639,
		Message:   "Error watching Ratelimit resources.",
	},
	Error2640: {
		ErrorCode: Error2640,
		Message:   "Error watching InterceptorService resources.",
	},
	Error2653: {
		ErrorCode: Error2653,
		Message:   "Gateway Label is invalid.",
	},
	Error2654: {
		ErrorCode: Error2654,
		Message:   "Error resolving certificate for Backend.",
	},
	Error2661: {
		ErrorCode: Error2661,
		Message:   "Error watching BackendJWT resources.",
	},
	Error2662: {
		ErrorCode: Error2662,
		Message:   "Error while getting BackendJWT.",
	},
	Error2663: {
		ErrorCode: Error2663,
		Message:   "Error creating Ratelimit controller.",
	},
	Error2621: {
		ErrorCode: Error2621,
		Message:   "Unable to find associated Backends for Secret.",
	},
	Error2627: {
		ErrorCode: Error2627,
		Message:   "Failed to decode certificate PEM",
	},
	Error2641: {
		ErrorCode: Error2641,
		Message:   "Error while parsing certificate.",
	},
	Error2642: {
		ErrorCode: Error2642,
		Message:   "Error while reading certificate from secretRef.",
	},
	Error2643: {
		ErrorCode: Error2643,
		Message:   "Error while reading certificate from configMapRef.",
	},
	Error2644: {
		ErrorCode: Error2644,
		Message:   "Error watching ConfigMap resources.",
	},
	Error2645: {
		ErrorCode: Error2645,
		Message:   "Error watching Secret resources.",
	},
	Error2646: {
		ErrorCode: Error2646,
		Message:   "Error while getting Backend.",
	},
	Error2647: {
		ErrorCode: Error2647,
		Message:   "Unable to find associated Backends for ConfigMap.",
	},
	Error2648: {
		ErrorCode: Error2648,
		Message:   "Error while reading key from secretRef.",
	},
	Error2649: {
		ErrorCode: Error2649,
		Message:   "Unable to find associated APIPolicies.",
	},
	Error2650: {
		ErrorCode: Error2650,
		Message:   "Error while getting custom rate limit policies.",
	},
	Error2651: {
		ErrorCode: Error2651,
		Message:   "Error while getting interceptor service.",
	},
	Error2652: {
		ErrorCode: Error2652,
		Message:   "Unable to create webhook for InterceptorService.",
	},
	Error2655: {
		ErrorCode: Error2655,
		Message:   "Unable to create webhook for Backend.",
	},
	Error2656: {
		ErrorCode: Error2656,
		Message:   "Error watching JWTIssuer resources.",
	},
	Error2657: {
		ErrorCode: Error2657,
		Message:   "Error creating JWTIssuer controller.",
	},
	Error2658: {
		ErrorCode: Error2658,
		Message:   "Error adding indexes.",
	},
	Error2659: {
		ErrorCode: Error2659,
		Message:   "Error resolving certificate for JWKS.",
	},
	Error2660: {
		ErrorCode: Error2660,
		Message:   "Unable to find associated JWTIssuers.",
	},
	Error3001: {
		ErrorCode: Error3001,
		Message:   "Error reading ssh key file.",
	},
	Error3002: {
		ErrorCode: Error3002,
		Message:   "Error creating ssh public key.",
	},
	Error3100: {
		ErrorCode: Error3100,
		Message:   "Error watching Gateway resources.",
	},
	Error3101: {
		ErrorCode: Error3101,
		Message:   "Error watching APIPolicy resources.",
	},
	Error3102: {
		ErrorCode: Error3102,
		Message:   "Error watching Backend resources.",
	},
	Error3103: {
		ErrorCode: Error3103,
		Message:   "Error watching ConfigMap resources.",
	},
	Error3104: {
		ErrorCode: Error3104,
		Message:   "Error watching Secret resources.",
	},
	Error3105: {
		ErrorCode: Error3105,
		Message:   "Error resolving listener certificates.",
	},
	Error3106: {
		ErrorCode: Error3106,
		Message:   "Unable to find associated Backends for Secret.",
	},
	Error3107: {
		ErrorCode: Error3107,
		Message:   "Unexpected object type, bypassing reconciliation.",
	},
	Error3108: {
		ErrorCode: Error3108,
		Message:   "Unable to find associated Backends for ConfigMap.",
	},
	Error3109: {
		ErrorCode: Error3109,
		Message:   "Error while updating Gateway status.",
	},
	Error3110: {
		ErrorCode: Error3110,
		Message:   "Error watching InterceptorService resources.",
	},
	Error3111: {
		ErrorCode: Error3111,
		Message:   "Error creating JWTIssuer controller.",
	},
	Error3112: {
		ErrorCode: Error3112,
		Message:   "Error adding indexes.",
	},
	Error3113: {
		ErrorCode: Error3113,
		Message:   "Error resolving certificate for JWKS.",
	},
	Error3114: {
		ErrorCode: Error3114,
		Message:   "Error creating JWT Issuer controller.",
	},
	Error3115: {
		ErrorCode: Error3115,
		Message:   "Route Timeout cannot be greater than the Max value defined.",
	},
	Error3116: {
		ErrorCode: Error3116,
		Message:   "Invalid Status Codes for Retry.",
	},
	Error3117: {
		ErrorCode: Error3117,
		Message:   "Retry Count Should be greater than 0.",
	},
	Error3118: {
		ErrorCode: Error3118,
		Message:   "Unable to find associated interceptorServices.",
	},
	Error3119: {
		ErrorCode: Error3119,
		Message:   "Error creating API controller .",
	},
	Error3120: {
		ErrorCode: Error3120,
		Message:   "Error adding indexes.",
	},
	Error3121: {
		ErrorCode: Error3121,
		Message:   "Error watching Ratelimit resources.",
	},
	Error3122: {
		ErrorCode: Error3122,
		Message:   "Error resolving Gateway State.",
	},
	Error3123: {
		ErrorCode: Error3123,
		Message:   "Unable to find associated secret.",
	},
	Error3124: {
		ErrorCode: Error3124,
		Message:   "Error while getting custom rate limit policies.",
	},
	Error3125: {
		ErrorCode: Error3125,
		Message:   "Unable to find associated APIPolicies.",
	},
	Error3126: {
		ErrorCode: Error3126,
		Message:   "Error watching BackendJWT resources.",
	},
}
