// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// Package main gRPC协议的tRPC服务
package main

import (
	"context"
	"os"
	"path/filepath"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/metadata"
	"trpc.group/trpc-go/trpc-codec/grpc"
	tgrpc "trpc.group/trpc-go/trpc-codec/grpc"
	"trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	pb "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/tgrpc"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

func main() {
	// generate TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	// register TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	trpc.ServerConfigPath = cfgPath()
	s := trpc.NewServer(server.WithStreamTransport(grpc.DefaultServerStreamTransport))

	// register service
	pb.RegisterGreeterService(s, &Greeter{})

	// starting all services
	if err := s.Serve(); err != nil {
		panic(err)
	}
}

// cfgPath get the config path
func cfgPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := filepath.Base(pwd)

	// get cfg by dir
	switch dir {
	case "tgrpc":
		return "cfg.yaml"
	case "servers":
		return "tgrpc/cfg.yaml"
	case "examples":
		return "servers/tgrpc/cfg.yaml"
	case "grpc":
		return "examples/servers/tgrpc/cfg.yaml"
	default:
		panic("unknown running dir " + dir)
	}
}

// Greeter struct
type Greeter struct{}

// Hello test
func (*Greeter) Hello(ctx context.Context, req *common.HelloReq) (rsp *common.HelloRsp, err error) {
	// 获取客户端发送的metadata
	md := tgrpc.ParseGRPCMetadata(ctx)
	// get the incoming metadata in ctx
	md1, _ := metadata.FromIncomingContext(ctx)
	_, sc := otelgrpc.Extract(ctx, &md1)
	log.WithContextFields(ctx, "id", sc.TraceID().String())
	log.Infof("get md: %v\n", md)
	msg := codec.Message(ctx)
	// log the frame head
	log.Info(msg.FrameHead())
	rsp = &common.HelloRsp{Msg: "Welcome " + req.Msg}
	// 设置服务端metadata
	for k, v := range md {
		tgrpc.WithServerGRPCMetadata(ctx, k, append(v, "value_from_server"))
	}
	return rsp, nil
}
