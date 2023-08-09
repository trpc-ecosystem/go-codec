// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/transport"
)

func init() {
	transport.RegisterClientTransport("grpc", DefaultClientTransport)
}

// DefaultClientTransport 默认的客户端通讯层
var DefaultClientTransport = &clientTransport{}

// clientTransport 实现了trpc-go的transport.ClientTransport接口，使用原生grpc通讯层替代trpc-go通讯层
type clientTransport struct {
	connectionPool pool
	streamClient   grpc.ClientStream
	streamDesc     *RegisterStreamsInfo
}

// RoundTrip 是实现transport.ClientTransport的方法，调用原生grpc客户端代码
func (c *clientTransport) RoundTrip(ctx context.Context, req []byte,
	roundTripOpts ...transport.RoundTripOption) (rsp []byte, err error) {
	// 从ctx中获取grpc Header，用来获取请求和设置响应
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return nil, errs.NewFrameError(errs.RetClientValidateFail,
			fmt.Sprintf("grpc header in context cannot be transfered to grpc.Header"))
	}
	reqbody := header.Req
	rspbody := header.Rsp
	// 默认值
	opts := &transport.RoundTripOptions{}

	// 将传入的func option写到opts字段中
	for _, o := range roundTripOpts {
		o(opts)
	}

	msg := codec.Message(ctx)
	// 获取超时设置
	timeout := msg.RequestTimeout()
	// 从ctx中获取metadata并调用grpc方法设置客户端metadata
	ctx, err = setGRPCMetadata(ctx, msg)
	if err != nil {
		return nil, err
	}

	// 从服务器获取metadata
	md := &metadata.MD{}
	var callOpts []grpc.CallOption
	callOpts = append(callOpts, grpc.Header(md))

	// 从连接池获取grpc连接
	conn, err := c.connectionPool.Get(opts.Address, timeout)
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientConnectFail, err.Error())
	}
	// 使用grpc客户端调用远端服务器方法
	if err = conn.Invoke(ctx, msg.ClientRPCName(),
		reqbody, rspbody, callOpts...); err != nil {
		return nil, fmt.Errorf("grpc invoke failed. err: %v", err)
	}

	// 将服务器的metadata写入ctx，使上层能够获取
	header.InMetadata = *md

	return nil, nil
}

// setGRPCMetadata 将grpc的Header信息塞入到 metadata中
func setGRPCMetadata(ctx context.Context, msg codec.Msg) (context.Context, error) {
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return nil, errs.NewFrameError(errs.RetClientValidateFail,
			fmt.Sprintf("grpc header disappeared when set md, code error"))
	}
	// 将grpc md设置到ctx，供发送端使用
	var kv []string
	for k, vals := range header.OutMetadata {
		for _, v := range vals {
			kv = append(kv, k, v)
		}
	}
	if kv != nil {
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
	}

	return ctx, nil
}
