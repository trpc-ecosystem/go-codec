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

// Package main initiates tRPC streaming requests using the gRPC protocol.
package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/metadata"
	common "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	pb "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/streams"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

func init() {
	// register trace provider
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
}

func main() {
	opts := []client.Option{
		client.WithTarget("ip://127.0.0.1:5051"),
		client.WithTimeout(1 * time.Second),
	}
	clientProxy := pb.NewGreeterClientProxy(opts...)
	// test Hello
	Hello(clientProxy)

	// test GetStream
	GetStream(clientProxy)

	// test PutStream
	PutStream(clientProxy)

	// test AllStream
	AllStream(clientProxy)
}

// Hello send msg
func Hello(c pb.GreeterClientProxy) {
	req := &common.HelloReq{
		Msg: "trpc_client_Hello",
	}
	tracer := otel.GetTracerProvider().Tracer("grpc-client")
	ctx, span := tracer.Start(trpc.BackgroundContext(), "hello")
	defer span.End()

	// injects correlation context and span context into the gRPC
	md := metadata.MD{}
	otelgrpc.Inject(ctx, &md)

	// creates a new context with outgoing md attached
	ctx = metadata.NewOutgoingContext(ctx, md)
	rsp, err := c.Hello(ctx, req)
	if err != nil {
		log.Error(err)
	}
	log.Info(rsp)
}

// GetStream get stream msg
func GetStream(c pb.GreeterClientProxy) {
	req := &common.HelloReq{
		Msg: "trpc_client_Hello",
	}
	// GetStream
	stream, err := c.GetStream(context.Background(), req)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		rsp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
			return
		}
		log.Info(rsp)
	}

}

// PutStream put stream msg
func PutStream(c pb.GreeterClientProxy) {
	stream, err := c.PutStream(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	// send stream 10 times
	for i := 1; i <= 10; i++ {
		if err := stream.Send(&common.HelloReq{Msg: fmt.Sprintf("req: %d", i)}); err != nil {
			log.Error(err)
		}
		time.Sleep(1 * time.Second)
	}
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		log.Error(err)
	}
	log.Info(rsp)

}

// AllStream all
func AllStream(c pb.GreeterClientProxy) {
	stream, err := c.AllStream(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Info("EOF")
					return
				}
				log.Error(err)
				return
			}
			log.Info(data)
		}
	}()
	for i := 1; i <= 10; i++ {
		req := &common.HelloReq{
			Msg: fmt.Sprintf("grpc_client_all_stream_%d", i),
		}
		err := stream.Send(req)
		if err != nil {
			log.Error(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
	stream.CloseSend()
}
