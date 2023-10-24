English | [中文](README.zh_CN.md)

# rawbinary 

[![Go Reference](https://pkg.go.dev/badge/trpc.group/trpc-go/trpc-codec/rawbinary.svg)](https://pkg.go.dev/trpc.group/trpc-go/trpc-codec/rawbinary)
[![Go Report Card](https://goreportcard.com/badge/trpc.group/trpc-go/trpc-codec/rawbinary)](https://goreportcard.com/report/trpc.group/trpc-go/trpc-codec/rawbinary)
[![Tests](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawbinary.yml/badge.svg)](https://github.com/trpc-ecosystem/go-codec/actions/workflows/rawbinary.yml)
[![Coverage](https://codecov.io/gh/trpc-ecosystem/go-codec/branch/main/graph/badge.svg?flag=rawbinary&precision=2)](https://app.codecov.io/gh/trpc-ecosystem/go-codec/tree/main/rawbinary)


The native binary protocol does not perform any processing on the protocol, and provides raw []byte requests and responses.

## Example

See examples/helloworld

## Precautions

- rawbinary only supports udp.
