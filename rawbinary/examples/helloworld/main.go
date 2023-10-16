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

// Package main is the main package.
package main

import (
	"context"
	"log"
	"time"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"

	"trpc.group/trpc-go/trpc-codec/rawbinary"
)

func main() {
	go func() {
		s := trpc.NewServer()
		rawbinary.Register(s, &helloworldServer2{})
		rawbinary.Register(s, &helloworldServer{})
		s.Serve()
	}()

	time.Sleep(time.Second)

	rsp, err := rawbinary.Do(context.Background(), []byte("helloworld"),
		client.WithTarget("ip://127.0.0.1:8001"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("recv: ", string(rsp))

	rsp, err = rawbinary.Do(context.Background(), []byte("helloworld2"),
		client.WithTarget("ip://127.0.0.1:8002"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("recv: ", string(rsp))

	// Call through proxy and specify calleeServiceName.
	ctx := trpc.BackgroundContext()
	proxy := rawbinary.NewClientProxyWithName("trpc.app.server.service", "/rawbinary/interface")
	rsp, err = proxy.Do(ctx, []byte("hello proxy do"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("third recv: ", string(rsp))
}

// helloworldServer is a struct for test server.
type helloworldServer struct{}

// ServiceName service name
func (s *helloworldServer) ServiceName() string {
	return "trpc.rawbinary.helloworld.Helloworld"
}

// Handle processes entry.
func (s *helloworldServer) Handle(ctx context.Context, req []byte) ([]byte, error) {

	rsp := make([]byte, len(req))
	copy(rsp, req)
	log.Println("helloworld")
	return rsp, nil
}

// helloworldServer2 is a struct for test server.
type helloworldServer2 struct{}

// ServiceName service name
func (s *helloworldServer2) ServiceName() string {
	return "trpc.rawbinary.helloworld.Helloworld2"
}

// Handle processes entry.
func (s *helloworldServer2) Handle(ctx context.Context, req []byte) ([]byte, error) {

	rsp := make([]byte, len(req))
	copy(rsp, req)
	log.Println("helloworld2")
	return rsp, nil
}
