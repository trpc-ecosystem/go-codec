English | [中文](README.zh_CN.md)

# go-codec

[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-green.svg)](https://github.com/trpc-ecosystem/go-codec/blob/main/LICENSE)

This repository provides implementations and examples of other protocols, currently including:

* grpc: Supports grpc protocol
* rawbinary: Supports the protocol where both request and response are `[]byte` under the udp transport layer protocol
* rawstring: Supports the protocol where both request and response are string types under the tcp transport layer protocol, with strings separated by `"\n"`
