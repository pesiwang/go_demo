package service

import (
	"context"
	"fmt"
	"grpc-server/protodemo"
)

type Demo struct {
	protodemo.UnimplementedDemoServiceServer
}

func (c *Demo) GetData(ctx context.Context, msg *protodemo.DemoReq) (*protodemo.DemoRsp, error) {
	params := msg.Id
	fmt.Println("recv grpc client request:", params)
	reply := fmt.Sprintf("reply_%v", params)
	return &protodemo.DemoRsp{
		Result: reply,
	}, nil
}
