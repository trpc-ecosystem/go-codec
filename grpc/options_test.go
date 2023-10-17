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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/naming/selector"
)

type MockSelector struct {
}

func (s *MockSelector) Select(serviceName string, opt ...selector.Option) (*registry.Node, error) {
	switch serviceName {
	case "errorCase":
		return nil, errors.New("error")
	case "emptyCase":
		return &registry.Node{}, nil
	case "rightCase":
		return &registry.Node{Address: "address"}, nil
	case "address":
		return &registry.Node{
			ServiceName: "name",
			Address:     "address",
			Network:     "tcp",
		}, nil
	case "":
		return &registry.Node{
			ServiceName: "name",
			Address:     "address",
			Network:     "tcp",
		}, nil
	default:
		return nil, nil
	}
}

func (s *MockSelector) Report(node *registry.Node, cost time.Duration, err error) error {
	return nil
}

func Test_getNode(t *testing.T) {
	opts := &client.Options{}
	mockSelector := &MockSelector{}
	opts.Selector = mockSelector
	_, err := getNode("errorCase", opts)
	assert.NotNil(t, err)

	_, err = getNode("emptyCase", opts)
	assert.NotNil(t, err)

	node, err := getNode("rightCase", opts)
	assert.Equal(t, node.Address, "address")
	assert.Nil(t, err)
}

func Test_updateMsg(t *testing.T) {
	opts := &client.Options{
		ServiceName:       "ServiceName",
		CalleeMethod:      "CalleeMethod",
		CallerServiceName: "CallerServiceName",
		SerializationType: 1,
		CompressType:      1,
		ReqHead:           &Header{},
		RspHead:           &Header{},
	}
	ctx := context.Background()
	msg := codec.Message(ctx)
	updateMsg(msg, opts)
}

func Test_getOptions(t *testing.T) {
	ctx := context.Background()
	msg := codec.Message(ctx)
	msg.WithNamespace("namespace")
	_, address, err := getOptions(msg)
	assert.Equal(t, "", address)
	assert.Nil(t, err)
}

func Test_selectNode(t *testing.T) {
	opts := &client.Options{}
	opts.Selector = &MockSelector{}

	ctx := context.Background()
	msg := codec.Message(ctx)
	_, err := selectNode(msg, opts, "address")
	assert.Nil(t, err)
}

func Test_setNamingInfo(t *testing.T) {
	opts := &client.Options{}
	addr, err := setNamingInfo(opts)
	assert.Nil(t, err)
	assert.Equal(t, "", addr)

	opts.Target = "target"
	addr, err = setNamingInfo(opts)
	assert.NotNil(t, err)
	assert.Equal(t, "", addr)

	opts.Target = "mock://target"
	addr, err = setNamingInfo(opts)
	assert.NotNil(t, err)

	opts.Target = "mock://target"
	selector.Register("mock", &MockSelector{})

	addr, err = setNamingInfo(opts)
	assert.Nil(t, err)
	assert.Equal(t, "target", addr)
}
