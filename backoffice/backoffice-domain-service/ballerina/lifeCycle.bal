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




final readonly & json lifeCycleStateTransitions = {
  "States": [
    {
      "State": "Created",
      "Transitions": [
        {
          "event": "Published",
          "targetState": "Published"
        },
        {
          "event": "Prototyped",
          "targetState": "Prototyped"
        }
      ]
    },
    {
      "State": "Prototyped",
      "Transitions": [
        {
          "event": "Published",
          "targetState": "Published"
        },
        {
          "event": "Created",
          "targetState": "Created"
        },
        {
          "event": "Prototyped",
          "targetState": "Prototyped"
        }
      ]
    },
    {
      "State": "Published",
      "Transitions": [
        {
          "event": "Blocked",
          "targetState": "Blocked"
        },
        {
          "event": "Prototyped",
          "targetState": "Prototyped"
        },
        {
          "event": "Created",
          "targetState": "Created"
        },
        {
          "event": "Deprecated",
          "targetState": "Deprecated"
        },
        {
          "Event": "Published",
          "targetState": "Published"
        }
      ]
    },
    {
      "State": "Blocked",
      "Transitions": [
        {
          "event": "Deprecated",
          "targetState": "Deprecated"
        },
        {
          "event": "Published",
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

