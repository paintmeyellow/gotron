protoc -I=./proto/tron -I=./proto/googleapis --go-grpc_out=paths=source_relative:./pkg/proto --go_out=paths=source_relative:./pkg/proto ./proto/tron/api/*.proto
protoc -I=./proto/tron -I=./proto/googleapis --go_out=paths=source_relative:./pkg/proto ./proto/tron/core/*.proto
protoc -I=./proto/tron -I=./proto/googleapis --go_out=paths=source_relative:./pkg/proto ./proto/tron/core/contract/*.proto
mv pkg/proto/core/contract/* pkg/proto/core/
rm -rf pkg/proto/core/contract
