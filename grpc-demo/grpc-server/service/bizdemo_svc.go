package service

import (
	"context"
	"fmt"
	"proto/bizdemo"
	"proto/common"
)

type BizDemo struct {
	//	protodemo.UnimplementedDemoServiceServer
}

func (c *BizDemo) GetData(ctx context.Context, msg *bizdemo.DemoReq) (*bizdemo.DemoResp, error) {
	params := msg.Id
	fmt.Println("recv grpc client GetData request, params", params)
	reply := fmt.Sprintf("reply_%v", params)
	// time.Sleep(3 * time.Second)

	return &bizdemo.DemoResp{
		Result: reply,
	}, nil
}

func (c *BizDemo) Test(ctx context.Context, msg *common.TestReq) (*common.TestResp, error) {
	i := msg.I
	data := msg.Data
	fmt.Printf("recv grpc client request,i=%v, data=%v\n", i, data)
	replyData := fmt.Sprintf("reply_%v", data)
	// time.Sleep(3 * time.Second)

	return &common.TestResp{
		I:    i,
		Data: replyData,
	}, nil
}
