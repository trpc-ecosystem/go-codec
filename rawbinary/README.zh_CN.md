[English](README.md) | 中文

# rawbinary

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/rawbinary.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/rawbinary)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/rawbinary)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/rawbinary)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawbinary.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawbinary.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/coverage/graph/badge.svg?flag=rawbinary&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/coverage/rawbinary)

基于 udp 的二进制协议，请求和响应均为 `[]byte`

## 示例

参考：examples/helloworld

## 注意事项

- rawbinary 只支持 udp
