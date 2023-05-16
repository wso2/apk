var express = require('express');
var router = express.Router();
const axios = require('axios');
const https = require('https');
const Settings = require('../../client/public/conf/Settings.js');

/* GET users listing. */
router.get('/', async function (req, res, next) {
    const accessToken = req.cookies.access_token;

    const instance = axios.create({
        httpsAgent: new https.Agent({
            rejectUnauthorized: false
        })
    });
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
    };
    console.log("*****************");
    console.log(headers);
    console.log("*****************");
    try {
        const response = await instance.get(`${Settings.app.rest_api}/application-rate-plans`, {
            headers: headers
        });
        const data = response.data;
        console.log(data);
        res.status(200).json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({ error: error.message });
    }
    const applicationThrottlePoliciesx = {
        "count": 2,
        "list": [
            {
                "defaultLimit": {
                    "type": "REQUESTCOUNTLIMIT",
                    "requestCount": {
                        "requestCount": 41,
                        "timeUnit": "min",
                        "unitTime": 1
                    }
                },
                "policyId": "cbee719f-ea93-4578-91a5-df94df99a008",
                "policyName": "42PerMin",
                "displayName": "41PerMin",
                "description": "Allows 30 request per minute",
                "isDeployed": false,
                "type": "ApplicationThrottlePolicy"
            },
            {
                "defaultLimit": {
                    "type": "REQUESTCOUNTLIMIT",
                    "requestCount": {
                        "requestCount": 32,
                        "timeUnit": "min",
                        "unitTime": 1
                    }
                },
                "policyId": "d94c27e8-3867-482d-8218-9dc69e027ebe",
                "policyName": "32PerMin",
                "displayName": "32PerMin",
                "description": "Allows 32 request per minute",
                "isDeployed": false,
                "type": "ApplicationThrottlePolicy"
            }
        ]
    };

    // Retrieve the access token from the HTTP-only cookie

    // Make a request to the API server with the access token
    // Assuming the API server is located at https://api.example.com/items

});

router.get('/am/admin/api-categories', async function (req, res, next) {
    const accessToken = req.cookies.access_token;

    const instance = axios.create({
        httpsAgent: new https.Agent({
            rejectUnauthorized: false
        })
    });
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
        'Host': 'api.am.wso2.com'
    };
    try {
        const response = await instance.get(`${Settings.app.rest_api}/api-categories`, {
            headers: headers
        });
        const data = response.data;
        res.status(200).json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({ error: error.message });
    }
});

router.post('/am/admin/api-categories', async function (req, res, next) {
    const accessToken = req.cookies.access_token;

    const instance = axios.create({
        httpsAgent: new https.Agent({
            rejectUnauthorized: false
        })
    });
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
        'Host': 'api.am.wso2.com'
    };

    try {
        const response = await instance.post(`${Settings.app.rest_api}/api-categories`, req.body, {
            headers: headers
        });
        const data = response.data;
        res.status(201).json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({ code: error.code, error: error.message });
    }
});

router.put('/am/admin/api-categories/:categoryId', async function (req, res, next) {
    const categoryId = req.params.categoryId;
    const accessToken = req.cookies.access_token;

    const instance = axios.create({
        httpsAgent: new https.Agent({
            rejectUnauthorized: false
        })
    });
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
        'Host': 'api.am.wso2.com'
    };

    try {
        const response = await instance.put(`${Settings.app.rest_api}/api-categories/${categoryId}`, req.body, {
            headers: headers
        });
        const data = response.data;
        res.status(200).json(data);
    } catch (error) {
        console.log(error);
        res.status(500).json({ code: error.code, error: error.message });
    }
});

router.delete('/am/admin/api-categories/:categoryId', async function (req, res, next) {
    const categoryId = req.params.categoryId;
    const accessToken = req.cookies.access_token;

    const instance = axios.create({
        httpsAgent: new https.Agent({
            rejectUnauthorized: false
        })
    });
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${accessToken}`,
        'Host': 'api.am.wso2.com'
    };

    try {
        const response = await instance.delete(`${Settings.app.rest_api}/api-categories/${categoryId}`, {
            headers: headers
        });
        res.status(200).json({ message: "Successfully deleted" });
    } catch (error) {
        console.log(error);
        res.status(500).json({ code: error.code, error: error.message });
    }
});

module.exports = router;
