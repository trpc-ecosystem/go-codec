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
	"sync"
	"time"

	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-go/errs"
)

// pool Implemented a simple grpc connection pool
type pool struct {
	connections sync.Map
}

// Get Obtain an available grpc client connection from the connection pool
func (p *pool) Get(address string, timeout time.Duration) (*grpc.ClientConn, error) {
	// TODO Consider timeouts when indexing connection pools
	if v, ok := p.connections.Load(address); ok {
		return v.(*grpc.ClientConn), nil
	}

	conn, err := grpc.Dial(address,
		grpc.WithInsecure(), // TODO Obtain certificate related configuration from ctx, support tls communication
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
