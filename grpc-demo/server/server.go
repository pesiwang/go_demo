/*
https://www.iarno.cn/article/grpc/

https://grpc.io/docs/languages/go/quickstart/
在 proto 根目录执行
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false *.proto
*/

package main

import (
	"fmt"
	"log"
	"net"
	"proto/bizdemo"
	"server/interceptor"
	"server/service"

	gpm "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

const grpcPort = 6655

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			gpm.ChainUnaryServer(
				interceptor.Recover,
				interceptor.Auth,
				interceptor.ServerLog,
			),
		),
	)

	bizdemo.RegisterBizDemoServer(grpcServer, &service.BizDemo{})

	fmt.Println("sever is listening....")

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
