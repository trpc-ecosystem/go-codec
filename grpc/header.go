// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Header stored in context to communicate with trpc
type Header struct {
	Req         interface{}         // request
	Rsp         interface{}         // response
	InMetadata  map[string][]string // metadata from client
	OutMetadata map[string][]string // metadata send to client
}

// ParseGRPCMetadata trpc-go服务端调用，来获取客户端的metadata
func ParseGRPCMetadata(ctx context.Context) map[string][]string {
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return nil
	}
	return header.InMetadata
}

// WithServerGRPCMetadata trpc-go服务端调用，来发送metadata
func WithServerGRPCMetadata(ctx context.Context, key string, value []string) {
	header, ok := ctx.Value(ContextKeyHeader).(*Header)
	if !ok {
		return
	}
	if header == nil {
		header = &Header{}
	}
	if header.OutMetadata == nil {
		header.OutMetadata = metadata.MD{}
	}
	header.OutMetadata[key] = value
}

// WithHeader trpc-go客户端调用，设置发送到服务端的md，或接受来自服务端的md
func WithHeader(ctx context.Context, header *Header) context.Context {
	return context.WithValue(ctx, ContextKeyHeader, header)
}
