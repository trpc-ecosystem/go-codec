English | [中文](README.zh_CN.md)

## Before running the sample code

1. Install [trpc](https://github.com/trpc-group/trpc-cmdline) tool, the version must be greater than v0.3.17.
2. Generate stub code: `cd ../testdata/protocols && make clean && make`.

## gRPC calls tRPC service

1. Start the tRPC service of the gRPC protocol: `go run servers/tgrpc/main.go`.
2. Start the gRPC client to call the tRPC service: `go run clients/grpc/main.go`.

## tRPC calls gRPC service

1. Start the tRPC service of the gRPC protocol: `go run servers/tgrpc/main.go`.
2. Start the tRPC client to call the tRPC service: `go run clients/tgrpc/main.go`.

## gRPC streaming call tRPC streaming service

1. Start the tRPC streaming service of the gRPC protocol: `go run servers/tgrpc_stream/main.go`.
2. Start the gRPC client to call the tRPC streaming service: `go run clients/tgrpc_stream/main.go`.

## call chain

For unary type connections, the example adds the use of trace links, and the directories are `clients/tgrpc/main.go` and `servers/tgrpc/main.go`
