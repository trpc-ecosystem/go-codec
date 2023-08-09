// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// Package trpc 实现接口
package trpc

import (
	"context"
	"log"

	tgrpc "trpc.group/trpc-go/trpc-codec/grpc"
	"trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
)

// Greeter struct
type Greeter struct{}

// Hello 实现hello接口
func (*Greeter) Hello(ctx context.Context, req *common.HelloReq, rsp *common.HelloRsp) error {
	// 获取客户端发送的metadata
	md := tgrpc.ParseGRPCMetadata(ctx)
	log.Printf("get md: %v\n", md)
	rsp.Msg = "Welcome " + req.Msg
	// 设置服务端metadata
	for k, v := range md {
		tgrpc.WithServerGRPCMetadata(ctx, k, append(v, "value_from_server"))
	}
	return nil
}
