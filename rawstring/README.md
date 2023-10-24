English | [中文](README.zh_CN.md)

# rawstring

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/rawstring.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/rawstring)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/rawstring)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/rawstring)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawstring.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawstring.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/main/graph/badge.svg?flag=rawstring&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/main/rawstring)

The rawstring protocol is a simple tcp-based protocol.

Its data format is: string `"(a=b&c=d)\n"`.

The default ending is `"\n"`, and the response is `"result=0&...\n"`.
