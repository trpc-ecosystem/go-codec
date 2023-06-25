// Package main is the main package.
package main

import (
	"context"
	"os"
	"path/filepath"

	"git.code.oa.com/trpc-go/trpc-codec/grpc"
	tgrpc "git.code.oa.com/trpc-go/trpc-codec/grpc"
	"git.code.oa.com/trpc-go/trpc-codec/grpc/testdata/protocols/common"
	pb "git.code.oa.com/trpc-go/trpc-codec/grpc/testdata/protocols/tgrpc"
	"git.code.oa.com/trpc-go/trpc-go"
	"git.code.oa.com/trpc-go/trpc-go/codec"
	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/server"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/metadata"
)

func main() {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	trpc.ServerConfigPath = cfgPath()
	s := trpc.NewServer(server.WithStreamTransport(grpc.DefaultServerStreamTransport))
	pb.RegisterGreeterService(s, &Greeter{})

	if err := s.Serve(); err != nil {
		panic(err)
	}
}

func cfgPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dir := filepath.Base(pwd)

	switch dir {
	case "tgrpc":
		return "cfg.yaml"
	case "servers":
		return "tgrpc/cfg.yaml"
	case "examples":
		return "servers/tgrpc/cfg.yaml"
	case "grpc":
		return "examples/servers/tgrpc/cfg.yaml"
	default:
		panic("unknown running dir " + dir)
	}
}

// Greeter TODO
type Greeter struct{}

// Hello TODO
func (*Greeter) Hello(ctx context.Context, req *common.HelloReq, rsp *common.HelloRsp) error {
	// Get the metadata sent by the client
	md := tgrpc.ParseGRPCMetadata(ctx)
	md1, _ := metadata.FromIncomingContext(ctx)
	_, sc := otelgrpc.Extract(ctx, &md1)
	l := log.WithFields("id", sc.TraceID().String())
	l.Infof("get md: %v\n", md)
	msg := codec.Message(ctx)
	l.Info(msg.FrameHead())
	rsp.Msg = "Welcome " + req.Msg
	// Set server metadata
	for k, v := range md {
		tgrpc.WithServerGRPCMetadata(ctx, k, append(v, "value_from_server"))
	}
	return nil
}
