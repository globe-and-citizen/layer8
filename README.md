# Layer8 Setup Instructions

This repo contains the Layer8 Resource/Authentication Server, Proxy, and Service Provider Mocks.

In conjunction with the Layer8_Interceptor & the Layer8_Middleware (available throught npm), this repo forms the proof of concept for the Layer8 E2E encryption proxy to be made available through a WASM module in the browser.

The directories in the folder sp_mocks represent test cases and experimental projects to test various aspects of the system. These can be used as reference for development. 

Note: the tunnel is created everytime the page is reloaded. Therefore, Layer8 works best with single page applications. 

Currently the proof of concept only works with node.js in the backend and is tightly coupled to the express.js package.


## Setting Up a Layer8 Local Development Sever
### Requirements

- Docker with Docker Compose installed
- 2 GB Memory or more (to run PostgreSQL, InfluxDB, and Telegraf)
- Golang v1.21.1+
- NodeJS v20+ (for service provider mock)
- The Make utility to run `make <command>`.

### Get Started

1. Ensure Docker is running locally. 

2. Run the setup script to set up all dependencies (PostgreSQL, InfluxDB, Telegraf) with the following `make` command:

```bash
make setup_local_dependency
```

3. "Secrets" and "secret" configurations, such as credentials, can be added to a `.env.secret` file located in the server directory.


4. Run the Layer8 server:

```bash
make run_layer8server_local
```

If everything is set up correctly, you should be able to access the Layer8 server at localhost:5001. One Client (username: "layer8" & pass:  "12341234") and one User (username "tester" & pass "12341234") come already created for testing.

### Testing with Mock Service Providers

We have two mock service providers to demonstrate how our client will be using our Layer8 server. Here's how to set them up:

#### We've Got Poems Mock

This mock project is stored in the `sp_mocks/wgp` path.

##### Backend Setup

1. Navigate to the `sp_mocks/wgp/backend` directory.
2. Set the `.env` values to match those in `.dev.env`.
3. The default setup is to have WGP registered with the default layer8 Client (username: "layer8" & password: "12341234"). This Client has the following credentials hard coded onto lines 23 & 24 of `layer8/sp_mocks/wgp/backend/server.js`
  ```
    clientId: "f0fe2f3f-cabe-4d44-a8d4-ed7154627867",
    clientSecret: "b85bc41f06a0f86b912a51a9688a3c78fb464e9b7ac692eff145f20e4bcae3e8",
  ```
4. To add a different client, log in to the Layer8 Client Portal @ `http://localhost:5001/client-login-page` and register a new Client. Use this new clientId and clientSecret in your layer8auth object: 
  ```
    const layer8Auth = new ClientOAuth2({
      clientId: "<new UUID>",
      clientSecret: "<new Secret>",
      accessTokenUri: `${LAYER8_URL}/api/oauth`,
      authorizationUri: `${LAYER8_URL}/authorize`,
      redirectUri: LAYER8_CALLBACK_URL,
      scopes: ["read:user"],
    });
  ```

After everything is set up, run the backend with:

```bash
npm run dev
```

##### Frontend Setup

1. Navigate to the sp_mocks/wgp/frontend directory.
2. Set the .env values to match those in .dev.env.

After everything is set up, run the frontend with:

```bash
npm run dev
```


#### Imsharer Mock

Imsharer doesn't require additional other than registering as a new client with the correct backend URL. 
You can run it with the following commands:

1) Navigate to "http://localhost:5001/client-register-page"
2) Register the following suggested new Client:
  Project Name: "Image Sharer"
  Redirect URL: "n/a" 
  Backend URL: "localhost:6001"
  Username: "tester2"
  Password: "12341234" 


To run the frontend, navigate to the root directory:
```bash
make run_imsharer_frontend // run frontend
```

To run the backend, navigate to the root directory:
```bash
make run_imsharer_backend // run backend
```

### Setting Up Additional Configurations

The Layer8 server requires a configuration file stored in the `/server/.env` directory to run the application. We have provided a configuration example in the `.env.dev` file. If additional configuration is needed that should not be shared on GitHub, it should be placed in the `.env.secret` file. After running the setup script, the contents of `/server/.env.dev` and `/server/.env.secret` will be automatically copied to `/server/.env`.

