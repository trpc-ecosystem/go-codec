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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/naming/selector"
)

type MockGrpcClientStream struct {
}

func (cs *MockGrpcClientStream) Header() (metadata.MD, error) {
	return nil, nil
}
func (cs *MockGrpcClientStream) Trailer() metadata.MD {
	return nil
}
func (cs *MockGrpcClientStream) CloseSend() error {
	return nil
}
func (cs *MockGrpcClientStream) Context() context.Context {
	return nil
}
func (cs *MockGrpcClientStream) SendMsg(m interface{}) error {
	return nil
}

func (cs *MockGrpcClientStream) RecvMsg(m interface{}) error {
	return nil
}

func TestNewStreamClient(t *testing.T) {
	ctx := context.Background()
	s := NewStreamClient()

	selector.DefaultSelector = &MockSelector{}
	_, err := s.NewStream(ctx, &client.ClientStreamDesc{}, "method")
	assert.NotNil(t, err)
}

func Test_clientStream_RecvMsg(t *testing.T) {
	type fields struct {
		ctx    context.Context
		stream grpc.ClientStream
	}
	type args struct {
		m interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test case",
			fields: fields{
				ctx:    context.Background(),
				stream: &MockGrpcClientStream{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &clientStream{
				ctx:    tt.fields.ctx,
				stream: tt.fields.stream,
			}
			if err := cs.RecvMsg(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("clientStream.RecvMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientStream_SendMsg(t *testing.T) {
	type fields struct {
		ctx    context.Context
		stream grpc.ClientStream
	}
	type args struct {
		m interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test case",
			fields: fields{
				ctx:    context.Background(),
				stream: &MockGrpcClientStream{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &clientStream{
				ctx:    tt.fields.ctx,
				stream: tt.fields.stream,
			}
			if err := cs.SendMsg(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("clientStream.SendMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientStream_CloseSend(t *testing.T) {
	type fields struct {
		ctx    context.Context
		stream grpc.ClientStream
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test case",
			fields: fields{
				ctx:    context.Background(),
				stream: &MockGrpcClientStream{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &clientStream{
				ctx:    tt.fields.ctx,
				stream: tt.fields.stream,
			}
			if err := cs.CloseSend(); (err != nil) != tt.wantErr {
				t.Errorf("clientStream.CloseSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clientStream_Context(t *testing.T) {
	type fields struct {
		ctx    context.Context
		stream grpc.ClientStream
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "test case",
			fields: fields{
				ctx:    context.Background(),
				stream: &MockGrpcClientStream{},
			},
			want: context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &clientStream{
				ctx:    tt.fields.ctx,
				stream: tt.fields.stream,
			}
			if got := cs.Context(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientStream.Context() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeGrpcDesc(t *testing.T) {
	type args struct {
		desc *client.ClientStreamDesc
	}
	tests := []struct {
		name string
		args args
		want *grpc.StreamDesc
	}{
		{
			name: "test case",
			args: args{
				desc: &client.ClientStreamDesc{
					StreamName:    "StreamName",
					ClientStreams: true,
				},
			},
			want: &grpc.StreamDesc{
				StreamName:    "StreamName",
				ClientStreams: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeGrpcDesc(tt.args.desc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeGrpcDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}
