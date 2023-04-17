var express = require('express');
var router = express.Router();
const axios = require('axios');
const https = require('https');
const jwt = require('jsonwebtoken');

const Settings = require('../../client/public/conf/Settings.js');

/* GET users listing. */
router.get('/update-token', async function (req, res, next) {
    try {
        res.send('respond with a resource');
    } catch (error) {
        console.log("Logging the error", error);
        next(error);
    }
});


module.exports = router;
