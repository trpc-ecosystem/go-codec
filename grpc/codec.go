// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"trpc.group/trpc-go/trpc-go/codec"
)

type contextHeader struct{}

// ContextKeyHeader key in context to store Header
var ContextKeyHeader = &contextHeader{}

var (
	// DefaultServerCodec 默认编解码实例
	DefaultServerCodec = &ServerCodec{}
	// DefaultClientCodec 默认的客户端编解码器
	DefaultClientCodec = &ClientCodec{}
)

// init 注册 grpc codec 与 grpc server transport
func init() {
	codec.Register("grpc", DefaultServerCodec, DefaultClientCodec)
}

// ServerCodec 服务端编解码
type ServerCodec struct {
}

// Decode ServerCodec.Decode 用于解码
func (s *ServerCodec) Decode(msg codec.Msg, reqbuf []byte) (reqbody []byte, err error) {
	return reqbuf, nil
}

// Encode ServerCodec.Encode 用于编码
func (s *ServerCodec) Encode(msg codec.Msg, reqbuf []byte) (reqbody []byte, err error) {
	return reqbuf, nil
}

// ClientCodec 是 grpc 客户端的编解码器，什么都不做
type ClientCodec struct{}

// Encode 是 grpc 客户端的编码器，什么都不做
func (c *ClientCodec) Encode(msg codec.Msg, rspbody []byte) (buffer []byte, err error) {
	return rspbody, nil
}

// Decode 是 grpc 客户端的解码器，什么都不做
func (c *ClientCodec) Decode(msg codec.Msg, buffer []byte) (rspbody []byte, err error) {
	return buffer, nil
}
