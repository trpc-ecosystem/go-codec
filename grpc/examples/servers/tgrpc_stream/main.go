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

// Package main gRPC 协议的 tRPC 流式服务
package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/metadata"
	"trpc.group/trpc-go/trpc-codec/grpc"
	common "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	pb "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/streams"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

// ServiceImpl impl
type ServiceImpl struct {
}

func main() {
	// generate TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// register TracerProvider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	trpc.ServerConfigPath = "cfg.yaml"
	s := trpc.NewServer(server.WithStreamTransport(grpc.DefaultServerStreamTransport))

	// register service
	pb.RegisterGreeterService(s, &ServiceImpl{})

	// starting all services
	if err := s.Serve(); err != nil {
		panic(err)
	}
}

// Hello test
func (s *ServiceImpl) Hello(ctx context.Context, req *common.HelloReq) (rsp *common.HelloRsp, err error) {
	// get the incoming metadata in ctx
	md, _ := metadata.FromIncomingContext(ctx)
	_, sc := otelgrpc.Extract(ctx, &md)
	log.WithContextFields(ctx, "id", sc.TraceID().String())
	log.Infof("Hello, req: %+v", req)
	rsp = &common.HelloRsp{Msg: "OK"}
	return rsp, nil
}

// GetStream GetStreamMsg
func (s *ServiceImpl) GetStream(req *common.HelloReq, stream pb.Greeter_GetStreamServer) error {
	log.Infof("GetStream, req: %+v", req)
	for i := 1; i <= 10; i++ {
		rsp := &common.HelloRsp{Msg: fmt.Sprintf("GetStream_Response_%d", i)}
		// send response
		err := stream.Send(rsp)
		if err != nil {
			log.Error(err)
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// PutStream test
func (s *ServiceImpl) PutStream(stream pb.Greeter_PutStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// close stream
				return stream.SendAndClose(&common.HelloRsp{Msg: "END"})
			}
			log.Error(err)
			return err
		}
		log.Infof("PutStream, req:%+v", req)
	}
}

// AllStream test
func (s *ServiceImpl) AllStream(stream pb.Greeter_AllStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
			return err
		}
		log.Infof("AllStream, req:%+v", req)
		// send stream
		err = stream.Send(&common.HelloRsp{Msg: "got" + req.Msg})
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}
