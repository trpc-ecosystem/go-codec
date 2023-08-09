// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"

	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/stream"
	"trpc.group/trpc-go/trpc-go/transport"
)

// DefaultStreamClient 生成新的StreamClient
var DefaultStreamClient = NewStreamClient()

// NewStreamClient  生成新的StreamClient
func NewStreamClient() stream.Client {
	return &StreamClient{}
}

// StreamClient grpc.Stream 客户端的实现
type StreamClient struct {
	connectionPool pool
}

// NewStream 生成streamConn并存储
func (s *StreamClient) NewStream(ctx context.Context, desc *client.ClientStreamDesc, method string,
	opt ...client.Option) (client.ClientStream, error) {
	cs := &clientStream{}
	cs.ctx = ctx
	msg := codec.Message(ctx)
	// 读取配置参数，设置用户输入参数
	opts, address, err := getOptions(msg, opt...)
	if err != nil {
		return nil, err
	}
	// 根据寻址选择器寻址到后端节点node
	if _, err = selectNode(msg, opts, address); err != nil {
		return nil, err
	}
	updateMsg(msg, opts)
	roundTripOpts := &transport.RoundTripOptions{}
	// 将传入的 call option 写到opts字段中
	for _, o := range opts.CallOptions {
		o(roundTripOpts)
	}
	timeout := msg.RequestTimeout()
	conn, err := s.connectionPool.Get(roundTripOpts.Address, timeout)
	if err != nil {
		return nil, err
	}
	grpcDesc := makeGrpcDesc(desc)
	stream, err := grpc.NewClientStream(ctx, grpcDesc, conn, method)
	if err != nil {
		return nil, err
	}
	cs.stream = stream
	return cs, nil
}

type clientStream struct {
	ctx    context.Context
	stream grpc.ClientStream
}

// RecvMsg 接收消息，返回error
func (cs *clientStream) RecvMsg(m interface{}) error {
	return cs.stream.RecvMsg(m)
}

// SendMsg  接收消息，返回error
func (cs *clientStream) SendMsg(m interface{}) error {
	return cs.stream.SendMsg(m)

}

// CloseSend 关闭发送那端，返回error
func (cs *clientStream) CloseSend() error {
	return cs.stream.CloseSend()
}

// Context 返回Context
func (cs *clientStream) Context() context.Context {
	return cs.ctx
}