Configuration for the cloud deployment will be stored in GitHub Action variables with the names DEVELOPMENT_APP_ENV and PRODUCTION_APP_ENV.

### Database Migration

#### Install required library

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

full documentation : 
[Installing golang-migrate library](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

#### Generate migration

example :

```bash
$ migrate create -ext sql -dir migrations -seq create_users_table
```

#### Running Migration

Each migration has up and down migration stored in [migrations](https://github.com/globe-and-citizen/layer8/tree/development/migrations)

```bash
$ migrate -database ${DB_URL} -path migrations up    Migration all in migrations
$ migrate -database ${DB_URL} -path migrations down  Revert migration all in migrations
$ migrate -database ${DB_URL} -path migrations up [N]    Apply N up migration
$ migrate -database ${DB_URL} -path migrations down [N]  Apply N down migration
$ migrate -database ${DB_URL} -path migrations version  Print current migration version
$ migrate -database ${DB_URL} -path migrations force V  Set version V but not run any migration
$ migrate -database ${DB_URL} -path migrations goto V  Migrate to version V
```

for full documentation see [migrate](https://github.com/golang-migrate/migrate)

## To Use the E2E Encrypted Tunnel:
### Frontend Code
#### To Open an Encrypted Tunnel:
After installing the NPM package, import the library and initialize an encrypted tunnel as follows (usually in the App.vue file):

```
// Top level imports
import layer8_interceptor from 'layer8_interceptor'

[...]

try{
  layer8_interceptor.initEncryptedTunnel({
    providers: [BACKEND_URLs], // An array of service providers you want to connect to
    proxy: PROXY_URL, // Your local proxy instance. Necessary only if the mode is "dev".
    "staticPath": "/anything", // The path to serve static files from.
  }, "dev") // If omitted, the default, hardcoded, production proxy is used. 
}catch(err){
  console.log(".initEncryptedTunnel error: ", err)
}
```

#### Standard API / Fetch Calls On the Frontend:
Once the encrypted tunnel is open, you can use the `layer8_interceptor.fetch()` method analogous to the native browser `fetch()`. Just like the native `fetch()`, a promise is returned and should be properly awaited within a asynchronous function.

```
// Top level imports
import layer8_interceptor from 'layer8_interceptor'

[...]

// Within asynchronous functions
let getResponse = await layer8_interceptor.fetch( URL )

let postResponse = await layer8_interceptor.fetch( URL, {
    method: "POST",
    headers: {
        "Content-Type": "Application/Json",
    },
    body: JSON.stringify({
        key: value,
    }),
});

```

#### Static Images and Other Assets:
```
// Top level imports
import layer8_interceptor from 'layer8_interceptor'

[...]

// Within asynchronous functions
const url = await layer8.static(image.url);

[...]

// Within the HTML
<img :src="url" />

```

### Backend Code
Use the Layer8 middleware as any other node package. At present, no configuration of the encrypted tunnel is necessary. Note: Currently Layer8 is being built tightly coupled with ExpresJS.

```
// Top Level Imports
const layer8 = require('layer8_middleware');


// Node Initialization
app.use(layer8.tunnel); // for the encrypted tunnel
app.use('/media', layer8.static('uploads')); // in order to connect a public folder, should be same as the `staticPath`

// Usage
app.get('/', (req, res) => {
  res.json({ key: val })
});

app.post("/route", async (req, res) => {
  const body = req.body;

  try {
    console.log(body)
    res.status(200).send("User registered successfully!");
  } catch (err) {
    console.log("err: ", err);
    res.status(500).send({ error: "Something went wrong!" });
  }
});
```

Note: During routine usage, there are no special calls necessary to make use of Layer8. The res.json() & res.send() & res.end() have been overwritten by Layer8 and will be used automatically. 

### Warnings and Gotcha's To Consider
1) Using express.json() as middleware in main file is unnecessary. The layer8_middleware automatically parses the incoming request. If you include the line app.use(express.json()) requests will get "caught" in the express.json() middleware and not reach your other endpoints.