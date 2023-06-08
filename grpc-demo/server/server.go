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
	"server/config"
	"server/interceptor"
	"server/service"

	gpm "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func startHttpServer(grpcServer *grpc.Server, serverConfig *config.ServerConfig) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := bizdemo.RegisterBizDemoHandlerFromEndpoint(ctx, mux, serverConfig.LocalGrpcServerAddr, opts)
	if err != nil {
		return err
	}

	fmt.Println("http server is listening....")

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(serverConfig.HttpServerAddr, mux)
}

func main() {
	config.Init()
	serverConfig := config.GetServerConfig()
	listener, err := net.Listen("tcp", serverConfig.GrpcServerAddr)
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

	go startHttpServer(grpcServer, serverConfig)

	bizdemo.RegisterBizDemoServer(grpcServer, &service.BizDemo{})

	fmt.Println("sever is listening....")

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
