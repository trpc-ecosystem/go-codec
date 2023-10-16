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
	"reflect"

	"trpc.group/trpc-go/trpc-go/server"
)

// RegisterInfo grpc Information required for registration
type RegisterInfo struct {
	Metadata    string
	ServerFunc  interface{}
	HandlerType interface{}
	MethodsInfo map[string]RegisterMethodsInfo
	StreamsInfo map[string]server.StreamDesc
}

// RegisterMethodsInfo Register the content of the method
type RegisterMethodsInfo struct {
	Method  server.Method
	ReqType reflect.Type
	RspType reflect.Type
}

// RegisterStreamsInfo Register the content of the stream
type RegisterStreamsInfo struct {
	server.StreamDesc
}

var (
	// grpcRegisterInfo: Record the registered information
	grpcRegisterInfo = make(map[string]*RegisterInfo)
)

// Register All external routes used to statically register grpc service, and the mapping of the return type
func Register(serviceName string, metadata string, methodInfos []RegisterMethodsInfo) error {
	registerInfo, ok := grpcRegisterInfo[serviceName]
	if !ok {
		registerInfo = &RegisterInfo{
			Metadata:    metadata,
			MethodsInfo: make(map[string]RegisterMethodsInfo),
			StreamsInfo: make(map[string]server.StreamDesc),
		}
		grpcRegisterInfo[serviceName] = registerInfo
	}

	for _, methodInfo := range methodInfos {
		if registerInfo.MethodsInfo == nil {
			registerInfo.MethodsInfo = make(map[string]RegisterMethodsInfo)
		}
		registerInfo.MethodsInfo[methodInfo.Method.Name] = methodInfo
	}

	return nil
}

// RegisterStream Register grpc stream description information
// Keep the previous RegisterMethod method, which can be collected later
func RegisterStream(serviceName, metadata string, streamInfos []server.StreamDesc,
	svr interface{}, handlerType interface{}) error {
	registerInfo, ok := grpcRegisterInfo[serviceName]
	if !ok {
		registerInfo = &RegisterInfo{
			Metadata:    metadata,
			MethodsInfo: make(map[string]RegisterMethodsInfo),
			StreamsInfo: make(map[string]server.StreamDesc),
		}
		grpcRegisterInfo[serviceName] = registerInfo
	}
	registerInfo.ServerFunc = svr
	registerInfo.HandlerType = handlerType
	for _, streamInfo := range streamInfos {
		if registerInfo.StreamsInfo == nil {
			registerInfo.StreamsInfo = make(map[string]server.StreamDesc)
		}
		registerInfo.StreamsInfo[streamInfo.StreamName] = streamInfo
	}
	return nil
}
