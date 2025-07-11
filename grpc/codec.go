//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 Tencent.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

package grpc

import (
	"trpc.group/trpc-go/trpc-go/codec"
)

type contextHeader struct{}

// ContextKeyHeader key in context to store Header
var ContextKeyHeader = &contextHeader{}

var (
	// DefaultServerCodec Default codec instance
	DefaultServerCodec = &ServerCodec{}
	// DefaultClientCodec Default client codec
	DefaultClientCodec = &ClientCodec{}
)

// init Register grpc codec and grpc server transport
func init() {
	codec.Register("grpc", DefaultServerCodec, DefaultClientCodec)
}

// ServerCodec Server codec
type ServerCodec struct {
}

// Decode ServerCodec.Decode for decoding
func (s *ServerCodec) Decode(msg codec.Msg, reqbuf []byte) (reqbody []byte, err error) {
	return reqbuf, nil
}

// Encode ServerCodec.Encode for coding
func (s *ServerCodec) Encode(msg codec.Msg, reqbuf []byte) (reqbody []byte, err error) {
	return reqbuf, nil
}

// ClientCodec is the codec for the grpc client, does nothing
type ClientCodec struct{}

// Encode is the encoder for the grpc client and does nothing
func (c *ClientCodec) Encode(msg codec.Msg, rspbody []byte) (buffer []byte, err error) {
	return rspbody, nil
}

// Decode is a decoder for the grpc client, does nothing
func (c *ClientCodec) Decode(msg codec.Msg, buffer []byte) (rspbody []byte, err error) {
	return buffer, nil
}
