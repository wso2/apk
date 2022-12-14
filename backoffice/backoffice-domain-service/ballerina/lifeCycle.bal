//
// Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//


public type StatesList record {
    States[] States;
};

public type States record {
    string State;
    Transitions[] Transitions?;
};

public type Transitions record {
    string event?;
    string targetState?;
};




json lifeCycleStateTransitions = {
  "States": [
    {
      "State": "Created",
      "Transitions": [
        {
          "event": "Publish",
          "targetState": "Published"
        },
        {
          "event": "Deploy as a Prototype",
          "targetState": "Prototyped"
        }
      ]
    },
    {
      "State": "Prototyped",
      "Transitions": [
        {
          "event": "Publish",
          "targetState": "Published"
        },
        {
          "event": "Demote to Created",
          "targetState": "Created"
        },
        {
          "event": "Deploy as a Prototype",
          "targetState": "Prototyped"
        }
      ]
    },
    {
      "State": "Published",
      "Transitions": [
        {
          "event": "Block",
          "targetState": "Blocked"
        },
        {
          "event": "Deploy as a Prototype",
          "targetState": "Prototyped"
        },
        {
          "event": "Demote to Created",
          "targetState": "Created"
        },
        {
          "event": "Deprecate",
          "targetState": "Deprecated"
        },
        {
          "Event": "Publish",
          "targetState": "Published"
        }
      ]
    },
    {
      "State": "Blocked",
      "Transitions": [
        {
          "event": "Deprecate",
          "targetState": "Deprecated"
        },
        {
          "event": "Re-Publish",
          "targetState": "Published"
        }
      ]
    },
    {
      "State": "Deprecated",
      "Transitions": [
        {
          "event": "Retire",
          "targetState": "Retired"
        }
      ]
    },
    {
      "State": "Retired"
    }
  ]
};

