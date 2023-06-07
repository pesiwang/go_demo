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

	"google.golang.org/grpc"
)

const (
	grpcAddr = "127.0.0.1:6655"
)

func main() {
	grpcClientConn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(grpcClientConn)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	grpcClient := bizdemo.NewBizDemoClient(grpcClientConn)

	resp, err := grpcClient.GetData(ctx, &bizdemo.DemoReq{
		Id: "123",
	})

	if err != nil {
		fmt.Printf("grpc get data call error: %s\n", err)
	} else {
		fmt.Printf("grpc get data call succ: %s\n", resp)
	}

	resp2, err2 := grpcClient.Test(ctx, &common.TestReq{
		I:    777,
		Data: "this is data field",
	})

	if err != nil {
		fmt.Printf("grpc  call error: %s\n", err2)
	} else {
		fmt.Printf("grpc call succ: %s\n", resp2)
	}
}
