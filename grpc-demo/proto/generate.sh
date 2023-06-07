
# run in proto directory
protoc -I . -I ..\..\..\protobuf\src\ --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false common/common.proto


protoc -I .  -I ..\..\..\protobuf\src\ --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false bizdemo/bizdemo.proto
