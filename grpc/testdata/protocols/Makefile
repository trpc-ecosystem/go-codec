.PHONY: all
all: trpc grpc tgrpc gstreams tstreams

.PHONY: common
common:
	protoc --go_out=common --go_opt=paths=source_relative \
		common.proto

.PHONY: trpc
trpc: common
	trpc create \
		-p greeter.proto \
		--rpconly \
		-o trpc

.PHONY: grpc
grpc: common
	protoc --go_out=grpc --go_opt=paths=source_relative \
		--go-grpc_out=grpc --go-grpc_opt=paths=source_relative \
		greeter.proto

.PHONY: tgrpc
tgrpc: common
	trpc create \
		-p greeter.proto \
		--rpconly \
		-o tgrpc \
		--protocol=grpc \

.PHONY: tstreams
tstreams: common
	trpc create \
		-p streams.proto \
		--rpconly \
		-o tstreams \
		--protocol=grpc \
		--mock=false

.PHONY: gstreams
gstreams: common
	protoc --go_out=gstreams --go_opt=paths=source_relative \
		--go-grpc_out=gstreams --go-grpc_opt=paths=source_relative \
		streams.proto

.PHONY: clean
clean:
	find . -name "*.go" | xargs rm
