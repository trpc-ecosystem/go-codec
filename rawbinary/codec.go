// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// Package rawbinary is for UDP.
package rawbinary

import (
	"io"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/transport"
)

func init() {
	codec.Register("rawbinary", &serverCodec{}, &clientCodec{})
	transport.RegisterFramerBuilder("rawbinary", DefaultFrameBuilder)
}

// DefaultFrameBuilder is a default data frame constructor.
var DefaultFrameBuilder = &FrameBuilder{}

// FrameBuilder is a data frame constructor.
type FrameBuilder struct {
}

// New generates a dataframe.
func (fb *FrameBuilder) New(reader io.Reader) transport.Framer {
	return &framer{
		reader: reader,
	}
}

// framer is a dataframe.
type framer struct {
	reader io.Reader
}

// ReadFrame reads out the full data frame.
func (f *framer) ReadFrame() ([]byte, error) {
	return io.ReadAll(f.reader)
}

// serverCodec is a server-side decoder.
type serverCodec struct{}

// Decode gets binary request data.
func (sc *serverCodec) Decode(msg codec.Msg, req []byte) ([]byte, error) {
	return req, nil
}

// Encode returns binary response data.
func (sc *serverCodec) Encode(msg codec.Msg, rsp []byte) ([]byte, error) {
	return rsp, nil
}

// serverCodec is a client decoder.
type clientCodec struct{}

// Encode packs binary request data.
func (cc *clientCodec) Encode(msg codec.Msg, reqBody []byte) ([]byte, error) {
	return reqBody, nil
}

// Decode parses binary response data.
func (cc *clientCodec) Decode(msg codec.Msg, rspBody []byte) ([]byte, error) {
	return rspBody, nil
}
