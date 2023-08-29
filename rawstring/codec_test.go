// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

/*
 * @Author: your name
 * @Date: 2020-09-17 14:55:23
 * @LastEditTime: 2020-09-17 14:56:46
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /rawstring/codec_test.go
 */
package rawstring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockClientCodec = &clientCodec{}
var mockServerCodec = &serverCodec{}

func Test_clientCodec_Decode(t *testing.T) {
	rsp, err := mockClientCodec.Decode(nil, []byte("mock"))
	assert.Nil(t, err)
	assert.Equal(t, rsp, []byte("mock"))
}

func Test_clientCodec_Encode(t *testing.T) {
	rsp, err := mockClientCodec.Encode(nil, []byte("mock"))
	assert.Nil(t, err)
	assert.Equal(t, rsp, []byte("mock\n"))
}

func Test_serverCodec_Decode(t *testing.T) {
	rsp, err := mockServerCodec.Decode(nil, []byte("mock"))
	assert.Nil(t, err)
	assert.Equal(t, rsp, []byte("mock"))
}

func Test_serverCodec_Encode(t *testing.T) {
	rsp, err := mockServerCodec.Encode(nil, []byte("mock"))
	assert.Nil(t, err)
	assert.Equal(t, rsp, []byte("mock"))
}
