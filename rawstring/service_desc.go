//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

package rawstring

import (
	"context"

	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/server"
)

// Server is the rawstring server interface.
type Server interface {
	ServiceName() string
	Handle(ctx context.Context, req []byte) ([]byte, error)
}

// Handler is the server-side handling function.
func Handler(svr interface{}, ctx context.Context, f server.FilterFunc) (interface{}, error) {
	msg := codec.Message(ctx)
	msg.WithSerializationType(codec.SerializationTypeNoop)
	msg.WithCompressType(codec.CompressTypeNoop)

	req := &codec.Body{}
	filters, err := f(req)
	if err != nil {
		return nil, err
	}

	handleFunc := func(ctx context.Context, reqI interface{}) (interface{}, error) {
		reqBody := reqI.(*codec.Body).Data
		return svr.(Server).Handle(ctx, reqBody)
	}

	rsp, err := filters.Filter(ctx, req, handleFunc)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// Register registers the service method.
func Register(s server.Service, svr Server) error {
	serviceDesc := server.ServiceDesc{
		ServiceName: svr.ServiceName(),
		HandlerType: (*Server)(nil),
		Methods: []server.Method{
			{
				Name: "*",
				Func: Handler,
			},
		},
	}
	return s.(*server.Server).Service(svr.ServiceName()).Register(&serviceDesc, svr)
}
