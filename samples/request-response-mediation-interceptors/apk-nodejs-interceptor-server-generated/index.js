'use strict';

const fs = require('fs');
var path = require('path');
const https = require('https');
const morgan = require('morgan')

var oas3Tools = require('oas3-tools');
var serverPort = 9081;

// swaggerRouter configuration
var options = {
    routing: {
        controllers: path.join(__dirname, './controllers')
    },
};

var expressAppConfig = oas3Tools.expressAppConfig(path.join(__dirname, 'api/openapi.yaml'), options);
var app = expressAppConfig.getApp();
app.use(morgan('combined'))

const serverOptions = {
    cert: fs.readFileSync('certs/tls.crt'),
    key: fs.readFileSync('certs/tls.key'),
    requestCert: false,
    rejectUnauthorized: false,
    ca: fs.readFileSync('certs/ca.crt'),
}

// Initialize the Swagger middleware
https.createServer(serverOptions, app).listen(serverPort, function () {
    console.log('Your server is listening on port %d (http://localhost:%d)', serverPort, serverPort);
    console.log('Swagger-ui is available on http://localhost:%d/docs', serverPort);
});
