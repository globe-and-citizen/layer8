# Layer8

This repo contains the Layer8 Resource/Authentication Server, Proxy, and Service Provider Mocks.

In conjunction with the Layer8_Interceptor & the Layer8_Middleware (available throught npm), this repo forms the proof of concept for the Layer8 E2E encryption proxy to be made available through a WASM module in the browser.

The directories in the folder sp_mocks represent test cases and experimental projects to test various aspects of the system. These can be used as reference for development. 

Note: the tunnel is created everytime the page is reloaded. Therefore, Layer8 works best with single page applications. 

Currently the proof of concept only works with node.js in the backend and is tightly coupled to the express.js package.


## Requirements

- Docker with Docker Compose installed
- 2 GB Memory or more (to run PostgreSQL, InfluxDB, and Telegraf)
- Golang v1.21.1+
- NodeJS (for service provider mock)
- Ability to run `make <command>` for simplicity

## Get Started

1. Run the setup script to set up all dependencies (PostgreSQL, InfluxDB, Telegraf):

```bash
make setup_local_dependency
```


2. Run the Layer8 server:

```bash
make run_layer8server_local
```

If everything is set up correctly, you should be able to access the Layer8 server at localhost:5001, and have one client and one user created for testing.

## Setting Up the Configuration

Layer8 server requires a configuration file stored in the `/server/.env` directory to run the application. We have provided a configuration example in the .env.dev file. After running the setup script, it will automatically copy the contents of `/server/.env.dev` to `/server/.env.`

Configuration for the cloud deployment will be stored in GitHub Action variables with the names DEVELOPMENT_APP_ENV and PRODUCTION_APP_ENV.

## Database Migration

### Install required library

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

full documentation : 
[Installing golang-migrate library](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

### Generate migration

example :

```bash
$ migrate create -ext sql -dir migrations -seq create_users_table
```

### Running Migration

Each migration has up and down migration stored in [migrations](https://gitlab.com/m7310/user-management-service/-/tree/develop/migrations)

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

# To Use the E2E Encrypted Tunnel:
## Frontend Code
### To Open an Encrypted Tunnel:
After installing the NPM package, import the library and initialize an encrypted tunnel as follows (usually in the App.vue file):

```
// Top level imports
import layer8_interceptor from 'layer8_interceptor'

[...]

try{
  layer8_interceptor.initEncryptedTunnel({
    providers: [BACKEND_URLs], // An array of service providers you want to connect to
    proxy: PROXY_URL, // Your local proxy instance. Necessary only if the mode is "dev".
  }, "dev") // If omitted, the default is production mode. 
}catch(err){
  console.log(".initEncryptedTunnel error: ", err)
}
```

### Standard API / Fetch Calls On the Frontend:
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

### Static Images and Other Assets:
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

## Backend Code
Use the Layer8 middleware as any other node package. At present, no configuration of the encrypted tunnel is necessary. Note: Currently Layer8 is being built tightly coupled with ExpresJS.

```
// Top Level Imports
const layer8 = require('layer8_middleware');


// Node Initialization
app.use(layer8.tunnel); // for the encrypted tunnel
app.use('/media', layer8.static('uploads')); // in order to connect a public folder

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

## Warnings and Gotcha's
1) Using express.json() as middleware in main file is unnecessary. The layer8_middleware automatically parses the incoming request. If you include the line app.use(express.json()) requests will get "caught" in the express.json() middleware and not reach your other endpoints.


## Testing with Mock Service Providers

We have two mock service providers to demonstrate how our client will be using our Layer8 server. Here's how to set them up:

### We've Got Poems Mock

This mock project is stored in the `sp_mocks/wgp` path.

#### Backend Setup

1. Navigate to the `sp_mocks/wgp/backend` directory.
2. Set the `.env` values to match those in `.dev.env`.
3. Log in to the Layer8 Client Portal via `http://localhost:5001/client-login-page` using the `TEST_CLIENT_USERNAME` and `TEST_CLIENT_PASSWORD` provided in `/server/.env.dev`.
4. In the `server.js` file, on lines 23-24, enter the `clientId` and `clientSecret` values from the Layer8 Client Dashboard (UUID and secret).

After everything is set up, run the backend with:

```bash
npm run dev
```

#### Frontend Setup

1. Navigate to the sp_mocks/wgp/frontend directory.
2. Set the .env values to match those in .dev.env.

After everything is set up, run the frontend with:

```bash
npm run dev
```


### Imsharer Mock

Imsharer doesn't require additional setup. You can run it with the following commands:

To run the frontend:
```bash
make run_imsharer_frontend // run frontend
```

To run the backend:
```bash
make run_imsharer_backend // run backend
```
