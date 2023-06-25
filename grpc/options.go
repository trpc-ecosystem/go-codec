// Tencent is pleased to support the open source community by making tRPC available.
// Copyright (C) 2023 THL A29 Limited, a Tencent company. All rights reserved.
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.

package grpc

import (
	"fmt"
	"strings"
	"time"

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

// getOptions Get the parameter data required for each request
func getOptions(msg codec.Msg, opt ...client.Option) (*client.Options, string, error) {
	// Each request constructs new parameter data to ensure concurrency safety
	opts := &client.Options{
		Transport:                transport.DefaultClientTransport,
		Selector:                 selector.DefaultSelector,
		CallOptions:              make([]transport.RoundTripOption, 0, defaultCallOptionsSize),
		SelectOptions:            make([]selector.Option, 0, defaultSelectOptionsSize),
		SerializationType:        -1, // Initial value -1, do not set
		CurrentSerializationType: -1, // The serialization method of the current client, the serialization method in
		// the protocol is based on the SerializationType, the forwarding proxy situation.
		CurrentCompressType: -1, // The current client transparent transmission body is not serialized, but the
		// backend of the business agreement needs to specify the serialization method.
	}

	// Use the servicename (package.service) of the protocol file of the transferred party as the key to obtain the
	//   relevant configuration.
	if err := loadClientConfig(opts, msg.CalleeServiceName()); err != nil {
		return nil, "", err
	}

	// Set service environment information
	opts.SelectOptions = append(opts.SelectOptions, getServiceInfoOptions(msg)...)

	// The input parameter is the highest priority to overwrite the original data
	for _, o := range opt {
		o(opts)
	}
	address, err := setNamingInfo(opts)
	if err != nil {
		return nil, "", err
	}
	return opts, address, nil
}

func loadClientConfig(opts *client.Options, key string) error {
	cfg := client.Config(key)
	if cfg.Timeout > 0 {
		opts.Timeout = time.Duration(cfg.Timeout) * time.Millisecond
	}
	if cfg.Serialization != nil {
		opts.SerializationType = *cfg.Serialization
	}

	if cfg.Network != "" {
		opts.Network = cfg.Network
		opts.CallOptions = append(opts.CallOptions, transport.WithDialNetwork(cfg.Network))
	}
	if cfg.Password != "" {
		opts.CallOptions = append(opts.CallOptions, transport.WithDialPassword(cfg.Password))
	}
	if cfg.CACert != "" {
		opts.CallOptions = append(opts.CallOptions,
			transport.WithDialTLS(cfg.TLSCert, cfg.TLSKey, cfg.CACert, cfg.TLSServerName))
	}
	return nil
}

func setNamingInfo(opts *client.Options) (string, error) {
	// By default, the name service servicename is used to obtain the address. If there is a specified target,
	//   the specified value will be used to obtain the address.
	if opts.Target == "" {
		return "", nil
	}
	// The format of Target is: selector://endpoint
	substr := "://"
	index := strings.Index(opts.Target, substr)
	if index == -1 {
		return "", errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client: target %s schema invalid", opts.Target))
	}
	opts.Selector = selector.Get(opts.Target[:index])
	// Check if selector is empty
	if opts.Selector == nil {
		return "", errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client: selector %s not exist",
			opts.Target[:index]))
	}
	address := opts.Target[index+len(substr):]
	return address, nil
}

// selectNode Address to the backend node according to the address selector set, and set msg
func selectNode(msg codec.Msg, opts *client.Options, address string) (*registry.Node, error) {
	node, err := getNode(address, opts)
	if err != nil {
		return nil, err
	}

	// Update the setting parameters through the node configuration information returned by the registration center.
	opts.LoadNodeConfig(node)
	msg.WithCalleeContainerName(node.ContainerName)
	msg.WithCalleeSetName(node.SetName)

	if len(msg.EnvTransfer()) > 0 {
		// Prioritize the use of upstream transparent transmission environment information
		msg.WithEnvTransfer(msg.EnvTransfer())
	} else {
		// If there is no transparent transmission upstream, use this environment information
		msg.WithEnvTransfer(node.EnvKey)
	}

	// Disable the service route to clear the environment information
	if opts.DisableServiceRouter {
		if len(msg.EnvTransfer()) > 0 {
			msg.WithEnvTransfer("")
		}
	}
	return node, nil
}

func getNode(address string, opts *client.Options) (*registry.Node, error) {
	// Get ipport request address
	node, err := opts.Selector.Select(address, opts.SelectOptions...)
	if err != nil {
		return nil, errs.NewFrameError(errs.RetClientRouteErr, "client Select: "+err.Error())
	}
	if node.Address == "" {
		return nil, errs.NewFrameError(errs.RetClientRouteErr, fmt.Sprintf("client Select: node address empty:%+v", node))
	}
	return node, nil
}

// updateMsg Update client request Msg context information
func updateMsg(msg codec.Msg, opts *client.Options) {
	// Set the service name of the called party. Generally, the service name is consistent with the package.service
	//   of the proto protocol, but the user can modify it through parameters.
	if len(opts.ServiceName) > 0 {
		msg.WithCalleeServiceName(opts.ServiceName) // From the perspective of the client, the caller is itself,
		//   and the callee is the downstream.
	}

	if len(opts.CalleeMethod) > 0 {
		msg.WithCalleeMethod(opts.CalleeMethod)
	}

	// Set backend transparent transmission parameters
	msg.WithClientMetaData(getMetaData(msg, opts))

	// When the client is used as a small tool, there is no caller, and you need to set it through the client option.
	if len(opts.CallerServiceName) > 0 {
		msg.WithCallerServiceName(opts.CallerServiceName)
	}
	if opts.SerializationType >= 0 {
		msg.WithSerializationType(opts.SerializationType)
	}
	if opts.CompressType > 0 {
		msg.WithCompressType(opts.CompressType)
	}

	// The user sets reqhead and wants to use his own request header
	if opts.ReqHead != nil {
		msg.WithClientReqHead(opts.ReqHead)
	}
	// The user sets rsphead and hopes to return the response header of the backend
	if opts.RspHead != nil {
		msg.WithClientRspHead(opts.RspHead)
	}

	msg.WithCallType(opts.CallType)
}

// getServiceInfoOptions Set service environment information
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

// getMetaData Obtain backend transparent transmission parameters
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
