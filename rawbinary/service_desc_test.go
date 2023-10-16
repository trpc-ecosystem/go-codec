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

package rawbinary

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "trpc.group/trpc-go/trpc-go"
    "trpc.group/trpc-go/trpc-go/filter"
    "trpc.group/trpc-go/trpc-go/server"
)

func TestHandler(t *testing.T) {
    ctx := trpc.BackgroundContext()
    _, err := Handler(&mockServer{}, ctx, func(reqbody interface{}) (filter.ServerChain, error) {
        return nil, nil
    })
    assert.Nil(t, err)
}

type mockServer struct{}

// ServiceName is service name.
func (s *mockServer) ServiceName() string {
    return "trpc.rawbinary.mock.mock"
}

// Handle is the processing method.
func (s *mockServer) Handle(ctx context.Context, req []byte) ([]byte, error) {

    rsp := make([]byte, len(req))
    copy(rsp, req)
    return rsp, nil
}

func TestRegister(t *testing.T) {
    svc := server.New(
        server.WithProtocol("rawbinary"),
    )
    s := &server.Server{}
    s.AddService("trpc.rawbinary.mock.mock", svc)
    err := Register(s, &mockServer{})
    assert.Nil(t, err)
}
