// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package rawstring

import (
	"context"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
)

// Client 接口
type Client interface {
	Do(ctx context.Context, req string, opts ...client.Option) (string, error)
}

// rawstingCli rawstring Client
type rawstingCli struct {
	Client client.Client
	opts   []client.Option
}

// NewClientProxy 新建一个 rawstingCli 代理
func NewClientProxy(opts ...client.Option) Client {
	c := &rawstingCli{
		Client: client.DefaultClient,
		opts: []client.Option{
			client.WithProtocol("rawstring"),
			client.WithNetwork("tcp"),
		},
	}
	c.opts = append(c.opts, opts...)
	return c
}

// Do 请求，按照\n进行分包
func (c *rawstingCli) Do(ctx context.Context, req string, opts ...client.Option) (string, error) {
	ctx, msg := codec.WithCloneMessage(ctx)
	msg.WithSerializationType(codec.SerializationTypeNoop)
	msg.WithCompressType((codec.CompressTypeNoop))

	opt := append(c.opts, opts...)
	byteReq := []byte(req)
	reqI := &codec.Body{
		Data: byteReq,
	}
	rspI := &codec.Body{}
	err := c.Client.Invoke(ctx, reqI, rspI, opt...)
	rsp := string(rspI.Data)
	return rsp, err
}

// Do 请求，按照\n进行分包
func Do(ctx context.Context, req string, opts ...client.Option) (string, error) {
	c := NewClientProxy(opts...)
	return c.Do(ctx, req)
}
