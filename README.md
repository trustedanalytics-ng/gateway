# WebSocket-based Gateway

The `gateway` is an ingestion bridge connecting edge devices with Cloud-based data processing. It exposes WebSocket interface which publishes all inbound messages to the configurable message back-end (Apache Kafka queue). It is compatible with any `Websocket-based` publisher.

Implemented:

* JSON message format, no mapping required
* Token-based client authentication (OAuth 2.0)
* Supports multiple backends (default to Apache Kafka)
* Configurable server (port, path, service etc.)

TODO:

* Client-level authorization
* Dynamic backend configuration

## Installation

You can install `gateway` by either cloning this repo (below) or by [downloading the latest binary](https://github.com/mchmarny/gateway/releases) distribution.

    git clone git@github.com:mchmarny/gateway.git
    cd ./gateway

To start the server invoke the `gateway` executable

    ./gateway

## Configuration

The `gateway` comes pre-configured with a default values.

```
args.ID = GetEnvVarAsString("GATEWAY_ID", "g1")
args.Index = GetEnvVarAsInt("GATEWAY_INDEX", 0)
args.Trace = GetEnvVarAsBool("GATEWAY_TRACE", false)

args.Server.Root = GetEnvVarAsString("GATEWAY_SERVER_ROOT", "/ws")
args.Server.Host = GetEnvVarAsString("GATEWAY_SERVER_HOST", "0.0.0.0")
args.Server.Port = GetEnvVarAsInt("GATEWAY_SERVER_PORT", 8080)
args.Server.Token = GetEnvVarAsString("GATEWAY_SERVER_TOKEN", "")
args.Server.AuthMethod = GetEnvVarAsString("GATEWAY_SERVER_AUTHMETHOD", "none")
args.Server.DeviceKeysURI = GetEnvVarAsString("GATEWAY_SERVER_DEVICEKEYSURI", "")
args.Server.TolerableJWTAge = GetEnvVarAsInt("GATEWAY_SERVER_TOLERABLEJWTAGE", 5)

var kafkaNodes string = GetEnvVarAsString("GATEWAY_PUB_URI", "docker:9091,docker:9092")

args.Pub.Topic = GetEnvVarAsString("GATEWAY_PUB_TOPIC", "messages")
args.Pub.Ack = GetEnvVarAsBool("GATEWAY_PUB_ACK", false)
args.Pub.Compress = GetEnvVarAsBool("GATEWAY_PUB_COMPRESS", true)
args.Pub.FlushFreq = GetEnvVarAsInt("GATEWAY_PUB_FLUSHFREQ", 1)
```

* `GATEWAY_PUB_TOPIC` will be automatically created if one does not exists
* `GATEWAY_PUB_ACK` if set to true will wait for acknowledgment from all brokers (slower)
* `GATEWAY_SERVER_AUTHMETHOD` can be one of `none`, `simple`, or `jwt`
* If you choose `none` as the authentication method, `gateway` will not attempt to authenticate any clients (all clients are authentic)
* If you wish to enable JWT authentication, set `GATEWAY_SERVER_AUTHMETHOD` environment variable to `jwt`. When JWT authentication is enabled, the environment variable `GATEWAY_SERVER_DEVICEKEYSURI` must be set to a GET REST API endpoint with the following properties:
  * the endpoint should have the following format:
    * `device_id` must be a path parameter and `alg` must be a URL parameter. An example format is http[s]://\<devices_repo_addr\>/devices/:device_id/key?alg=<alg>.
    * `alg` will hold values of the form 'ESXXX' or 'RSXXX'
  * the endpoint should return a JSON containing a `public_key` property with its value being the hex string containing the contents of the .pem public key file of the sending device with id `device_id`.

When JWT authentication is enabled, clients attempting to connect to `gateway` must include a JWT (in the Authorization header field) with the following properties:
* The payload must include two fields:
  * `device_id` with the device id of the client (a public key for that `device_id` must be available from the API endpoint set by `GATEWAY_SERVER_DEVICEKEYSURI`)
  * `iat` which is a Unix epoch time in seconds (visit [this link](https://tools.ietf.org/html/rfc7519#section-4.1.6) for details)
    * The acceptable age for a JWT can be changed (using unit of minutes) by setting `GATEWAY_SERVER_TOLERABLEJWTAGE` environment variable
* The JWT must be signed using an ES\* or RS\* algorithm.

## Testing

In addition to the integrated `Go` test, the `gateway` application also includes a simple Node.js test client: `etc/client.js`. This client is intended for perform a simple smoke-test of the deployed `gateway` application.

```
node client -s 'wss' \
            -e '${$APP_NAME}.<platform_domain>' \
            -p 4443 \
            -r '/ws' \
            -t '${APP_TOKEN}' \
            -f 1000
```

The test client will loop through and send to the gateway individual events at the frequency specified `-f` in milliseconds looking like this:

```
{
   "source_id": "ip-172-1-2-3",
   "event_id": "c4e6e427-7e45-41d8-9ba2-2223062398dc",
   "event_ts": 1417763522224,
   "metrics": [
      {
         "key": "cpu_load_5min",
         "value": 0.0029296875
      },
      {
         "key": "cpu_load_10min",
         "value": 0.0146484375
      },
      {
         "key": "cpu_load_15min",
         "value": 0.04541015625
      },
      {
         "key": "free_memory",
         "value": 3402670080
      }
   ],
   "event_number": 30
}
```

If you are not sure of the arguments, execute `node client.js --help` for some help.

## Backends

Currently `gateway` supports:

* [Apache Kafka](http://kafka.apache.org/)

#### Message

The `gateway` decorates the inbound messages with following attributes:

```
{
    id: [v4 uuid],
    on: [UTC timestamp]
    body: [inbound message content in UTF-8 encoded string]
}
```

## License

This project is under the MIT License. See the [LICENSE](https://github.com/mchmarny/gateway/blob/master/LICENSE) file for the full license text.
