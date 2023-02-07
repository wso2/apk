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

app.use(function(req, res, next){
  var err = req.session.error;
  var msg = req.session.success;
  delete req.session.error;
  delete req.session.success;
  res.locals.message = '';
  if (err) res.locals.message = '<p class="msg error">' + err + '</p>';
  if (msg) res.locals.message = '<p class="msg success">' + msg + '</p>';
  next();
});

//Initial request coming here and redirect to login
app.get('/', function(req, res){
  req.session.sessionDataKey = "testKey";
  res.redirect('/login');
});

// Logout function
app.get('/logout', function(req, res){
  // destroy the user's session to log them out
  // will be re-created next request
  req.session.destroy(function(){
    res.redirect('/');
  });
});

// Login GET request
app.get('/login', function(req, res) {
  // request sessiondata key add to cookie
  var minute = 60000;
  //Read cookie "req.cookies.sessionDataKey"
  //Read queryParam "req.param.sessionDataKey"
  console.log(req.cookies.sessionDataKey);
  if (req.session.sessionDataKey) {
    res.cookie('sessionDataKey', 1, { maxAge: minute });
    res.render('login');
  }
  else {
    req.session.error = 'Auth 302 error';
    res.redirect('loginError');
  }
});

//Login post request
app.post('/login', function (req, res, next) {
  // redirection IDP
  res.redirect(`url?username=${req.body.username}&password=${req.body.password}&organization=${req.body.org}`)
});

//Login callback
app.get('/login-callback', function (req, res, next) {
  
});


/* Listen Port */
if (!module.parent) {
  app.listen(3000);
  console.log('Express started on port 3000');
}
