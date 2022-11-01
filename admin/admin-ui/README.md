

## Configure

Open conf/Settings.js to configure UI side
Open .env to configure nodejs side of things.

## General

When nodejs server starts, it serves the web application as well as the nodejs server side endpoints.

### Running the server

Install node 14+

Go to the project root folder and execute the following command

```console
npm i
```

Then start the watch for client side code.

```console
npm run react-dev
```

Open another tab and watch for server side code changes

```console
npm run watch
```

Now base on the config provided visit the front end app. With the default port config, it's https://localhost:4000.
Also NodeJs routes can be accessed example: https://localhost:4000/users

