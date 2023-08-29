// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"sync"
	"time"

	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-go/errs"
)

// pool 实现了简单的 grpc 连接池
type pool struct {
	connections sync.Map
}

// Get 从连接池中获取一个可用的 grpc 客户端连接
func (p *pool) Get(address string, timeout time.Duration) (*grpc.ClientConn, error) {
	// TODO 索引连接池时考虑超时时间
	if v, ok := p.connections.Load(address); ok {
		return v.(*grpc.ClientConn), nil
	}

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(), // TODO 从 ctx 中获取证书相关配置，支持 tls 通讯
		grpc.WithTimeout(timeout))
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientConnectFail, err.Error())
	}
	v, loaded := p.connections.LoadOrStore(address, conn)
	if !loaded {
		return conn, nil
	}
	return v.(*grpc.ClientConn), nil
}
