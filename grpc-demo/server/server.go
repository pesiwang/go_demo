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
	"math/rand"
	"net"
	"net/http"
	"proto/bizdemo"
	"server/config"
	"server/interceptor"
	"server/service"
	"time"

	_ "net/http/pprof"

	gpm "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// pprof 分析命令
// 获取分析文件
// go tool pprof http://127.0.0.1:6061/debug/pprof/profile?seconds=60

// 进行分析：
// go tool pprof -http=:8080 ./pprof.server.exe.samples.cpu.003.pb.gz

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

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(serverConfig.HttpServerAddr, mux)
}

var datas []string

func Add(str string) int {
	data := []byte(str)
	datas = append(datas, string(data))
	return len(datas)
}

func startPprof() {
	go func() {
		for {
			log.Printf("len: %d\n", Add("go-programming-tour-book"))
			time.Sleep(time.Second * 3)
		}
	}()

	go func() {
		for {
			j := 0
			for i := 0; i <= 1000000; i++ {
				j = i*j + i + rand.Int()

			}
			Add("ddddd")
			fmt.Printf("j: %v\n", j)
			time.Sleep(time.Second * 3)
		}
	}()

	_ = http.ListenAndServe("0.0.0.0:6061", nil)
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
	go startPprof()

	bizdemo.RegisterBizDemoServer(grpcServer, &service.BizDemo{})

	fmt.Println("sever is listening....")

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
