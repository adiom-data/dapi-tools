# DAPI Tools

## Mongoinfer

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
