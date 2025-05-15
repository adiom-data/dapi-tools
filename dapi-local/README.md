# DAPI

DAPI makes it really easy to build a declarative gRPC enabled access layer for your databases. Currently supports MongoDB and Postgres as a database and can perform authorization checks based on a JWT Authorization: Bearer token. It is built on top of [Connect RPC](https://connectrpc.com) so it supports JSON over HTTP in addition to gRPC.

* Enable easy and secure access to your databases through protobufs/gRPC
* Standardize API interfaces and communication in your stack on protobufs/gRPC

You can for example use this to set up a gRPC server to access your MongoDB databases. Then you can use [grpcmcp](https://github.com/adiom-data/grpcmcp) to make it available to Claude or another MCP enabled LLM client.

## Quick Start

* Start out in this directory
* Create a docker network: `docker create network dapi`
* Run a MongoDB Server on port 27017 (e.g. `docker run --name mongodb -p 27017:27017 --network dapi -d mongodb/mongodb-community-server:latest`)
* Run the server on port 8090 in a terminal `docker run -v "./config.yml:/config.yml" -v "./out.pb:/out.pb" -p 8090:8090 --network dapi -d markadiom/dapi`
* Issue some curl queries in another terminal:
```
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4ifQ.ha_SXZjpRN-ONR1vVoKGkrtmKR5S-yIjzbdCY0x6R3g" -H 'Content-Type: application/json' -d '{"data": {"fullplot": "hi"}}' localhost:8090/com.example.ExampleService/CreateTestMovies

curl -H 'Content-Type: application/json' -d '{}' localhost:8090/com.example.ExampleService/ListTestMovies
```
* You can also use grpcurl if you have it installed
```
grpcurl --plaintext -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYWRtaW4ifQ.ha_SXZjpRN-ONR1vVoKGkrtmKR5S-yIjzbdCY0x6R3g" -d '{"data": {"fullplot": "hi"}}' localhost:8090 com.example.ExampleService/CreateTestMovies

grpcurl --plaintext localhost:8090 com.example.ExampleService/ListTestMovies
```

Note that the example `CreateTestMovies` endpoint runs an authorization check that looks for the "role" of "admin" and so we pass an appropriate header whereas the `ListTestMovies` endpoint is unauthenticated and so does not require the header.

## Common Expression Language (CEL) Expressions

CEL is used in the config. See https://github.com/google/cel-spec for more information.

These are variables that *may* be available in the various CEL expressions used in the configuration.

* `req`
  * For a gRPC endpoint, this is a protobuf message representing the incoming request type.
* `claims`
  * This is a `map<string, any>` representing what was decoded from the claims portion of authorization. This may be empty if authorization did not happen.
* `headers`
  * This is a `map<string, list<string>>` representing the incoming headers.
* `resp`
  * This is likely *not* a protobuf message representing the outgoing response type. Rather, this is a nested map type that represents an intermediate result of a query to a backend database.

## Config

The configuration file is the core of running `dapi`. It provides some configuration for the server set up and chain of interceptors to apply. The main part is the `services` section which has declarations for special behavior for the specified gRPC endpoints found in `server.descriptors`.

For example, in the config below we configure the `/com.package.ServiceName/MethodName` endpoint to perform authorization by check that the claims role extracted from a Bearer JWT in the `Authorization` header is set to admin. Then if that passes the authorization check, the endpoint will perform a find in MongoDB.

Sample config.yml:
```
server:
  cleartext: true
  hostport: localhost:8080
  descriptors: protos.pb
interceptors:
  - name: Auth
    config:
      token: secretkey
  - name: MongoDB
    config:
      url: mongodb://localhost:27017
services:
  com.package.ServiceName:
    database: mongodatabase
    collection: mycollection
    endpoints:
      MethodName:
        auth: 'claims.role == "admin"'
        options: 'options.FindOptions {Limit: 3}'
        find:
          filter: '{ "name": req.name }'
      AnotherMethod:
        ...
  com.another.Service:
    ...
```

### Server

This section of the configuration governs the gRPC server itself.

* `hostport` string (e.g. localhost:8080)
  * **Required**. The host and port to bind the server to.

* `descriptors` string (e.g. /path/to/protos.pb)
  * **Required**. A path pointing to a file storing a binary encoding of your protobuf definitions. The service definitions in your config should also exist here. You can create this file via `protoc --descriptor_set_out=protos.pb` or using `buf build -o protos.pb`.

* `cleartext` bool
  * Run a cleartext server (no TLS). Do not set if specifying `certfile` and `keyfile`.

* `certfile` string (e.g. /path/to/mycert.crt)
  * Specify the certificate and key files using `certfile` and `keyfile` for TLS. Do not set `cleartext` if you are using this.
* `keyfile` string (e.g. /path/to/mykey.key)
  * See `certfile`.

* `reflect` bool
  * Run a gRPC reflection service.

* `health` bool
  * Run a gRPC health service.

* `shutdowntimeout` duration (e.g. 5s)
  * Timeout to apply when shutting down for a graceful exit.

### Interceptors

Provide the list of interceptors to apply in that general order and their configuration. There are 3 types of interceptors which we will label `http`, `grpc`, and `impl` as they apply at different parts of the stack. Within each type the ordering is as specified, but note that all `http` wraps around `grpc`, which wraps around `impl`. So to keep things organized just make sure to group the same type together in the list.

```
- name: <Interceptor Name - Case Sensitive>
  config:
    <Interceptor Config>
- ...
```

Including an interceptor will activate the respective part of the services config (if supported). For example, the `Auth` interceptor will activate the `auth` field in an endpoint configuration. The `CORS` interceptor doesn't use the service config at all.

Currently Available:
* `http` - This is just HTTP middleware with no gRPC context.
  * [CORS](#cors)
* `grpc` - This is what you'd normally understand as a gRPC interceptor.
  * [Auth](#auth)
  * [MongoDB](#mongodb)
  * [Postgres](#postgres)
* `impl` - This is a wrapper for the underlying implementation of an endpoint.
  * [Proxy](#proxy)

#### CORS

If you are exposing this server directly to web clients, you may want to enable support for CORS headers. Specify `allowedorigins` and other default configuration will also be provided.

* `allowedorigins` string list
  * **Required**. Must set at least this field if using CORS.
* `allowedheaders` string list
  * Any application specific headers cross-domain clients are allowed to send.
* `exposedheaders` string list
  * Any application specific headers cross-domain clients expect to receive.

#### Otel

Emit traces and metrics to an Otel collector. By default it uses the well known ports on localhost (port 4317 for gRPC or 4318 for HTTP).

* `servicename` string
  * The name of the service to use
* `trace`
  * `disabled` bool
    * Disable tracing
  * `endpointurl` string (e.g. `http://localhost:4317`, `http://localhost:4318`)
    * URL for the otel collector
  * `http` bool
    * Use HTTP instead of gRPC.
* `metric`
  * `disabled` bool
    * Disable metrics
  * `endpointurl` string (e.g. `http://localhost:4317`, `http://localhost:4318`)
    * URL for the otel collector
  * `http` bool
    * Use HTTP instead of gRPC.

#### Auth

The Auth interceptor performs authorization on a passed `Authorization` HTTP header. It looks for values of the form `Bearer <JWT>` and validates the JWT. Use only one of the configuration fields below.

* `token` string (e.g. mysecretkey)
  * A hard coded key used to validate the JWT.
* `jwksurl` string (e.g. http://localhost:8081/.well-known/jwks.json)
  * An URL pointing to a well known location hosting a JWKS endpoint for validating the JWT.

##### Services

* `unauthenticated` bool
  * If `true`, then any specified endpoints will not perform an authorization check and do not need to define any auth checks. This is overridden if an endpoint specifically defines an auth check.

##### Endpoints

* `unauthenticated` bool
  * If `true`, then this specific endpoint will not perform an authorization check and will ignore any defined auth checks.

* `auth` string
  * A CEL expression for checking authorization. It has access to `req`, `claims`, and `headers` and must return a `bool`. If this is defined for client and bidi streaming endpoints, this is checked on each client message.

* `streamclientauth` string
  * A CEL expression for checking authorization for client and bidi streaming endpoints only. It has access to `claims` and `headers` and must return a `bool`.

#### MongoDB

The MongoDB interceptor will activate MongoDB endpoint functionality for services. Usually this will be one of the later interceptors in the chain.

* `url` string (e.g. mongodb://localhost:27017)
  * The MongoDB URL.

##### Services

* `database` string
  * The default name of the database for all specified endpoints in the service.

* `collection` string
  * The default name of the collection for all specified endpoints in the service.

##### Endpoints

* `database` string
  * The name of the database to use for this endpoint.

* `collection` string
  * The name of the collection to use for this endpoint.

* `options` string (e.g. 'options.FindOptions { Limit: 2 }')
  * A CEL expression to specify options. It has access to `req`, `claims`, and `headers` and it must return the corresponding options type for the type of Mongo query defined. For example, if using a `find` query this must be `options.FindOptions { ... }`.

* `mapresponse` string
  * A CEL expression to specify a response mapping. It has access to `req`, `resp`, `claims`, and `headers` and it must return with the appropriate protobuf type for the endpoint or a generic map with a string key.

The following fields show the various available MongoDB queries. Most queries contain a single string field but some have two such as for updates. All of the string fields are CEL expressions that have access to `req`, `claims`, and `headers`. Each entry below shows the expected options type. The contents of the options and behavior follows that of the MongoDB Golang driver. [See Options Code](https://github.com/mongodb/mongo-go-driver/tree/master/mongo/options).

* `find: filter` string [`options.FindOptions`]
  * List data is nested in a "data" field.
* `findone: filter` string [`options.FindOneOptions`]
* `insertone: document` string [`options.InsertOneOptions`]
  * Inserted id is nested in a "_id" field.
* `insertmany: documents` string [`options.InsertManyOptions`]
  * Inserted ids is nested in a "_ids" field.
* `deleteone: filter` string [`options.DeleteOptions`]
  * Deleted count is nested in a "count" field.
* `deletemany: filter` string [`options.DeleteOptions`]
  * Deleted count is nested in a "count" field.
* `updateone: filter, update` string, string [`options.UpdateOptions`]
  * Nested "matched_count", "modified_count", "upserted_count", "upserted_id" are provided.
* `updatemany: filter, update` string, string [`options.UpdateOptions`]
  * Nested "matched_count", "modified_count", "upserted_count", "upserted_id" are provided.
* `replaceone: filter, replacement` string, string [`options.ReplaceOptions`]
  * Nested "matched_count", "modified_count", "upserted_count", "upserted_id" are provided.
* `findoneandreplace: filter, replacement` string, string [`options.FindOneAndReplaceOptions`]
* `findoneandupdate: filter, update` string, string [`options.FindOneAndUpdateOptions`]
* `findoneanddelete: filter` string [`options.FindOneAndDeleteOptions`]
* `aggregate: pipeline` string [`options.AggregateOptions`]
  * Data is nested in a "data" field.
* `paginatedfind: filter, cursor` string, string [`options.FindOptions`]
  * Special `find` implementation that does cursor based pagination. Specify `cursor` to be the field from `req` that contains a `bytes` or `string` cursor field (e.g. `req.cursor`). You can specify `Sort` and `Limit` (default: 100) in the find options. List data is nested in a "data" field. The next cursor is nested in a "next_cursor field. On the client side, to get the next page, simply pass the last "next_cursor" as the cursor in the next request until the "next_cursor" is empty.

##### CEL cheat sheet

* Example querying by Mongo ID from a hex string: `{ "_id": ObjectID(req._id) }`
  * `ObjectID is a special macro to support converting a hex string to a Mongo ID
* Sorting example `'options.FindOptions{ Sort: [kv("field1", 1), kv("field2" -1)]}'`
  * `kv` is a special macro to support sorting

#### Postgres

The Postgres interceptor will activate Postgres endpoint functionality for services. Usually this will be one of the later interceptors in the chain.

* `url` string (e.g. postgres://localhost:5432)
  * The Postgres DB URL.

##### Endpoints

* `exec` Query
  * Runs an exec command that returns "delete", "insert", "rows_affected", "select", "update", "string" fields reflecting the result of the exec.

* `query` Query
  * Runs a query command that returns a list of results nested in a "data" field.

* `queryone` Query
  * Runs a query command that returns the first result.

* `bulkinsert`
  * Inserts a list of items into a table
  * `table` string
    * **Required**. The table name.
  * `schema` string
    * The schema name.
  * `values` string
    * A CEL expression representing the items to insert. It should have a resulting type of `list<map<string, any>>` and each item must have the same columns. It has access to `req`, `claims`, and `headers`. It returns a "count" field indicating the number inserted.
  * `mapresponse` string
    * A CEL expression to specify a response mapping. It has access to `req`, `resp`, `claims`, and `headers` and it must return with the appropriate protobuf type for the endpoint or a generic map with a string key.

Example
```
query:
  sql: SELECT * from mytable where field=$1 and other=$2
  args:
    - req.field
    - req.field2
```

###### Query

* `sql` string
  * **Required**. A plain string that represents the SQL query. Use `$1`, `$2`, ... to indicate substitutions

* `args` string list
  * A list of CEL expressions that indicate the corresponding positional argument from `sql`. It has access to `req`, `claims`, and `headers`.

* `mapresponse` string
  * A CEL expression to specify a response mapping. It has access to `req`, `resp`, `claims`, and `headers` and it must return with the appropriate protobuf type for the endpoint or a generic map with a string key.

* `unsafesql` string
  * **WARNING**. A CEL expression to specify `sql` that can be used instead of `sql`. It has access to `req`, `claims`, and `headers`. Unsafe due to potential for SQL injection.

#### Proxy

The Proxy interceptor exists as a fallback to redirect to an implementation hosted elsewhere.

* `url` string (e.g. http://localhost:8091)
  * **Required**. The base url of the gRPC server to fall back to. `http` scheme will signify plain text requests, while `https` will use TLS.

* `gzip` bool
  * If set, use gzip when sending requests to the backend.

* `connect` bool
  * If set, use the connect protocol instead of gRPC.

* `requestheaders` string list
  * List of request header keys to forward.

* `responseheaders` string list
  * List of response header keys to forward.

* `responsetrailers` string list
  * List of response trailer keys to forward.

#### Services

* `proxy` bool
  * If set, all defined endpoints under will automatically be proxied.

#### Endpoints

* `proxy` bool
  * If set, this field will override the service level setting for the endpoint.

### Services

Services contains a map of a fully qualified gRPC service name to their respective configs.

```
com.service.ServiceName:
  ... # Service level config provided by Interceptors
  endpoints:
    MethodName:
      ... # Endpoint level config provided by Interceptors
    ... # More methods
... # More services
```

See the appropriate interceptor for the available configuration fields on both the Services and Endpoints level.

## Help

Join our Discord at https://discord.gg/hDjx3DehwG
