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

// Package grpc tRPC-Go grpc 协议
package grpc

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/transport"
)

// GrpcToTrpcer is an interface to represent handler of grpc server.
type GrpcToTrpcer interface {
	Handle(srv interface{}, ctx context.Context, dec func(interface{}) error,
		interceptor grpc.UnaryServerInterceptor) (out interface{}, err error)
}

// GrpcToTrpcLayer implements GrpcToTrpcer and offers a handler of grpc server.
type GrpcToTrpcLayer struct {
	Handler transport.Handler
}

// Handle req and resp are passed to trpc-go through ctx, and req is obtained from ctx and written to resp in the
//
//	generated stub without repeated serialization.
//
// Obtaining the input and output types of the method from GrpcRegisterInfoMap is still indispensable. If this piece
//
//	of code can be placed in the stub, there is no need to record the input and output types.
func (g *GrpcToTrpcLayer) Handle(srv interface{}, ctx context.Context, dec func(interface{}) error,
	_ grpc.UnaryServerInterceptor) (out interface{}, err error) {
	method, _ := grpc.Method(ctx)
	ctx, msg := codec.WithNewMessage(ctx)
	msg.WithServerRPCName(method)

	if pr, ok := peer.FromContext(ctx); ok {
		if addr := pr.Addr.String(); addr != "" {
			msg.WithRemoteAddr(pr.Addr)
		}
	}

	index := strings.LastIndex(method, "/")
	if index < 0 {
		return nil, fmt.Errorf("GrpcToTrpcLayer: method: `%s` format error. ", method)
	}

	serviceName := method[1:index]
	methodName := method[index+1:]

	registerInfo, ok := grpcRegisterInfo[serviceName]
	if !ok {
		return nil, fmt.Errorf("serviceName: %s has not been registered. ", serviceName)
	}

	methodInfo, ok := registerInfo.MethodsInfo[methodName]
	if !ok {
		return nil, fmt.Errorf("methodName: %s has not been registered. ", methodName)
	}

	req := reflect.New(methodInfo.ReqType).Interface()
	if err = dec(req); err != nil {
		return nil, err
	}

	// put Header struct to ctx
	grpcData := &Header{
		Req:         req,
		InMetadata:  map[string][]string{},
		OutMetadata: metadata.MD{},
	}
	if gmd, ok := metadata.FromIncomingContext(ctx); ok {
		grpcData.InMetadata = gmd
	}

	innerCtx := context.WithValue(ctx, ContextKeyHeader, grpcData)

	var reqbuffer []byte
	if _, err := g.Handler.Handle(innerCtx, reqbuffer); err != nil {
		return nil, err
	}
	if err := msg.ServerRspErr(); err != nil && err.Code != errs.RetOK {
		return nil, err
	}
	// send md
	if grpcData.OutMetadata != nil {
		err = grpc.SendHeader(ctx, grpcData.OutMetadata)
		if err != nil {
			return nil, err
		}
	}

	return grpcData.Rsp, nil
}

// StreamHandler Encapsulate trpc.Handler as grpcHandler
func StreamHandler(srv interface{}, s grpc.ServerStream) error {
	ctx := s.Context()
	method, _ := grpc.Method(ctx)
	ctx, msg := codec.WithNewMessage(ctx)
	index := strings.LastIndex(method, "/")
	if index < 0 {
		return fmt.Errorf("GrpcToTrpcLayer：method: `%s` format error. ", method)
	}

	serviceName := method[1:index]
	methodName := method[index+1:]
	msg.WithServerRPCName(method)
	if pr, ok := peer.FromContext(ctx); ok {
		if addr := pr.Addr.String(); addr != "" {
			msg.WithRemoteAddr(pr.Addr)
		}
	}
	streamMap, ok := grpcRegisterInfo[serviceName]
	if !ok {
		return errs.NewFrameError(errs.RetServerNoFunc, fmt.Sprintf("not registered service: %s", serviceName))
	}
	desc, ok := streamMap.StreamsInfo[methodName]
	if !ok {
		return errs.NewFrameError(errs.RetServerNoFunc, fmt.Sprintf("not registered method: %s", msg.CalleeMethod()))
	}

	return desc.Handler(srv, s)
}

// makeGrpcDesc maps stream.ClientStreamDesc to grpc.StreamDesc
func makeGrpcDesc(desc *client.ClientStreamDesc) *grpc.StreamDesc {
	grpcDesc := &grpc.StreamDesc{
		StreamName:    desc.StreamName,
		ServerStreams: desc.ServerStreams,
		ClientStreams: desc.ClientStreams,
	}
	return grpcDesc
}
