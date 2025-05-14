# DAPI Tools

## Mongoinfer

Basic tool to bootstrap from a MongoDB database. Tweaking the output may be necessary. Generates protobufs and a Dapi config with some CRUD operations. A separate step is needed to build a descriptors file from the protobufs using a tool like `buf`.

You can see an example of the output in the `example` directory.

### Quick Start

The example below assumes you have a mongodb instance running on localhost with your data. Running it will overwrite the example `yml` and `proto` in the `example` directory.

```
# Example run- customize the options as needed
go run cmd/mongoinfer/main.go --url=mongodb://localhost:27017 --dapi-config-file=example/config.yml --proto-file=example/out.proto --package com.example --service ExampleService

# Then you can use buf to generate a descriptors file
buf build -o out.pb

# For additional options see the help
go run cmd/mongoinfer/main.go --help
```

## Help

Join our Discord at https://discord.gg/hDjx3DehwG
