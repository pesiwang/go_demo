/*
https://www.iarno.cn/article/grpc/

https://grpc.io/docs/languages/go/quickstart/
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false *.proto
*/

package main

import (
	"context"
	"fmt"
	"proto/bizdemo"
	"proto/common"
	"time"

	gpm "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientLog(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md := metadata.Pairs("token", "this is token string")
	ctx = metadata.NewOutgoingContext(ctx, md)

	fmt.Println("client log interceptor: before handler")
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	timeCost := time.Since(start)
	fmt.Printf("client log interceptor: after handler, time cost:%v milliseconds\n", timeCost.Milliseconds())

	return err
}

func main() {

	grpcClientConn, err := grpc.Dial(
		"127.0.0.1:9000",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			gpm.ChainUnaryClient(
				ClientLog,
			),
		),
	)
	if err != nil {
		panic(grpcClientConn)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	// defer cancel()
	ctx := context.Background()

	grpcClient := bizdemo.NewBizDemoClient(grpcClientConn)

	resp, err := grpcClient.GetData(ctx, &bizdemo.DemoReq{
		Id: "123",
	})

	if err != nil {
		fmt.Printf("grpc get data call error: %s\n", err)
	} else {
		if resp.Data.MessageIs(&common.TestResp{}) {
			fmt.Println("resp.Data is common.TestResp type")
		} else {
			fmt.Printf("resp.Data is unknown type:%v\n", resp.Data)
		}

		testResp := &common.TestResp{}
		anyErr := resp.Data.UnmarshalTo(testResp)
		if anyErr != nil {
			fmt.Printf("grpc get data call, resp.Data.UnmarshalTo(testResp) failed: %v, %v\n", anyErr, testResp)
		} else {
			fmt.Printf("grpc get data call succ, resp.Data: %v\n", testResp)
		}
		fmt.Printf("grpc get data call succ, resp: %s\n", resp)
	}

	resp2, err2 := grpcClient.Test(ctx, &common.TestReq{
		I:    777,
		Data: "this is data field",
	})

	if err2 != nil {
		fmt.Printf("grpc  call error: %s\n", err2)
	} else {
		fmt.Printf("grpc call succ: %s\n", resp2)
	}
}
