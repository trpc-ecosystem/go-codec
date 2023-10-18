English | [中文](README.zh_CN.md)

# go-codec

This repository provides implementation examples of some business protocols, currently including:

* grpc: Supports grpc protocol
* rawbinary: Supports the protocol where both request and response are `[]byte` under the udp transport layer protocol
* rawstring: Supports the protocol where both request and response are string types under the tcp transport layer protocol, with strings separated by `"\n"`
