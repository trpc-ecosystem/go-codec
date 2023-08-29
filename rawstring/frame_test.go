// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package rawstring

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFramerBuilder_New(t *testing.T) {
	var b []byte
	f := DefaultFramerBuilder.New(bytes.NewReader(b))
	assert.NotNil(t, f)
}

func TestFramer_ReadFrame(t *testing.T) {
	b := []byte("hello\n")
	f := DefaultFramerBuilder.New(bytes.NewReader(b))
	msg, err := f.ReadFrame()
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
}
