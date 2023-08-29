// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/codec"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/naming/registry"
	"trpc.group/trpc-go/trpc-go/naming/selector"
	"trpc.group/trpc-go/trpc-go/transport"
)

var (
	defaultCallOptionsSize   = 4
	defaultSelectOptionsSize = 4
)

// getOptions 获取每次请求所需的参数数据
func getOptions(msg codec.Msg, opt ...client.Option) (*client.Options, string, error) {
	// 每次请求构造新的参数数据 保证并发安全
	opts := &client.Options{
		Transport:                transport.DefaultClientTransport,
		Selector:                 selector.DefaultSelector,
		CallOptions:              make([]transport.RoundTripOption, 0, defaultCallOptionsSize),
		SelectOptions:            make([]selector.Option, 0, defaultSelectOptionsSize),
		SerializationType:        -1, // 初始值 -1，不设置
		CurrentSerializationType: -1, // 当前 client 的序列化方式，协议里面的序列化方式以 SerializationType 为准，转发代理情况，
		CurrentCompressType:      -1, // 当前 client 透传 body 不序列化，但是业务协议后端需要指定序列化方式
	}
	// 设置服务环境信息
	opts.SelectOptions = append(opts.SelectOptions, getServiceInfoOptions(msg)...)

	// 输入参数为最高优先级 覆盖掉原有数据
	for _, o := range opt {
		o(opts)
	}
	address, err := setNamingInfo(opts)
	if err != nil {
		return nil, "", err
	}
	return opts, address, nil
}

func setNamingInfo(opts *client.Options) (string, error) {
	// 默认使用名字服务 servicename 获取地址，如果有指定 target，则由指定的值来获取
	if opts.Target == "" {
		return "", nil
	}
	// Target 的格式为：selector://endpoint
	substr := "://"
	index := strings.Index(opts.Target, substr)
	if index == -1 {
		return "", errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client: target %s schema invalid", opts.Target))
	}
	opts.Selector = selector.Get(opts.Target[:index])
	// 检查 selector 是否为空
	if opts.Selector == nil {
		return "", errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client: selector %s not exist",
			opts.Target[:index]))
	}
	address := opts.Target[index+len(substr):]
	return address, nil
}

// selectNode 根据设置的寻址选择器寻址到后端节点，并设置 msg
func selectNode(msg codec.Msg, opts *client.Options, address string) (*registry.Node, error) {
	node, err := getNode(address, opts)
	if err != nil {
		return nil, err
	}

	// 通过注册中心返回的节点配置信息更新设置参数
	opts.LoadNodeConfig(node)
	msg.WithCalleeContainerName(node.ContainerName)
	msg.WithCalleeSetName(node.SetName)

	if len(msg.EnvTransfer()) > 0 {
		// 优先使用上游的透传环境信息
		msg.WithEnvTransfer(msg.EnvTransfer())
	} else {
		// 上游没有透传则使用本环境信息
		msg.WithEnvTransfer(node.EnvKey)
	}

	// 禁用服务路由则清空环境信息
	if opts.DisableServiceRouter {
		if len(msg.EnvTransfer()) > 0 {
			msg.WithEnvTransfer("")
		}
	}
	return node, nil
}

func getNode(address string, opts *client.Options) (*registry.Node, error) {
	// 获取 ipport 请求地址
	node, err := opts.Selector.Select(address, opts.SelectOptions...)
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientRouteErr, "client Select: "+err.Error())
	}
	if node.Address == "" {
		return nil, errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client Select: node address empty:%+v", node))
	}
	return node, nil
}

// updateMsg 更新客户端请求 Msg 上下文信息
func updateMsg(msg codec.Msg, opts *client.Options) {
	// 设置被调方 service name 一般 service name 和 proto 协议的 package.service 一致，但是用户可以通过参数修改
	if len(opts.ServiceName) > 0 {
		msg.WithCalleeServiceName(opts.ServiceName) // 以 client 角度看，caller 是自身，callee 是下游
	}

	if len(opts.CalleeMethod) > 0 {
		msg.WithCalleeMethod(opts.CalleeMethod)
	}

	// 设置后端透传参数
	msg.WithClientMetaData(getMetaData(msg, opts))

	// 以 client 作为小工具时，没有 caller，需要自己通过 client option 设置进来
	if len(opts.CallerServiceName) > 0 {
		msg.WithCallerServiceName(opts.CallerServiceName)
	}
	if opts.SerializationType >= 0 {
		msg.WithSerializationType(opts.SerializationType)
	}
	if opts.CompressType > 0 {
		msg.WithCompressType(opts.CompressType)
	}

	// 用户设置 reqhead，希望使用自己的请求包头
	if opts.ReqHead != nil {
		msg.WithClientReqHead(opts.ReqHead)
	}
	// 用户设置 rsphead，希望回传后端的响应包头
	if opts.RspHead != nil {
		msg.WithClientRspHead(opts.RspHead)
	}

	msg.WithCallType(opts.CallType)
}

// getServiceInfoOptions 设置服务环境信息
func getServiceInfoOptions(msg codec.Msg) []selector.Option {
	if len(msg.Namespace()) > 0 {
		return []selector.Option{
			selector.WithSourceNamespace(msg.Namespace()),
			selector.WithSourceServiceName(msg.CallerServiceName()),
			selector.WithSourceEnvName(msg.EnvName()),
			selector.WithEnvTransfer(msg.EnvTransfer()),
			selector.WithSourceSetName(msg.SetName()),
		}
	}
	return nil
}

// getMetaData 获取后端透传参数
func getMetaData(msg codec.Msg, opts *client.Options) codec.MetaData {
	md := msg.ClientMetaData()
	if md == nil {
		md = codec.MetaData{}
	}
	for k, v := range opts.MetaData {
		md[k] = v
	}
	return md
}
