/*
https://www.iarno.cn/article/grpc/

https://grpc.io/docs/languages/go/quickstart/
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false *.proto
*/

package main

import (
	"context"
	"fmt"
	"grpc-client/protodemo"
	"time"

	"google.golang.org/grpc"
)

const (
	grpcAddr = "127.0.0.1:6655"
)

func main() {
	c, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	rsp, err := protodemo.NewDemoServiceClient(c).
		GetData(ctx, &protodemo.DemoReq{
			Id: "123",
		})

	if err != nil {
		fmt.Printf("grpc call error: %s", err)
	} else {
		fmt.Printf("grpc call succ: %s", rsp)
	}
}
