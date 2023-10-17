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

package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"trpc.group/trpc-go/trpc-go/server"
)

func TestRegister(t *testing.T) {
	var testCases = []struct {
		serviceName string
		metadata    string
		methodInfos []RegisterMethodsInfo
		Expected    error
	}{
		{
			serviceName: "trpc.test.hellogrpc.Greeter",
			metadata:    "hellogrpc.proto",
			methodInfos: []RegisterMethodsInfo{
				{
					Method: server.Method{
						Name: "SayHello",
						// Func: pb.GreeterService_SayHello_Handler,
					},
				},
			},
			Expected: nil,
		},
	}
	for _, item := range testCases {
		err := Register(item.serviceName, item.metadata, item.methodInfos)
		assert.Equal(t, item.Expected, err)
	}
}

func TestRegisterStream(t *testing.T) {
	type args struct {
		serviceName string
		metadata    string
		streamInfos []server.StreamDesc
		svr         interface{}
		handlerType interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test casr",
			args: args{
				serviceName: "serviceName",
				metadata:    "metadata.proto",
				streamInfos: []server.StreamDesc{},
				svr:         nil,
				handlerType: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterStream(tt.args.serviceName, tt.args.metadata, tt.args.streamInfos, tt.args.svr, tt.args.handlerType); (err != nil) != tt.wantErr {
				t.Errorf("RegisterStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
