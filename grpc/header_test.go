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

package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGRPCMetadata(t *testing.T) {
	// value is not *Header
	ctx := context.WithValue(context.TODO(), ContextKeyHeader, "")
	parseV := ParseGRPCMetadata(ctx)
	assert.Nil(t, parseV)

	// normal case
	inMetadata := map[string][]string{"key": {"1", "2"}}
	v := &Header{
		InMetadata: inMetadata,
	}
	ctx = context.WithValue(context.TODO(), ContextKeyHeader, v)
	parseV = ParseGRPCMetadata(ctx)
	assert.Equal(t, v.InMetadata, parseV)
}

func TestWithServerGRPCMetadata(t *testing.T) {
	k := "key"
	v := []string{"1", "2"}
	h := &Header{}
	ctx := context.WithValue(context.TODO(), ContextKeyHeader, h)
	WithServerGRPCMetadata(ctx, k, v)
	assert.Contains(t, h.OutMetadata, k)
	assert.Equal(t, v, h.OutMetadata[k])
}

func TestWithHeader(t *testing.T) {
	ctx := context.TODO()
	h := ctx.Value(ContextKeyHeader)
	assert.Nil(t, h)
	header := &Header{}
	ctx = WithHeader(ctx, header)
	h = ctx.Value(ContextKeyHeader)
	assert.NotNil(t, h)
}
