## 运行示例代码前

1. 安装 [trpc](https://trpc.group/trpc-go/trpc-go-cmdline) 工具。
2. 生成桩代码：`cd ../testdata/protocols && make clean && make`。

## gRPC 调用 tRPC 服务

1. 启动 gRPC 协议的 tRPC 服务：`go run servers/tgrpc/main.go`。
2. 启动 gRPC client 调用 tRPC 服务：`go run clients/grpc/main.go`。

## tRPC 调用 gRPC 服务

1. 启动 gRPC 协议的 tRPC 服务：`go run servers/tgrpc/main.go`。
2. 启动 tRPC client 调用 tRPC 服务：`go run clients/tgrpc/main.go`。

## gRPC流式 调用 tRPC流式 服务

1. 启动 gRPC 协议的 tRPC 流式服务：`go run servers/tgrpc_stream/main.go`。
2. 启动 gRPC client 调用 tRPC 流式服务: `go run clients/tgrpc_stream/main.go`。

## 调用链

对于unary类型的连接，示例增加了trace链路的使用，目录分别是`clients/tgrpc/main.go`和`servers/tgrpc/main.go`
