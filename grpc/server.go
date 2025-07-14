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

// DefaultServerTransport : Construct and encapsulate the grpc server transport instance
var DefaultServerTransport = NewServerTransport(transport.WithReusePort(true))

// ServerTransport transport layer
type ServerTransport struct {
	opts *transport.ServerTransportOptions
}

// NewServerTransport create transport
func NewServerTransport(opt ...transport.ServerTransportOption) transport.ServerTransport {
	opts := &transport.ServerTransportOptions{}

	// Write the incoming func option into the opts field
	for _, o := range opt {
		o(opts)
	}

	s := &ServerTransport{
		opts: opts,
	}

	return s
}

// ListenAndServe process configuration
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
	// Process all grpc server routes through GrpcToTrpcer and forward them to the trpc-go framework for processing
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
