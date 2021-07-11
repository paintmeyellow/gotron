protoc -I=./proto/tron -I=./proto/googleapis --go_out=plugins=grpc,paths=source_relative:. ./proto/tron/api/*.proto
protoc -I=./proto/tron -I=./proto/googleapis --go_out=plugins=grpc,paths=source_relative:. ./proto/tron/core/*.proto
protoc -I=./proto/tron -I=./proto/googleapis --go_out=plugins=grpc,paths=source_relative:. ./proto/tron/core/contract/*.proto