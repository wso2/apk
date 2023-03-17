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
	Error2624 = 2624
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
)

// Error Log Pkg auth(3001-3099) Config Constants
const (
	Error3001 = 3001
	Error3002 = 3002
)

// Mapper used to keep error details for error logs
var Mapper = map[int]ErrorDetails{
	Error2600: {
		ErrorCode: Error2600,
		Message:   "unable to start manager: %v",
		Severity:  BLOCKER,
	},
	Error2601: {
		ErrorCode: Error2601,
		Message:   "Unable to create webhook API: %v",
		Severity:  BLOCKER,
	},
	Error2602: {
		ErrorCode: Error2602,
		Message:   "unable to set up health check: %v",
		Severity:  BLOCKER,
	},
	Error2603: {
		ErrorCode: Error2603,
		Message:   "unable to set up ready check: %v",
		Severity:  BLOCKER,
	},
	Error2604: {
		ErrorCode: Error2604,
		Message:   "uproblem running manager: %v",
		Severity:  BLOCKER,
	},
	Error2605: {
		ErrorCode: Error2605,
		Message:   "unable to list APIs for API context validation: %v",
		Severity:  CRITICAL,
	},
	Error2606: {
		ErrorCode: Error2606,
		Message:   "Error creating Application controller: %v",
		Severity:  BLOCKER,
	},
	Error2607: {
		ErrorCode: Error2607,
		Message:   "Error watching Application resources: %v",
		Severity:  BLOCKER,
	},
	Error2608: {
		ErrorCode: Error2608,
		Message:   "Error creating Subscription controller: %v",
		Severity:  BLOCKER,
	},
	Error2609: {
		ErrorCode: Error2609,
		Message:   "Error watching Subscription resources: %v",
		Severity:  BLOCKER,
	},
	Error2610: {
		ErrorCode: Error2610,
		Message:   "Error creating API controller : %v",
		Severity:  BLOCKER,
	},
	Error2611: {
		ErrorCode: Error2611,
		Message:   "Error watching API resources: %v",
		Severity:  BLOCKER,
	},
	Error2612: {
		ErrorCode: Error2612,
		Message:   "Error adding indexes: %v",
		Severity:  BLOCKER,
	},
	Error2613: {
		ErrorCode: Error2613,
		Message:   "Error watching HTTPRoute resources: %v",
		Severity:  BLOCKER,
	},
	Error2614: {
		ErrorCode: Error2614,
		Message:   "Error watching Service resources: %v",
		Severity:  BLOCKER,
	},
	Error2615: {
		ErrorCode: Error2615,
		Message:   "Error watching BackendPolicy resources: %v",
		Severity:  BLOCKER,
	},
	Error2616: {
		ErrorCode: Error2616,
		Message:   "Error watching Authentication resources: %v",
		Severity:  BLOCKER,
	},
	Error2617: {
		ErrorCode: Error2617,
		Message:   "Error watching APIPolicy resources: %v",
		Severity:  BLOCKER,
	},
	Error2618: {
		ErrorCode: Error2618,
		Message:   "Error watching scope resources: %v",
		Severity:  BLOCKER,
	},
	Error2619: {
		ErrorCode: Error2619,
		Message: "Api CR related to the reconcile request with key: %s returned error." +
			" Assuming API is already deleted, hence ignoring the error : %v",
		Severity: TRIVIAL,
	},
	Error2620: {
		ErrorCode: Error2620,
		Message:   "Error retrieving ref CRs for API in namespace : %s, %v",
		Severity:  TRIVIAL,
	},
	Error2621: {
		ErrorCode: Error2621,
		Message:   "Unable to find associated BackendPolicies for service: %s",
		Severity:  CRITICAL,
	},
	Error2622: {
		ErrorCode: Error2622,
		Message:   "Unexpected object type, bypassing reconciliation: %v",
		Severity:  TRIVIAL,
	},
	Error2623: {
		ErrorCode: Error2623,
		Message:   "Unable to find associated APIs: %s",
		Severity:  CRITICAL,
	},
	Error2624: {
		ErrorCode: Error2624,
		Message:   "Unexpected object type, bypassing reconciliation: %v",
		Severity:  TRIVIAL,
	},
	Error2625: {
		ErrorCode: Error2625,
		Message:   "Unable to find associated HTTPRoutes: %s",
		Severity:  CRITICAL,
	},
	Error2626: {
		ErrorCode: Error2626,
		Message:   "Unsupported object type %T",
		Severity:  BLOCKER,
	},
	Error2627: {
		ErrorCode: Error2627,
		Message:   "Unable to find associated Service for BackendPolicy: %s",
		Severity:  BLOCKER,
	},
	Error2628: {
		ErrorCode: Error2628,
		Message:   "API Event is nil",
		Severity:  BLOCKER,
	},
	Error2629: {
		ErrorCode: Error2629,
		Message:   "API deployment failed for %s event : %v",
		Severity:  MAJOR,
	},
	Error2630: {
		ErrorCode: Error2630,
		Message: "Error undeploying prod httpRoute of API : %v in Organization %v from environments %v." +
			" Hence not checking on deleting the sand httpRoute of the API",
		Severity: MAJOR,
	},
	Error2631: {
		ErrorCode: Error2631,
		Message:   "Error setting HttpRoute CR info to mgwSwagger. %v",
		Severity:  MAJOR,
	},
	Error2632: {
		ErrorCode: Error2632,
		Message:   "Error validating mgwSwagger intermediate representation. %v",
		Severity:  MAJOR,
	},
	Error2633: {
		ErrorCode: Error2633,
		Message:   "Error updating the API : %s:%s in vhosts: %s. %v",
		Severity:  "MAJOR",
	},
	Error2634: {
		ErrorCode: Error2634,
		Message:   "Error creating connection for %v : %v",
		Severity:  MAJOR,
	},
	Error2635: {
		ErrorCode: Error2635,
		Message:   "Error sending API to APK management server : %v",
		Severity:  MAJOR,
	},
	Error2636: {
		ErrorCode: Error2636,
		Message:   "Error while generating the soap fault message. %s",
		Severity:  MINOR,
	},
	Error2637: {
		ErrorCode: Error2637,
		Message:   "Unable to create webhook for Ratelimit: %v",
		Severity:  BLOCKER,
	},
	Error2638: {
		ErrorCode: Error2638,
		Message:   "Unable to create webhook for APIPolicy: %v",
		Severity:  BLOCKER,
	},
	Error3001: {
		ErrorCode: Error3001,
		Message:   "Error reading ssh key file: %s",
		Severity:  CRITICAL,
	},
	Error3002: {
		ErrorCode: Error3002,
		Message:   "Error creating ssh public key: %s",
		Severity:  CRITICAL,
	},
}
