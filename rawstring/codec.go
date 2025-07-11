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

package rawstring

import (
	"trpc.group/trpc-go/trpc-go/codec"
)

func init() {
	codec.Register("rawstring", &serverCodec{}, &clientCodec{})
}

// serverCodec Server 端解码器
type serverCodec struct{}

// Decode 获取二进制请求数据
func (sc *serverCodec) Decode(msg codec.Msg, req []byte) ([]byte, error) {
	return req, nil
}

// Encode 回包二进制响应数据
func (sc *serverCodec) Encode(msg codec.Msg, rsp []byte) ([]byte, error) {
	return rsp, nil
}

// serverCodec Client 端解码器
type clientCodec struct{}

// Encode 打包二进制请求数据
func (cc *clientCodec) Encode(msg codec.Msg, reqBody []byte) ([]byte, error) {
	reqstr := string(reqBody)
	reqstr += "\n"
	data := []byte(reqstr)
	return data, nil
}

// Decode 解析二进制响应数据
func (cc *clientCodec) Decode(msg codec.Msg, rspBody []byte) ([]byte, error) {
	return rspBody, nil
}
