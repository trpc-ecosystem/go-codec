[English](README.md) | 中文

# tRPC-Go grpc 协议

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/grpc.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/grpc)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/grpc)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/grpc)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/grpc.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/grpc.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/main/graph/badge.svg?flag=grpc&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/main/grpc)

tRPC-Go 框架通过包引入和 grpc server 的封装，来达到支持 grpc 协议的目的。它通过 grpc server transport 和编解码来支持 grpc server 处理 grpc client 的请求。

## 快速开始

以下是通过示例 demo 的创建，演示使用流程。

假设我们现在的业务项目 app 是 **test**，我们要开发的服务 server 是 **hellogrpc**。

使用的 git 工程是 `http://trpc.group/trpc-go/trpc-codec.git`，并将本示例放置在该工程下的 `grpc/examples` 路径下。

大家在操作的过程中，可以设置自己的 app 以及 server 名，但需要再后续步骤中，注意相应字段的替换。

#### 准备工作

1. 具备 golang 编译环境
4. [安装 trpc 工具](https://trpc.group/trpc-go/trpc-go-cmdline)
5. [安装 grpc_cli 工具](https://grpc.github.io/grpc/core/md_doc_command_line_tool.html)

#### 开始

1. clone 工程：`git clone git@github.com:trpc-ecosystem/go-codec.git`

2. `cd trpc-codec/grpc/examples`

3. `mkdir hellogrpc && cd hellogrpc && mkdir protocol`

4. 初始化 golang mod：`go mod init trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc`

5. 在 protocol 路径下，编写服务协议文件 `vim protocol/hellogrpc.proto`：

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

    > **注意 proto 中的 package 以及 go_package 的定义。**

6. 通过命令行生成服务模型：`trpc create --protocol=grpc --protofile=protocol/hellogrpc.proto --output .`。
7. 为了方便测试，替换远程协议成本地 `go mod edit -replace=trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc/protocol=./stub/trpc.group/trpc-go/trpc-codec/grpc/examples/hellogrpc/protocol`

8. 编写业务逻辑：

- 修改 `main.go`，添加 `trpc-grpc` package，并且在 main 函数中注册：

> 后续会修改成 trpc 工具直接支持，暂时需要手动引入和注册

```go
// 引入库文件
import "trpc.group/trpc-go/trpc-codec/grpc"
...
func main() {
    s := trpc.NewServer()
      
    pb.RegisterGreeterService(s, &greeterServiceImpl{})
      
    s.Serve()
}
```
      
- 修改 service 接口 `greeter.go` 文件，形如：

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
    // 新增内容
    rsp.Msg = "hello grpc client: " + req.Msg
      
    return nil
}
      
// SayHi ...
func (s *greeterServiceImpl) SayHi(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
    // implement business logic here ...
    // 新增内容
    rsp.Msg = "hi grpc client: " + req.Msg
      
    return nil
}
```

9. 编译：`go build`，将会生成 `hellogrpc` 的可执行文件。

10. 修改当前路径下的启动配置 `trpc_go.yaml` 文件中的 `service` 下的 protocol 字段，从 `trpc` 修改为 `grpc`：

```yaml
  service:                                         #业务服务提供的 service，可以有多个
      - name: trpc.test.hellogrpc.Greeter      #service 的路由名称
        ip: 127.0.0.1                            #服务监听 ip 地址 可使用占位符 ${ip},ip 和 nic 二选一，优先 ip
        #nic: eth0
        port: 8000                #服务监听端口 可使用占位符 ${port}
        network: tcp                             #网络监听类型  tcp udp
        protocol: grpc               #修改为 grpc
        timeout: 1000                            #请求最长处理时间 单位 毫秒
```

11. 启动服务：`./hellogrpc &`

12. 使用 grpc-cli 执行测试：

```shell
# 查看服务
$ grpc_cli ls localhost:8000
grpc.reflection.v1alpha.ServerReflection
trpc.test.hellogrpc.Greeter
# 查看 Greeter 服务的详细信息
$ grpc_cli ls localhost:8000 trpc.test.hellogrpc.Greeter -l
filename: hellogrpc.proto
package: trpc.test.hellogrpc;
service Greeter {
  rpc SayHello(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
  rpc SayHi(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
}
# 查看 Greeter.SayHi 方法的详细信息
$ grpc_cli ls localhost:8000 trpc.test.hellogrpc.Greeter.SayHi -l
rpc SayHi(trpc.test.hellogrpc.HelloRequest) returns (trpc.test.hellogrpc.HelloReply) {}
# 调试 Greeter.SayHi 接口
$ grpc_cli call localhost:8000 'trpc.test.hellogrpc.Greeter.SayHi' "msg: 'I am a test.'"
msg: "hi grpc client: I am a test."
Rpc succeeded with OK status
```

13. 编写客户端代码
使用 grpc-go 生成的客户端代码。
```shell
# 生成 grpc-go 的客户端代码
$ protoc --go_out=plugins=grpc:. protocol/hellogrpc.proto
```

14. 使用 grpc—stream 方式
详见 [examples](/examples/README.zh_CN.md)

## 相关参考

[grpc protocol](https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md)
[http2 frame](https://http2.github.io/http2-spec/#FramingLayer)
