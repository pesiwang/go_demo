/*
https://www.iarno.cn/article/grpc/

https://grpc.io/docs/languages/go/quickstart/
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false *.proto
*/

package main

import (
	"fmt"
	"grpc-server/protodemo"
	"grpc-server/service"
	"log"
	"net"

	"google.golang.org/grpc"
)

const grpcPort = 6655

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	g := grpc.NewServer()
	protodemo.RegisterDemoServiceServer(g, &service.Demo{})

	fmt.Println("sever is listening....")

	if err = g.Serve(l); err != nil {
		log.Fatal(err)
	}

}
