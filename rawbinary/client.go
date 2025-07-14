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

package rawbinary

import (
	"context"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
)

// Client interface.
type Client interface {
	Do(ctx context.Context, req []byte, opts ...client.Option) ([]byte, error)
}

// rawbinaryCli rawbinary Client
type rawbinaryCli struct {
	Client client.Client
	opts   []client.Option
	// calleeServiceName has four-stage, trpc.app.server.service.
	calleeServiceName string
	// rpcName is for reporting, two-stage, /rawbinary/interface1.
	rpcName string
}

// NewClientProxy creates a new rawbinaryCli proxy.
func NewClientProxy(opts ...client.Option) Client {
	c := &rawbinaryCli{
		Client: client.DefaultClient,
		opts: []client.Option{
			client.WithProtocol("rawbinary"),
			client.WithNetwork("udp"),
		},
	}
	c.opts = append(c.opts, opts...)
	return c
}

// NewClientProxyWithName can specify serviceName and rpcName.
// serviceName trpc.app.server.service has four-stage.
// rpcName has two-stage /rawbinary/interface1
func NewClientProxyWithName(serviceName, rpcName string, opts ...client.Option) Client {
	c := &rawbinaryCli{
		Client: client.DefaultClient,
		opts: []client.Option{
			client.WithProtocol("rawbinary"),
			client.WithNetwork("udp"),
		},
		calleeServiceName: serviceName,
		rpcName:           rpcName,
	}
	c.opts = append(c.opts, opts...)
	return c
}

// Do sends request.
func (c *rawbinaryCli) Do(ctx context.Context, req []byte, opts ...client.Option) ([]byte, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	msg.WithSerializationType(codec.SerializationTypeNoop)
	msg.WithCompressType(codec.CompressTypeNoop)
	msg.WithCalleeServiceName(c.calleeServiceName)
	msg.WithClientRPCName(c.rpcName)
	opt := append(c.opts, opts...)

	reqbody := &codec.Body{
		Data: req,
	}
	rspbody := &codec.Body{}

	err := c.Client.Invoke(ctx, reqbody, rspbody, opt...)
	return rspbody.Data, err
}

// Do sends request.
func Do(ctx context.Context, req []byte, opts ...client.Option) ([]byte, error) {
	c := NewClientProxy(opts...)
	return c.Do(ctx, req)
}
