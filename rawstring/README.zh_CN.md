[English](README.md) | 中文

# rawstring

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/rawstring.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/rawstring)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/rawstring)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/rawstring)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawstring.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawstring.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/main/graph/badge.svg?flag=rawstring&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/main/rawstring)

rawstring 协议是一种简单的基于 tcp 的调用协议

其数据格式为：字符串 `"(a=b&c=d)\n"`

结尾默认是 `"\n"`，响应回复 `"result=0&...\n"`
