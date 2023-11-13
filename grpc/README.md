English | [中文](README.zh_CN.md)

# tRPC-Go grpc protocol

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/grpc.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/grpc)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/grpc)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/grpc)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/grpc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/grpc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/main/graph/badge.svg?flag=grpc&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/main/grpc)

The tRPC-Go framework achieves the purpose of supporting the grpc protocol through package introduction and grpc server encapsulation. It supports grpc server to process grpc client requests through grpc server transport and codec.

## Quick start

The following is the creation of a sample demo to demonstrate the usage process.

Suppose our current business project app is **test**, and the service server we want to develop is **hellogrpc**.

During the operation, you can set your own app and server name, but you need to pay attention to the replacement of the corresponding fields in the subsequent steps.

#### Preparation

1. An environment with a golang compilation environment.
4. [Install trpc tool](https://github.com/trpc-group/trpc-cmdline)
5. [Install grpc_cli tool](https://grpc.github.io/grpc/core/md_doc_command_line_tool.html)

#### Start

1. clone project: `git clone git@github.com:trpc-ecosystem/go-codec.git`

2. `cd trpc-codec/grpc/examples`

3. `mkdir hellogrpc && cd hellogrpc && mkdir protocol`

4.  init golang mod：`go mod init github.com/examples/hellogrpc`

5. On the protocol path, write the service agreement file `vim protocol/hellogrpc.proto`：

```proto
syntax = "proto3";  
package trpc.app.server;
option go_package="trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/streams";

message Req {
  string msg = 1;
}

message Rsp {
  string msg = 1;
}

service Greeter {
  rpc Hello(Req) returns (Rsp) {}
  rpc GetStream (Req) returns (stream Rsp){}
  rpc PutStream (stream Req) returns (Rsp){}
  rpc AllStream (stream Req) returns (stream Rsp){}
}
```

6. Generate a serving model via the command line: `trpc create --protocol=grpc --protofile=protocol/hellogrpc.proto --output .`.
7. In order to facilitate testing, replace the remote protocol with local `go mod edit -replace=trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc/protocol=./stub/trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc/protocol`

8. Write business logic:

- Modify `main.go`, add `trpc-grpc` package, and register in the main function:

> It will be modified to be directly supported by the trpc tool in the future, and it needs to be manually introduced and registered for the time being

```go
// Import library files
import "trpc.group/trpc-go/trpc-codec/grpc"
...
func main() {
      
    s := trpc.NewServer()
      
    pb.RegisterGreeterService(s, &greeterServiceImpl{})
      
    s.Serve()
}
```
      
- Modify the `greeter.go` file of the service interface, as follows:

```go
// Package main is the main package.
package main
      
import (
    "context"
      
    pb "trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc/protocol"
)
      
// SayHello ...
func (s *greeterServiceImpl) SayHello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
    // implement business logic here ...
    // new content
    rsp.Msg = "hello grpc client: " + req.Msg
      
    return nil
}
      
// SayHi ...
func (s *greeterServiceImpl) SayHi(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
    // implement business logic here ...
    // new content
    rsp.Msg = "hi grpc client: " + req.Msg
      
    return nil
}
```

9. Compile: `go build`, will generate the executable file of `hellogrpc`.

10. Modify the protocol field under `service` in the startup configuration `trpc_go.yaml` file under the current path, from `trpc` to `grpc`:

```yaml
  service:                                     # The service provided by the business service can have multiple
    - name: trpc.test.hellogrpc.Greeter      # service route name
      ip: 127.0.0.1                          # The service listens to the ip address. You can use the placeholder ${ip}, choose one of ip and nic, and give priority to ip
      #nic: eth0
      port: 8000                             # Service listening port can use placeholder ${port}
      network: tcp                           # Network monitoring type tcp/udp
      protocol: grpc                         # Change to grpc
      timeout: 1000                          # Request maximum processing time, at milliseconds
```

11. Start the service: `./hellogrpc &`

12. Execute tests with grpc-cli:

```shell
# view service
$ grpc_cli ls localhost:8000
grpc.reflection.v1alpha.ServerReflection
trpc.test.hellogrpc.Greeter
# View details of the Greeter service
$ grpc_cli ls localhost:8000 trpc.test.hellogrpc.Greeter -l
filename: hellogrpc.proto
package: trpc.test.hellogrpc;
service Greeter {
  rpc SayHello(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
  rpc SayHi(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
}
# See the details of the Greeter.SayHi method
$ grpc_cli ls localhost:8000 trpc.test.hellogrpc.Greeter.SayHi -l
rpc SayHi(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
# Debug Greeter.SayHi interface
$ grpc_cli call localhost:8000 'trpc.test.hellogrpc.Greeter.SayHi' "msg: 'I am a test.'"
msg: "hi grpc client: I am a test."
Rpc succeeded with OK status
```

13. Write client code
Client code generated using grpc-go.
```shell
# Generate client code for grpc-go
$ protoc --go_out=plugins=grpc:. protocol/hellogrpc.proto
```

14. Use the grpc-stream method

See [examples](/examples/README.md) for details

## Related References

[grpc protocol](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
[http2 frame](https://http2.github.io/http2-spec/#FramingLayer)
