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

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockGrpcServerStream struct{}

func (m *mockGrpcServerStream) SetHeader(metadata.MD) error {
	return nil
}

func (m *mockGrpcServerStream) SendHeader(metadata.MD) error {
	return nil
}

func (m *mockGrpcServerStream) SetTrailer(metadata.MD) {
}

func (m *mockGrpcServerStream) Context() context.Context {
	return context.Background()
}

func (m *mockGrpcServerStream) SendMsg(i interface{}) error {
	return nil
}

func (m *mockGrpcServerStream) RecvMsg(i interface{}) error {
	return nil
}

func TestStreamHandler(t *testing.T) {
	type args struct {
		srv interface{}
		s   grpc.ServerStream
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				srv: "data",
				s:   &mockGrpcServerStream{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StreamHandler(tt.args.srv, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("StreamHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
