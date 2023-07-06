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

// Error Log Pkg operator(2600-2699) Config Constants
// - LoggerAPKOperator
const (
	error2600 = 2600
	error2601 = 2601
	error2602 = 2602
	error2603 = 2603
	error2604 = 2604
	error2605 = 2605
	error2606 = 2606
	error2607 = 2607
	error2608 = 2608
	error2609 = 2609
	error2610 = 2610
	error2611 = 2611
	error2612 = 2612
	error2613 = 2613
	error2614 = 2614
	error2615 = 2615
	error2616 = 2616
	error2617 = 2617
	error2618 = 2618
	error2619 = 2619
	error2621 = 2621
	error2622 = 2622
	error2623 = 2623
	error2624 = 2624
	error2625 = 2625
	error2626 = 2626
	error2627 = 2627
	error2628 = 2628
	error2629 = 2629
	error2630 = 2630
	error2631 = 2631
	error2632 = 2632
	error2633 = 2633
	error2634 = 2634
	error2635 = 2635
	error2636 = 2636
	error2637 = 2637
	error2638 = 2638
	error2639 = 2639
	error2640 = 2640
	error2641 = 2641
	error2642 = 2642
	error2643 = 2643
	error2644 = 2644
	error2645 = 2645
	error2646 = 2646
	error2647 = 2647
	error2648 = 2648
	error2649 = 2649
	error2650 = 2650
	error2651 = 2651
	error2652 = 2652
	error2653 = 2653
	error2654 = 2654
	error2655 = 2655
	error2656 = 2656
)

// Error Log Pkg auth(3001-3099) Config Constants
const (
	error3001 = 3001
	error3002 = 3002
)

// Error codes gateway controller (3100-3199)
const (
	error3100 = 3100
	error3101 = 3101
	error3102 = 3102
	error3103 = 3103
	error3104 = 3104
	error3105 = 3105
	error3106 = 3106
	error3107 = 3107
	error3108 = 3108
	error3109 = 3109
	error3110 = 3110
)
const (
	error3111 = 3111
	error3112 = 3112
	error3113 = 3113
	error3114 = 3114
	error3115 = 3115
	error3116 = 3116
	error3117 = 3117
)

