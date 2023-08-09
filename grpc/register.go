// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"reflect"

	"trpc.group/trpc-go/trpc-go/server"
)

// RegisterInfo grpc 注册时需要的信息
type RegisterInfo struct {
	Metadata    string
	ServerFunc  interface{}
	HandlerType interface{}
	MethodsInfo map[string]RegisterMethodsInfo
	StreamsInfo map[string]server.StreamDesc
}

// RegisterMethodsInfo 注册 method 的内容
type RegisterMethodsInfo struct {
	Method  server.Method
	ReqType reflect.Type
	RspType reflect.Type
}

// RegisterStreamsInfo 注册 stream 的内容
type RegisterStreamsInfo struct {
	server.StreamDesc
}

var (
	// grpcRegisterInfo： 记录注册的信息
	grpcRegisterInfo = make(map[string]*RegisterInfo)
)

// Register 用于静态注册grpc service 的所有对外路由，以及返回类型的映射
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

// RegisterStream 注册grpc stream描述信息
// 保留前面RegisterMethod方法,后续可以收归到一起
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
