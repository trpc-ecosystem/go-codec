// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

// Package main 发起 gRPC 协议的 tRPC 请求
package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	pb "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/grpc"
)

func main() {
	// create a client connection
	conn, err := grpc.Dial("127.0.0.1:5051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// creates an Exporter
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	// generate TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	// register TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := otel.GetTracerProvider().Tracer("grpc-client")
	ctx, span := tracer.Start(context.Background(), "hello")
	defer span.End()

	// injects correlation context and span context into the gRPC
	md := metadata.MD{}
	otelgrpc.Inject(ctx, &md)

	// creates a new context with outgoing md attached
	ctx = metadata.NewOutgoingContext(ctx, md)

	c := pb.NewGreeterClient(conn)
	// send rpc request
	if rsp, err := c.Hello(ctx, &common.HelloReq{Msg: "abc"}); err != nil {
		panic(err)
	} else {
		fmt.Print(rsp)
	}

}