// Mapper used to keep error details for error logs
var Mapper = map[int]ErrorDetails{
	error2600: {
		ErrorCode: error2600,
		Message:   "unable to start manager: %v",
		Severity:  BLOCKER,
	},
	error2601: {
		ErrorCode: error2601,
		Message:   "Unable to create webhook API: %v",
		Severity:  BLOCKER,
	},
	error2602: {
		ErrorCode: error2602,
		Message:   "unable to set up health check: %v",
		Severity:  BLOCKER,
	},
	error2603: {
		ErrorCode: error2603,
		Message:   "unable to set up ready check: %v",
		Severity:  BLOCKER,
	},
	error2604: {
		ErrorCode: error2604,
		Message:   "uproblem running manager: %v",
		Severity:  BLOCKER,
	},
	error2605: {
		ErrorCode: error2605,
		Message:   "unable to list APIs for API context validation: %v",
		Severity:  CRITICAL,
	},
	error2606: {
		ErrorCode: error2606,
		Message:   "Error creating Application controller: %v",
		Severity:  BLOCKER,
	},
	error2607: {
		ErrorCode: error2607,
		Message:   "Error watching Application resources: %v",
		Severity:  BLOCKER,
	},
	error2608: {
		ErrorCode: error2608,
		Message:   "Error creating Subscription controller: %v",
		Severity:  BLOCKER,
	},
	error2609: {
		ErrorCode: error2609,
		Message:   "Error watching Subscription resources: %v",
		Severity:  BLOCKER,
	},
	error2610: {
		ErrorCode: error2610,
		Message:   "Error creating API controller : %v",
		Severity:  BLOCKER,
	},
	error2611: {
		ErrorCode: error2611,
		Message:   "Error watching API resources: %v",
		Severity:  BLOCKER,
	},
	error2612: {
		ErrorCode: error2612,
		Message:   "Error adding indexes: %v",
		Severity:  BLOCKER,
	},
	error2613: {
		ErrorCode: error2613,
		Message:   "Error watching HTTPRoute resources: %v",
		Severity:  BLOCKER,
	},
	error2614: {
		ErrorCode: error2614,
		Message:   "Error watching Service resources: %v",
		Severity:  BLOCKER,
	},
	error2615: {
		ErrorCode: error2615,
		Message:   "Error watching Backend resources: %v",
		Severity:  BLOCKER,
	},
	error2616: {
		ErrorCode: error2616,
		Message:   "Error watching Authentication resources: %v",
		Severity:  BLOCKER,
	},
	error2617: {
		ErrorCode: error2617,
		Message:   "Error watching APIPolicy resources: %v",
		Severity:  BLOCKER,
	},
	error2618: {
		ErrorCode: error2618,
		Message:   "Error watching scope resources: %v",
		Severity:  BLOCKER,
	},
	error2619: {
		ErrorCode: error2619,
		Message:   "Error applying startup APIs: %v",
		Severity:  BLOCKER,
	},
	error2622: {
		ErrorCode: error2622,
		Message:   "Unexpected object type, bypassing reconciliation: %v",
		Severity:  TRIVIAL,
	},
	error2623: {
		ErrorCode: error2623,
		Message:   "Unable to find associated APIs: %s",
		Severity:  CRITICAL,
	},
	error2624: {
		ErrorCode: error2624,
		Message:   "Unexpected object type, bypassing reconciliation: %v",
		Severity:  TRIVIAL,
	},
	error2625: {
		ErrorCode: error2625,
		Message:   "Unable to find associated HTTPRoutes: %s",
		Severity:  CRITICAL,
	},
	error2626: {
		ErrorCode: error2626,
		Message:   "Unsupported object type %T",
		Severity:  BLOCKER,
	},
	error2628: {
		ErrorCode: error2628,
		Message:   "API Event is nil",
		Severity:  BLOCKER,
	},
	error2629: {
		ErrorCode: error2629,
		Message:   "API deployment failed for %s event : %v",
		Severity:  MAJOR,
	},
	error2630: {
		ErrorCode: error2630,
		Message: "Error undeploying prod httpRoute of API : %v in Organization %v from environments %v." +
			" Hence not checking on deleting the sand httpRoute of the API",
		Severity: MAJOR,
	},
	error2631: {
		ErrorCode: error2631,
		Message:   "Error setting HttpRoute CR info to adapterInternalAPI. %v",
		Severity:  MAJOR,
	},
	error2632: {
		ErrorCode: error2632,
		Message:   "Error validating adapterInternalAPI intermediate representation. %v",
		Severity:  MAJOR,
	},
	error2633: {
		ErrorCode: error2633,
		Message:   "Error updating the API : %s:%s in vhosts: %s. %v",
		Severity:  "MAJOR",
	},
	error2634: {
		ErrorCode: error2634,
		Message:   "Error creating connection for %v : %v",
		Severity:  MAJOR,
	},
	error2635: {
		ErrorCode: error2635,
		Message:   "Error sending API to APK management server : %v",
		Severity:  MAJOR,
	},
	error2636: {
		ErrorCode: error2636,
		Message:   "Error while generating the soap fault message. %s",
		Severity:  MINOR,
	},
	error2637: {
		ErrorCode: error2637,
		Message:   "Unable to create webhook for Ratelimit: %v",
		Severity:  BLOCKER,
	},
	error2638: {
		ErrorCode: error2638,
		Message:   "Unable to create webhook for APIPolicy: %v",
		Severity:  BLOCKER,
	},
	error2639: {
		ErrorCode: error2639,
		Message:   "Error watching Ratelimit resources: %v",
		Severity:  BLOCKER,
	},
	error2640: {
		ErrorCode: error2640,
		Message:   "Error watching InterceptorService resources: %v",
		Severity:  BLOCKER,
	},
	error2653: {
		ErrorCode: error2653,
		Message:   "Gateway Label is invalid: %s",
		Severity:  CRITICAL,
	},
	error2654: {
		ErrorCode: error2654,
		Message:   "Error resolving certificate for Backend %v",
		Severity:  BLOCKER,
	},
	error2621: {
		ErrorCode: error2621,
		Message:   "Unable to find associated Backends for Secret: %s",
		Severity:  CRITICAL,
	},
	error2627: {
		ErrorCode: error2627,
		Message:   "Failed to decode certificate PEM",
		Severity:  CRITICAL,
	},
	error2641: {
		ErrorCode: error2641,
		Message:   "Error while parsing certificate: %s",
		Severity:  CRITICAL,
	},
	error2642: {
		ErrorCode: error2642,
		Message:   "Error while reading certificate from secretRef: %s",
		Severity:  CRITICAL,
	},
	error2643: {
		ErrorCode: error2643,
		Message:   "Error while reading certificate from configMapRef: %s",
		Severity:  CRITICAL,
	},
	error2644: {
		ErrorCode: error2644,
		Message:   "Error watching ConfigMap resources: %v",
		Severity:  BLOCKER,
	},
	error2645: {
		ErrorCode: error2645,
		Message:   "Error watching Secret resources: %v",
		Severity:  BLOCKER,
	},
	error2646: {
		ErrorCode: error2646,
		Message:   "Error while getting backend: %v, error: %v",
		Severity:  CRITICAL,
	},
	error2647: {
		ErrorCode: error2647,
		Message:   "Unable to find associated Backends for ConfigMap: %s",
		Severity:  CRITICAL,
	},
	error2648: {
		ErrorCode: error2648,
		Message:   "Error while reading key from secretRef: %s",
		Severity:  CRITICAL,
	},
	error2649: {
		ErrorCode: error2649,
		Message:   "Unable to find associated APIPolicies: %s",
		Severity:  CRITICAL,
	},
	error2650: {
		ErrorCode: error2650,
		Message:   "Error while getting custom rate limit policies: %s",
		Severity:  MAJOR,
	},
	error2651: {
		ErrorCode: error2651,
		Message:   "Error while getting interceptor service %s, %s",
		Severity:  CRITICAL,
	},
	error2652: {
		ErrorCode: error2652,
		Message:   "Unable to create webhook for InterceptorService: %v",
		Severity:  BLOCKER,
	},
	error2655: {
		ErrorCode: error2655,
		Message:   "Unable to create webhook for Authentication: %v",
		Severity:  BLOCKER,
	},
	error2656: {
		ErrorCode: error2656,
		Message:   "Unable to create webhook for Backend: %v",
		Severity:  BLOCKER,
	},
	error3001: {
		ErrorCode: error3001,
		Message:   "Error reading ssh key file: %s",
		Severity:  CRITICAL,
	},
	error3002: {
		ErrorCode: error3002,
		Message:   "Error creating ssh public key: %s",
		Severity:  CRITICAL,
	},
	error3100: {
		ErrorCode: error3100,
		Message:   "Error watching Gateway resources: %v",
		Severity:  BLOCKER,
	},
	error3101: {
		ErrorCode: error3101,
		Message:   "Error watching APIPolicy resources: %v",
		Severity:  BLOCKER,
	},
	error3102: {
		ErrorCode: error3102,
		Message:   "Error watching Backend resources: %v",
		Severity:  BLOCKER,
	},
	error3103: {
		ErrorCode: error3103,
		Message:   "Error watching ConfigMap resources: %v",
		Severity:  BLOCKER,
	},
	error3104: {
		ErrorCode: error3104,
		Message:   "Error watching Secret resources: %v",
		Severity:  BLOCKER,
	},
	error3105: {
		ErrorCode: error3105,
		Message:   "Error resolving listener certificates: %v",
		Severity:  BLOCKER,
	},
	error3106: {
		ErrorCode: error3106,
		Message:   "Unable to find associated Backends for Secret: %s",
		Severity:  CRITICAL,
	},
	error3107: {
		ErrorCode: error3107,
		Message:   "Unexpected object type, bypassing reconciliation: %v",
		Severity:  TRIVIAL,
	},
	error3108: {
		ErrorCode: error3108,
		Message:   "Unable to find associated Backends for ConfigMap: %s",
		Severity:  CRITICAL,
	},
	error3109: {
		ErrorCode: error3109,
		Message:   "Error while updating Gateway status %v",
		Severity:  BLOCKER,
	},
	error3110: {
		ErrorCode: error3110,
		Message:   "Error watching InterceptorService resources: %v",
		Severity:  BLOCKER,
	},
	error3111: {
		ErrorCode: error3111,
		Message:   "Error creating JWTIssuer controller: %v",
		Severity:  MAJOR,
	},
	error3112: {
		ErrorCode: error3112,
		Message:   "Error adding indexes: %v",
		Severity:  BLOCKER,
	},
	error3113: {
		ErrorCode: error3113,
		Message:   "Error resolving certificate for JWKS %v",
		Severity:  MAJOR,
	},
	error3114: {
		ErrorCode: error3114,
		Message:   "Error creating JWT Issuer controller: %v",
		Severity:  BLOCKER,
	},
	error3115: {
		ErrorCode: error3115,
		Message:   "Route Timeout cannot be greater than the Max value defined : %s",
		Severity:  MAJOR,
	},
	error3116: {
		ErrorCode: error3116,
		Message:   "Invalid Status Codes for Retry: %s",
		Severity:  MAJOR,
	},
	error3117: {
		ErrorCode: error3117,
		Message:   "Retry Count Should be greater than 0: %s",
		Severity:  MAJOR,
	},
}
