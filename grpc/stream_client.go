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

	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/stream"
	"trpc.group/trpc-go/trpc-go/transport"
)

// DefaultStreamClient Generate a new StreamClient
var DefaultStreamClient = NewStreamClient()

// NewStreamClient  Generate a new StreamClient
func NewStreamClient() stream.Client {
	return &StreamClient{}
}

// StreamClient grpc.Stream client implementation
type StreamClient struct {
	connectionPool pool
}

// NewStream Generate streamConn and store
func (s *StreamClient) NewStream(ctx context.Context, desc *client.ClientStreamDesc, method string,
	opt ...client.Option) (client.ClientStream, error) {
	cs := &clientStream{}
	cs.ctx = ctx
	msg := codec.Message(ctx)
	// Read configuration parameters, set user input parameters
	opts, address, err := getOptions(msg, opt...)
	if err != nil {
		return nil, err
	}
	// Address to the backend node node according to the address selector
	if _, err = selectNode(msg, opts, address); err != nil {
		return nil, err
	}
	updateMsg(msg, opts)
	roundTripOpts := &transport.RoundTripOptions{}
	// Write the incoming call option into the opts field
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

// RecvMsg Receive message and return error
func (cs *clientStream) RecvMsg(m interface{}) error {
	return cs.stream.RecvMsg(m)
}

// SendMsg Receive message and return error
func (cs *clientStream) SendMsg(m interface{}) error {
	return cs.stream.SendMsg(m)

}

// CloseSend Close the sending end and return error
func (cs *clientStream) CloseSend() error {
	return cs.stream.CloseSend()
}

// Context Return to Context
func (cs *clientStream) Context() context.Context {
	return cs.ctx
}
