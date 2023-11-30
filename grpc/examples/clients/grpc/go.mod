module trpc.group/trpc-go/trpc-codec/grpc/examples/clients/grpc_unary

go 1.18

replace trpc.group/trpc-go/trpc-codec/grpc => ../../../../grpc

require (
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.29.0
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.4.1
	go.opentelemetry.io/otel/sdk v1.4.1
	google.golang.org/grpc v1.44.0
	trpc.group/trpc-go/trpc-codec/grpc v0.0.0-20220216035248-7863fc4506d7
)

require (
	github.com/go-logr/logr v1.2.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.opentelemetry.io/otel/trace v1.4.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20211019152133-63b7e35f4404 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
