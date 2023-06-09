set -e
function generateAll() {
    protoc -I . -I .. -I ../../../../protobuf/src/ -I ../../../../googleapis/ \
        --go_out . --go_opt paths=source_relative \
        --go-grpc_out . --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
        --include_imports  --include_source_info --descriptor_set_out $2.pb  \
        *.proto

    protoc -I .  -I .. -I ../../../../protobuf/src/ -I ../../../../googleapis/ \
        --grpc-gateway_out . \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        *.proto

    protoc -I . -I ..  -I ../../../../protobuf/src/ -I ../../../../googleapis/ \
        --swagger_out=logtostderr=true:. *.proto
}
case $1 in
clean)
    echo "clean all generated files"
    for i in $(ls -F | grep '/$' | sed 's#/##g'); do
        echo $i
        cd $i
        rm -rf *.pb.go *.pb *.json *.pb.gw.go
        cd ..
    done
    ;;
*)
    echo "build all *.proto"
    for i in $(ls -F | grep '/$' | grep -v admin | sed 's#/##g'); do
        echo $i
        cd $i
        rm -rf *.pb.go *.pb
        generateAll *.proto $i
        cd ..
    done
    ;;
esac



# protoc -I . -I .. -I ../../../../protobuf/src/ -I ../../../../googleapis/ \
#     --go_out . --go_opt paths=source_relative \
#     --go-grpc_out . --go-grpc_opt=paths=source_relative --go-grpc_opt=require_unimplemented_servers=false \
#     *.proto

# protoc -I .  -I .. -I ../../../../protobuf/src/ -I ../../../../googleapis/ \
#         --grpc-gateway_out . \
#         --grpc-gateway_opt logtostderr=true \
#         --grpc-gateway_opt paths=source_relative \
#         --grpc-gateway_opt generate_unbound_methods=true \
#         *.proto

#protoc -I .  -I ../../../protobuf/src/ -I ../../../googleapis/ --go_out=. --go_opt paths=source_relative --go-grpc_out . --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --go-grpc_opt paths=source_relative,require_unimplemented_servers=false bizdemo/bizdemo.proto

