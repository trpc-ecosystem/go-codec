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

// DefaultClientTransport default client communication layer
var DefaultClientTransport = &clientTransport{}

// clientTransport Implemented the transport.ClientTransport interface of trpc-go, using the native grpc communication
//
//	layer instead of the trpc-go communication layer.
type clientTransport struct {
	connectionPool pool
	streamClient   grpc.ClientStream
	streamDesc     *RegisterStreamsInfo
}

// RoundTrip It is a method to implement transport.ClientTransport, calling the native grpc client code.
func (c *clientTransport) RoundTrip(ctx context.Context, req []byte,
	roundTripOpts ...transport.RoundTripOption) (rsp []byte, err error) {
	// Get the grpc Header from ctx to get the request and set the response
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return nil, errs.NewFrameError(errs.RetClientValidateFail,
			fmt.Sprintf("grpc header in context cannot be transfered to grpc.Header"))
	}
	reqbody := header.Req
	rspbody := header.Rsp
	// Defaults
	opts := &transport.RoundTripOptions{}

	// Write the incoming func option into the opts field
	for _, o := range roundTripOpts {
		o(opts)
	}

	msg := codec.Message(ctx)
	// Get timeout settings
	timeout := msg.RequestTimeout()
	// Get metadata from ctx and call grpc method to set client metadata
	ctx, err = setGRPCMetadata(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Get metadata from the server
	md := &metadata.MD{}
	var callOpts []grpc.CallOption
	callOpts = append(callOpts, grpc.Header(md))

	// Get a grpc connection from the connection pool
	conn, err := c.connectionPool.Get(opts.Address, timeout)
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientConnectFail, err.Error())
	}
	// Use the grpc client to call the remote server method
	if err = conn.Invoke(ctx, msg.ClientRPCName(),
		reqbody, rspbody, callOpts...); err != nil {
		return nil, fmt.Errorf("grpc invoke failed. err: %v", err)
	}

	// Write the metadata of the server to ctx, so that the upper layer can obtain
	header.InMetadata = *md

	return nil, nil
}

// setGRPCMetadata Insert grpc Header information into metadata
func setGRPCMetadata(ctx context.Context, msg codec.Msg) (context.Context, error) {
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return nil, errs.NewFrameError(errs.RetClientValidateFail,
			fmt.Sprintf("grpc header disappeared when set md, code error"))
	}
	// Set grpc md to ctx for use by the sender
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
