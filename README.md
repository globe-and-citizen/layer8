# layer8
This repo contains the Layer8 Resource/Authentication Server, Proxy, and Service Provider Mocks.

In conjunction with the Layer8_Interceptor & the Layer8_Middleware (available throught npm), this repo forms the proof of concept for the Layer8 E2E encryption proxy to be made available through a WASM module in the browser.

The directories in the folder sp_mocks represent test cases and experimental projects to test various aspects of the system. These can be used as reference for development. 

Note: the tunnel is created everytime the page is reloaded. Therefore, Layer8 works best with single page applications. 

Currently the proof of concept only works with node.js in the backend and is tightly coupled to the express.js package. 

# To Run
1) Download and install Golang
2) Navigate to your project frontend directory (or the sp_mock you want to run). Run `npm install layer8_interceptor` (v0.0.17 at the time of writting). 
3) Naviate to your project backend directory (or the sp_mock you want to run). Run `npm install layer8_middleware` (v0.0.17 at the time of writting).
4) From the layer8 home directory, run:
    - `$cd ./server && go mod tidy`
    - `cd server && go run main.go -port=5001 -jwtKey=secret -MpKey=secret -UpKey=secret -ProxyURL=http://localhost:5001`
5) Clone `.env.dev` to `.env` in the `frontend` and `backend` directories of the sp_mock you are using.
6) Run your frontend / backend. If using that provided:
    - "We've Got Poems":
        - `cd sp_mocks/wgp/frontend && npm i && npm run dev`
        - `cd sp_mocks/wgp/backend && npm i && npm run dev`
    - "Image sharer":
        - `cd sp_mocks/imsharer/frontend && npm run dev`
        - `cd sp_mocks/imsharer/backend && npm run dev`

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

## Setup Metrics Collector

### Prerequisite
- Docker up and running
- Docker compose

### Setup

1. Start InfluxDB v2 as a Docker container by running the following command:
```
docker compose -f docker-compose-influx.yml up 
```
2. Open the InfluxDB dashboard via a browser using the defined credentials in docker-compose-influx.yml on port 8086.
3. Create an access token in the InfluxDB UI (https://docs.influxdata.com/influxdb/v2/admin/tokens/create-token/).
4. Add `INFLUXDB_URL` as `host.docker.internal` and `INFLUXDB_TOKEN` value as based on the created token variable to `.env` file to run the telegraf container.
5. Start Telegraf by running the following command:
```
docker compose -f docker-compose-telegraf.yml up 
```
6. After Telegraf is up and running, any metrics collected by the OpenTelemetry SDK could be sent via the gRPC protocol to port 4317.
7. For our case, set `OTEL_EXPORTER_OTLP_ENDPOINT` to localhost:4317.
