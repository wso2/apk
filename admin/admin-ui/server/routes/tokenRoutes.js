var express = require('express');
var router = express.Router();
const axios = require('axios');
const https = require('https');
const jwt = require('jsonwebtoken');

const Settings = require('../../client/public/conf/Settings.js');

/* GET users listing. */
router.get('/', async function (req, res, next) {
    process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0'; 
    const instance = axios.create({
        httpsAgent: new https.Agent({  
          rejectUnauthorized: false
        })
      });
      
    // Disables SSL verification
    // TODO we need to get this client_secret from an Environment variable
    /* eslint-disable no-undef */
    const base64EncodedKeyAndSecret = Buffer.from(`${Settings.idp.client_id}:${Settings.idp.client_secret}`).toString('base64');
    // /Send post request to authorization endpoint get the token
    const tokenRequestPayload = {
        grant_type: 'authorization_code',
        code: req.query.code,
        redirect_uri: `${Settings.idp.redirect_uri}/token`,
        client_id: Settings.idp.client_id
    };
    const headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': `Basic ${base64EncodedKeyAndSecret}`,
        'Host': Settings.idp.host,
    };
    try {
        console.log('sending token request to', Settings.idp.token_endpoint);
        const response = await instance.post(Settings.idp.token_endpoint, tokenRequestPayload, {
            headers: headers
        });
        const data = response.data;
        console.log(data);
        const { access_token, expires_in, refresh_token } = data;
        const cookieOptions = {
            httpOnly: true,
            maxAge: expires_in * 1000, // expires_in is the token expiration time in seconds
            sameSite: 'strict',
            secure: true, //process.env.NODE_ENV === 'production' // Set to true in production
        };
        res.cookie('access_token', access_token, cookieOptions);
        // decode the token and print the fields
        const decodedToken = jwt.decode(access_token, { complete: true });

        // Extract the fields from the decoded token
        const { header, payload } = decodedToken;
        const { iss, sub, exp } = payload;

        console.log('Header:', header);
        console.log('Issuer:', iss);
        console.log('Subject:', sub);
        console.log('Expiration Time:', new Date(exp * 1000));
        // redirect to the home page
        res.redirect(`/?user=${sub}&exp=${exp}`);
    } catch (error) {
        console.log("Logging the error", error);
        next(error);
    }
});


module.exports = router;
