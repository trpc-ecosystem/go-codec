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
    "fmt"
    "log"
    "sync"
    "time"

    "trpc.group/trpc-go/trpc-codec/rawstring"
    "trpc.group/trpc-go/trpc-go"
    "trpc.group/trpc-go/trpc-go/client"
)

var (
    max = 5
)

func main() {
    go func() {
        s := trpc.NewServer()
        rawstring.Register(s, &helloworldServer2{})
        rawstring.Register(s, &helloworldServer{})
        s.Serve()

    }()
    time.Sleep(1 * time.Second)
    wg := sync.WaitGroup{}

    for i := 0; i < max; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            req := fmt.Sprintf("uin=286989429&i=%v\n", i)
            ctx := context.Background()
            rsp, err := rawstring.Do(ctx, req,
                client.WithTarget("ip://127.0.0.1:8000\b"))
            if err != nil {
                fmt.Printf("err %v\n", err)
            }
            fmt.Printf("rsp %v", rsp)

        }(i)
    }

    for i := 0; i < max; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            req := fmt.Sprintf("test2uin=286989429&i=%v", i)
            ctx := context.Background()
            rsp, err := rawstring.Do(ctx, req,
                client.WithTarget("ip://127.0.0.1:8001"))
            if err != nil {

                fmt.Printf("err %v\n", err)
                log.Fatalln(err)
            }
            fmt.Printf("rsp %v", rsp)
        }(i)
    }
    wg.Wait()
    log.Println("finish succ")
}

// helloworldServer 测试 server
type helloworldServer struct{}

// ServiceName service name
func (s *helloworldServer) ServiceName() string {
    return "trpc.rawstring.helloworld.hellorawstring"
}

// Handle 处理入口
func (s *helloworldServer) Handle(ctx context.Context, req []byte) ([]byte, error) {
    rsp := make([]byte, len(req))
    copy(rsp, req)
    return rsp, nil
}

// helloworldServer2 测试 server
type helloworldServer2 struct{}

// ServiceName service name
func (s *helloworldServer2) ServiceName() string {
    return "trpc.rawstring.helloworld.hellorawstring1"
}

// Handle 处理入口
func (s *helloworldServer2) Handle(ctx context.Context, req []byte) ([]byte, error) {
    rsp := make([]byte, len(req))
    copy(rsp, req)
    return rsp, nil
}
