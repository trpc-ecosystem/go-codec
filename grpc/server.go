// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"trpc.group/trpc-go/trpc-go/transport"
)

func init() {
	transport.RegisterServerTransport("grpc", DefaultServerTransport)
}

// DefaultServerTransport : 构建并封装 grpc server transport 实例
var DefaultServerTransport = NewServerTransport(transport.WithReusePort(true))

// ServerTransport 传输层
type ServerTransport struct {
	opts *transport.ServerTransportOptions
}

// NewServerTransport 创建 transport
func NewServerTransport(opt ...transport.ServerTransportOption) transport.ServerTransport {
	opts := &transport.ServerTransportOptions{}

	// 将传入的 func option 写到 opts 字段中
	for _, o := range opt {
		o(opts)
	}

	s := &ServerTransport{
		opts: opts,
	}

	return s
}

// ListenAndServe 处理配置
func (t *ServerTransport) ListenAndServe(ctx context.Context, opt ...transport.ListenServeOption) error {
	opts := &transport.ListenServeOptions{
		Network: "tcp",
	}
	for _, o := range opt {
		o(opts)
	}
	if opts.Handler == nil {
		return errors.New("trpc server transport handler empty")
	}

	var serveFuncIn = &GrpcToTrpcLayer{
		Handler: opts.Handler,
	}
	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	// 把所有的 grpc server 路由全部通过 GrpcToTrpcer 统一处理，并转发给 trpc-go 框架处理
	var methodDesc []grpc.MethodDesc
	for serviceName, serviceInfo := range grpcRegisterInfo {
		for methodName := range serviceInfo.MethodsInfo {
			methodDesc = append(methodDesc, grpc.MethodDesc{
				MethodName: methodName,
				Handler:    serveFuncIn.Handle,
			})
		}
		s.RegisterService(&grpc.ServiceDesc{
			ServiceName: serviceName,
			HandlerType: (*GrpcToTrpcer)(nil),
			Methods:     methodDesc,
			Metadata:    serviceInfo.Metadata,
		}, serveFuncIn)
	}
	reflection.Register(s)
	go s.Serve(lis)

	return nil
}
