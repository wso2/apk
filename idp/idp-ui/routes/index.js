'use strict'

/**
 * Module dependencies.
 */

var express = require('express');
var hash = require('pbkdf2-password')()
var path = require('path');
var session = require('express-session');

var app = module.exports = express();

// config

app.set('view engine', 'ejs');
app.set('views', path.join("/home/krish/Documents/apk/idp/idp-ui/", 'views'));

// middleware

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

app.get('/', function(req, res){
  req.session.sessionDataKey = "testKey";
  res.redirect('/login');
});


app.get('/logout', function(req, res){
  // destroy the user's session to log them out
  // will be re-created next request
  req.session.destroy(function(){
    res.redirect('/');
  });
});

app.get('/login', function(req, res) {
  // request sessiondata key add to cookie
  var minute = 60000;
  console.log(req.session.sessionDataKey);
  if (req.session.sessionDataKey) {
    res.cookie('sessionDataKey', 1, { maxAge: minute });
    res.render('login');
  }
  else {
    req.session.error = 'Auth 302 error';
    res.render('login');
  }
});

app.post('/login', function (req, res, next) {
  // redirection IDP
  res.redirect(`url?username=${req.body.username}&password=${req.body.password}&organization=${req.body.org}`)
});

app.get('/login-callback', function (req, res, next) {

});


/* istanbul ignore next */
if (!module.parent) {
  app.listen(3000);
  console.log('Express started on port 3000');
}