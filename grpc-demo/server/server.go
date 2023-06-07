/*
https://www.iarno.cn/article/grpc/

https://grpc.io/docs/languages/go/quickstart/
在 proto 根目录执行
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false *.proto
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"proto/bizdemo"
	"server/interceptor"
	"server/service"

	gpm "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const grpcPort = 6655

func startHttpServer(grpcServer *grpc.Server) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := bizdemo.RegisterBizDemoHandlerFromEndpoint(ctx, mux, "localhost:6655", opts)
	if err != nil {
		return err
	}

	fmt.Println("http server is listening....")

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

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

	go startHttpServer(grpcServer)

	bizdemo.RegisterBizDemoServer(grpcServer, &service.BizDemo{})

	fmt.Println("sever is listening....")

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
