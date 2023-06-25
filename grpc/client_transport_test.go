// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	grpcpb "trpc.group/trpc-go/trpc-codec/grpc/testdata/protocols/grpc"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/transport"
)

// FakeGreeterServer pile service
type FakeGreeterServer struct {
	handler func(ctx context.Context, req *common.HelloReq, rsp *common.HelloRsp) error
	grpcpb.UnimplementedGreeterServer
}

// Hello interface
func (f *FakeGreeterServer) Hello(ctx context.Context, req *common.HelloReq) (*common.HelloRsp, error) {
	rsp := &common.HelloRsp{}
	err := f.handler(ctx, req, rsp)
	return rsp, err
}

func _Greeter_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(common.HelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(grpcpb.GreeterServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/helloworld.Greeter/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(grpcpb.GreeterServer).Hello(ctx, req.(*common.HelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var greeterServiceDesc = &grpc.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*grpcpb.GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _Greeter_Hello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/helloworld/helloworld/helloworld.proto",
}

func TestClientTransport(t *testing.T) {
	c := &clientTransport{
		connectionPool: pool{
			connections: sync.Map{},
		},
	}

	// no header
	ctx := context.TODO()
	_, err := c.RoundTrip(ctx, nil)
	assert.NotNil(t, err)

	header := &Header{
		Req: &common.HelloReq{Msg: "req"},
		Rsp: &common.HelloRsp{},
	}
	ctx = context.WithValue(ctx, ContextKeyHeader, header)
	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRPCName("/helloworld.Greeter/Hello")
	msg.WithCalleeServiceName(greeterServiceDesc.ServiceName)
	msg.WithCalleeApp("app")
	msg.WithCalleeServer("server")
	msg.WithCalleeService("Greeter")
	msg.WithCalleeMethod(greeterServiceDesc.Methods[0].MethodName)
	msg.WithSerializationType(codec.SerializationTypePB)

	// cannot access to address
	_, err = c.RoundTrip(ctx, nil,
		transport.WithDialAddress("12345"))
	assert.NotNil(t, err)

	// start grpc listen on available port
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Nil(t, err)
	s := grpc.NewServer()
	fakeGreeterServer := &FakeGreeterServer{}
	s.RegisterService(greeterServiceDesc, fakeGreeterServer)
	go func() {
		err = s.Serve(lis)
		assert.Nil(t, err)
	}()

	// timeout
	fakeGreeterServer.handler = func(ctx context.Context, req *common.HelloReq, rsp *common.HelloRsp) error {
		time.Sleep(3 * time.Second)
		return nil
	}

	msg.WithRequestTimeout(1 * time.Second)
	_, err = c.RoundTrip(ctx, nil,
		transport.WithDialAddress("12345"))
	assert.NotNil(t, err)

	// normal rsp
	fakeGreeterServer.handler = func(ctx context.Context, req *common.HelloReq, rsp *common.HelloRsp) error {
		rsp.Msg = "server copy " + req.Msg
		return nil
	}

	addr := lis.Addr().String()
	_, err = c.RoundTrip(ctx, nil,
		transport.WithDialAddress(addr))
	assert.Nil(t, err)
	rsp := header.Rsp.(*common.HelloRsp)
	assert.Equal(t, "server copy req", rsp.Msg)

	// with metadata
	mdValue := []string{"1"}
	mdByte, err := json.Marshal(mdValue)
	assert.Nil(t, err)
	msg.WithClientMetaData(map[string][]byte{"1": mdByte})
	_, err = c.RoundTrip(ctx, nil,
		transport.WithDialAddress(addr))
	assert.Nil(t, err)
	rsp = header.Rsp.(*common.HelloRsp)
	assert.Equal(t, "server copy req", rsp.Msg)
}

func TestSetGRPCMetadata(t *testing.T) {
	ctx := context.TODO()
	var err error
	ctx, err = setGRPCMetadata(ctx, nil)
	assert.NotNil(t, err)
}
