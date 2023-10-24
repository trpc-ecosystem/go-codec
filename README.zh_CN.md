[English](README.md) | 中文

# go-codec

[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-codec/blob/main/LICENSE)

本仓库提供了部分业务协议的实现示例，目前包括：

* grpc: 支持 grpc 协议
* rawbinary: 在 udp 传输层协议下支持请求响应皆为 `[]byte` 的协议
* rawstring: 在 tcp 传输层协议下支持请求响应皆为 string 类型的协议，string 通过 `"\n"` 来进行分割
