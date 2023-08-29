// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package rawstring

import (
	"bufio"
	"io"

	"trpc.group/trpc-go/trpc-go/transport"
)

func init() {
	transport.RegisterFramerBuilder("rawstring", DefaultFramerBuilder)
}

var (
	// DefaultFramerBuilder cmd 默认数据帧构造器
	DefaultFramerBuilder = &FramerBuilder{}
)

// FramerBuilder cmd数据帧构造器
type FramerBuilder struct{}

// New 生成一个cmd数据帧
func (fd *FramerBuilder) New(reader io.Reader) transport.Framer {
	return &framer{
		reader: reader,
	}
}

// framer framer
type framer struct {
	reader io.Reader
}

// ReadFrame 从 io reader 中取出完整的数据帧
func (f *framer) ReadFrame() (msg []byte, err error) {
	reader := bufio.NewReader(f.reader)
	return reader.ReadBytes('\n')
}
