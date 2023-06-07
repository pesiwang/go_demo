
# run in proto directory
protoc -I . -I ..\..\..\protobuf\src\ -I ..\..\..\googleapis\ --go_out . --go_opt paths=source_relative --go-grpc_out . --grpc-gateway_out .  --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --go-grpc_opt paths=source_relative,require_unimplemented_servers=false common/common.proto




protoc -I .  -I ..\..\..\protobuf\src\ -I ..\..\..\googleapis\ --go_out=. --go_opt paths=source_relative --go-grpc_out . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --go-grpc_opt paths=source_relative,require_unimplemented_servers=false bizdemo/bizdemo.proto
