// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-go/server"
	"trpc.group/trpc-go/trpc-go/transport"
)

func Test_registerServices(t *testing.T) {
	s := grpc.NewServer()
	registerInfo := map[string]*RegisterInfo{
		"service1": {
			MethodsInfo: map[string]RegisterMethodsInfo{
				"method1": {},
			},
			StreamsInfo: map[string]server.StreamDesc{
				"stream1": {},
			},
		},
	}
	registerServices(s, registerInfo, nil)
}

func TestServerStream(t *testing.T) {
	s := NewServerStreamTransport(transport.WithIdleTimeout(100))
	ctx := context.Background()
	s.Send(ctx, []byte("abc"))
	s.Close(ctx)
	a := transport.WithHandler(transport.Handler(&h{}))
	s.ListenAndServe(ctx, a)
}
