'use strict'

/**
 * Module dependencies.
 */

var express = require('express');
var hash = require('pbkdf2-password')()
var path = require('path');
var session = require('express-session');
var cookieParser = require('cookie-parser');
var app = module.exports = express();

// config

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

// middleware
app.use(cookieParser());
app.use(express.urlencoded({ extended: false }))
app.use(session({
  resave: false, // don't save session if unmodified
  saveUninitialized: false, // don't create session until something stored
  secret: 'secret'
}));

// Session-persisted message middleware

app.use(function (req, res, next) {
  var err = req.session.error;
  var msg = req.session.success;
  delete req.session.error;
  delete req.session.success;
  res.locals.message = '';
  res.locals.loginURl = process.env.IDP_LOGIN_URL;
  if (err) res.locals.message = '<p class="msg error">' + err + '</p>';
  if (msg) res.locals.message = '<p class="msg success">' + msg + '</p>';
  next();
});

//Initial request coming here and redirect to login
app.get('/', function (req, res) {
  res.redirect('/login');
});


// Login GET request
app.get('/login', function (req, res) {
  var sessionKey = req.query.stateKey;
  res.locals.sessionKey = sessionKey;
  var sessionCookieName = "session-" + sessionKey;
  var sessionCookie = req.cookies[sessionCookieName];
  if (sessionCookie){
    res.render('login');
  }
});


//Login callback
app.get('/login-callback', function (req, res, next) {
  var stateKey = req.query.stateKey;
  var url = process.env.IDP_AUTH_CALLBACK_URL + "?sessionKey=" + stateKey;
  res.redirect(url);
});


/* Listen Port */
if (!module.parent) {
  app.listen(3000);
  console.log('Express started on port 3000');
}
