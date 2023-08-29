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
	transport.RegisterServerStreamTransport("grpc", DefaultServerStreamTransport)
}

// DefaultServerStreamTransport default server stream transport
var DefaultServerStreamTransport = NewServerStreamTransport()

// ServerStreamTransport 传输层
type ServerStreamTransport struct {
	opts *transport.ServerTransportOptions
}

// NewServerStreamTransport 创建 grpc_stream  transport
func NewServerStreamTransport(opt ...transport.ServerTransportOption) transport.ServerStreamTransport {
	opts := &transport.ServerTransportOptions{}
	// 将传入的 func option 写到 opts 字段中
	for _, o := range opt {
		o(opts)
	}

	s := &ServerStreamTransport{
		opts: opts,
	}

	return s
}

// ListenAndServe 启动 grpc 监听
func (t *ServerStreamTransport) ListenAndServe(ctx context.Context, opt ...transport.ListenServeOption) error {
	opts := &transport.ListenServeOptions{
		Network: "tcp",
	}
	for _, o := range opt {
		o(opts)
	}
	if opts.Handler == nil {
		return errors.New("trpc server transport handler empty")
	}

	s := grpc.NewServer()
	registerServices(s, grpcRegisterInfo, opts.Handler)
	reflection.Register(s)
	lis, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return err
	}
	go s.Serve(lis)

	return nil
}

// Send 执行发送逻辑，此处下层被 grpcShtreamHandle 接住，Send 未用到
func (t *ServerStreamTransport) Send(ctx context.Context, req []byte) error {
	return nil
}

// Close 服务端异常时候调用 Close 进行现场清理
func (t *ServerStreamTransport) Close(ctx context.Context) {
	return
}

// registerServices 注册 server 信息到 grpc.Server
func registerServices(s *grpc.Server, registerInfo map[string]*RegisterInfo, handler transport.Handler) {
	var serveFuncIn = &GrpcToTrpcLayer{
		Handler: handler,
	}
	for serviceName, serviceInfo := range grpcRegisterInfo {
		var (
			methodDesc []grpc.MethodDesc
			streamDesc []grpc.StreamDesc
		)
		for methodName := range serviceInfo.MethodsInfo {
			methodDesc = append(methodDesc, grpc.MethodDesc{
				MethodName: methodName,
				Handler:    serveFuncIn.Handle,
			})
		}
		for streamName, info := range serviceInfo.StreamsInfo {
			streamDesc = append(streamDesc, grpc.StreamDesc{
				StreamName:    streamName,
				Handler:       StreamHandler,
				ServerStreams: info.ServerStreams,
				ClientStreams: info.ClientStreams,
			})
		}
		s.RegisterService(&grpc.ServiceDesc{
			ServiceName: serviceName,
			HandlerType: serviceInfo.HandlerType,
			Methods:     methodDesc,
			Streams:     streamDesc,
			Metadata:    serviceInfo.Metadata,
		}, serviceInfo.ServerFunc)
	}
}
