require('dotenv').config();
var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
const fs = require('fs');
// Live reloading libs import
var livereload = require("livereload");
const connectLivereload = require("connect-livereload");
// Routes
var usersRouter = require('./routes/users');
// Initialize express
var app = express();

// Live reload browser
const liveReloadServer = livereload.createServer();
liveReloadServer.watch(path.join(__dirname, '../client/public'));
liveReloadServer.server.once("connection", () => {
  setTimeout(() => {
    liveReloadServer.refresh("/");
  }, 100);
});
// app.use(connectLivereload());
// Live reload browser end

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

// Server side route definitions

app.use('/users', usersRouter);
// app.use('/specs', specs);

// Serving the static react files
/* ******************************** */
/* ******************************** */
app.use(
  express.static(path.join(__dirname, "../client/public"))
);

app.get("*", (req, res) => {
  res.sendFile(
    path.join(__dirname, "../client/public/build/index.html")
  );
});

// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.render('error');
});

module.exports = app;
