// --------------------------------------------------------------------
//  Copyright (c) 2023, WSO2 LLC. (http://wso2.com) All Rights Reserved.

//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at

//  http://www.apache.org/licenses/LICENSE-2.0

//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
// -----------------------------------------------------------------------

const express = require("express");
const app = express();
const bodyParser = require("body-parser");
const jwt = require("jsonwebtoken");

// Middleware to parse JSON requests
app.use(bodyParser.json());

// POST route to handle JSON request
app.post("/api/v1/handle-request", (req, res) => {
  const data = req.body;
  console.log(data);

  if (data["requestHeaders"][":path"] === "/http-bin-api-basic/1.0.8/get") {
    res.send({
      rateLimitKeys: {
        rlkey_user: "bob",
      },
    });
  } else {
    res.send({
      rateLimitKeys: {},
    });
  }
});

// Start the server on port 8080
app.listen(8080, () => {
  console.log("Server is listening on port 8080");
});
