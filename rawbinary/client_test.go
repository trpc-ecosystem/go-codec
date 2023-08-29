// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package rawbinary

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
)

type fakeClient struct {
}

// Invoke is an rpc call。
func (c *fakeClient) Invoke(ctx context.Context, reqbody interface{}, rspbody interface{}, opt ...client.Option) (err error) {
	body, ok1 := reqbody.(*codec.Body)
	rbody, ok2 := rspbody.(*codec.Body)
	if ok1 && ok2 {
		rbody.Data = make([]byte, len(body.Data))
		copy(rbody.Data, body.Data)
	}
	return nil
}

func TestDo(t *testing.T) {
	cli := client.DefaultClient
	client.DefaultClient = &fakeClient{}
	defer func() {
		client.DefaultClient = cli
	}()

	// Method 1, rawbinary.Do method
	reqbody := []byte("helloworld")
	_, err := Do(context.Background(), reqbody, client.WithDisableServiceRouter())
	assert.Nil(t, err)

	// Method 2, proxy.Do, and specify serviceName
	ctx := trpc.BackgroundContext()
	proxy := NewClientProxyWithName("trpc.app.server.service", "/rawbinary/interface1")
	rspbody, err := proxy.Do(ctx, reqbody)
	assert.Equal(t, reqbody, rspbody)
	assert.Nil(t, err)
}
