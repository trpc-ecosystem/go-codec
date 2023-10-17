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

package rawstring

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/server"
)

type mockRawStringServer struct{}

func (s *mockRawStringServer) ServiceName() string {
	return "trpc.rawstring.mock.Mock"
}

func (s *mockRawStringServer) Handle(ctx context.Context, req []byte) ([]byte, error) {
	rsp := make([]byte, len(req))
	copy(rsp, req)
	return rsp, nil
}

func TestHandler(t *testing.T) {
	_, err := Handler(&mockRawStringServer{}, trpc.BackgroundContext(), func(_ interface{}) (filter.ServerChain, error) {
		return nil, nil
	})

	assert.Nil(t, err)
}

func TestRegister(t *testing.T) {
	svc := server.New(server.WithProtocol("rawstring"))
	s := &server.Server{}
	s.AddService("trpc.rawstring.mock.Mock", svc)
	err := Register(s, &mockRawStringServer{})
	assert.Nil(t, err)
}
