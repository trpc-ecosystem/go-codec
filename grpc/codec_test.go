// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/transport"
)

type h struct{}

// Handle Handle
func (*h) Handle(ctx context.Context, reqbuf []byte) (rsp []byte, err error) {
	fmt.Println("recv grpc req")
	return nil, nil
}

type hErr struct{}

// Handle Handle
func (*hErr) Handle(ctx context.Context, reqbuf []byte) (rsp []byte, err error) {
	return nil, errors.New("handle error")
}

func TestListenAndServe(t *testing.T) {
	grpcRegisterInfo["serviceName"] = &RegisterInfo{
		MethodsInfo: map[string]RegisterMethodsInfo{
			"method": {},
		},
	}
	opt := transport.WithListenAddress(":8080")
	err := DefaultServerTransport.ListenAndServe(context.Background(), opt)
	assert.Equal(t, err, errors.New("trpc server transport handler empty"))

	a := transport.WithHandler(transport.Handler(&h{}))
	err = DefaultServerTransport.ListenAndServe(context.Background(), opt, a)
	assert.Nil(t, err)

	opt = transport.WithListenAddress("sss")
	a = transport.WithHandler(transport.Handler(&h{}))
	err = DefaultServerTransport.ListenAndServe(context.Background(), opt, a)
	tmp := &net.AddrError{Err: "missing port in address", Addr: "sss"}
	assert.Equal(t, err.(*net.OpError).Unwrap().Error(), tmp.Error())
}

func TestServerCodecDecode(t *testing.T) {
	sc := ServerCodec{}
	msg := codec.Message(context.Background())
	srcBts := []byte("hello,world")
	dest, err := sc.Decode(msg, srcBts)
	assert.Nil(t, err)
	assert.Equal(t, dest, srcBts)
}

func TestServerCodecEncode(t *testing.T) {
	sc := ServerCodec{}
	msg := codec.Message(context.Background())
	srcBts := []byte("hello,world")
	dest, err := sc.Encode(msg, srcBts)
	assert.Nil(t, err)
	assert.Equal(t, dest, srcBts)
}

func TestGrpcToTrpcLayerHandle(t *testing.T) {
	// empty ctx
	gToTrpcLayer := &GrpcToTrpcLayer{
		Handler: &h{},
	}
	msg := codec.Message(context.Background())
	ctx := msg.Context()
	msg.WithSerializationType(codec.SerializationTypeJSON)
	_, err := gToTrpcLayer.Handle(
		nil,
		ctx,
		func(src interface{}) error {
			return errors.New("parse error.")
		},
		nil,
	)
	assert.Equal(t, err, errors.New("GrpcToTrpcLayer: method: `` format error. "))

	peer.NewContext(ctx, &peer.Peer{})
	metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"key": "value"}))
	grpc.NewContextWithServerTransportStream(ctx, new(STS))

}

// STS STS 桩结构体
type STS struct {
}

// Method Method
func (s *STS) Method() string {
	return "/Service/Method"
}

// SetHeader SetHeader
func (s *STS) SetHeader(md metadata.MD) error {
	return nil
}

// SendHeader SendHeader
func (s *STS) SendHeader(md metadata.MD) error {
	return nil
}

// SetTrailer SetTrailer
func (s *STS) SetTrailer(md metadata.MD) error {
	return nil
}

// Addr 地址
type Addr struct {
}

// String 字符串
func (a *Addr) String() string {
	return "addr"
}

// Network 地址
func (a *Addr) Network() string {
	return ""
}

func TestGrpcToTrpcLayer_Handle(t *testing.T) {
	grpcRegisterInfo["Service"] = &RegisterInfo{
		MethodsInfo: map[string]RegisterMethodsInfo{
			"Method": {
				ReqType: reflect.TypeOf("a"),
			},
		},
	}
	type fields struct {
		Handler transport.Handler
	}
	type args struct {
		srv interface{}
		ctx context.Context
		dec func(interface{}) error
		in3 grpc.UnaryServerInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		args   args

		wantErr bool
	}{
		{
			name: "empty ctx",
			fields: fields{
				Handler: &h{},
			},
			args: args{
				srv: nil,
				ctx: context.Background(),
				dec: func(src interface{}) error {
					return errors.New("parse error.")
				},
				in3: nil,
			},
			wantErr: true,
		},
		{
			name: "more ctx",
			fields: fields{
				Handler: &h{},
			},
			args: args{
				srv: nil,
				ctx: getContext(),
				dec: func(src interface{}) error {
					return nil
				},
				in3: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcToTrpcLayer{
				Handler: tt.fields.Handler,
			}
			_, err := g.Handle(tt.args.srv, tt.args.ctx, tt.args.dec, tt.args.in3)
			if (err != nil) != tt.wantErr {
				t.Errorf("GrpcToTrpcLayer.Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func getContext() context.Context {
	msg := codec.Message(context.Background())
	ctx := msg.Context()
	ctx = peer.NewContext(ctx, &peer.Peer{Addr: new(Addr)})
	ctx = metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"key": "value"}))
	ctx = grpc.NewContextWithServerTransportStream(ctx, new(STS))
	return ctx
}

func TestClientCodecd(t *testing.T) {
	ctx := context.Background()
	msg := codec.Message(ctx)
	body := []byte("test")
	CCodec := &ClientCodec{}
	rsp, err := CCodec.Decode(msg, body)
	assert.Equal(t, body, rsp)
	assert.Nil(t, err)

	buf, err := CCodec.Encode(msg, body)
	assert.Equal(t, body, buf)
	assert.Nil(t, err)
}
