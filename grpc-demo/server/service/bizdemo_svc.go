package service

import (
	"context"
	"fmt"
	"proto/bizdemo"
	"proto/common"

	"google.golang.org/protobuf/types/known/anypb"
)

type BizDemo struct {
	//	protodemo.UnimplementedDemoServiceServer
}

func (c *BizDemo) GetData(ctx context.Context, msg *bizdemo.DemoReq) (*bizdemo.DemoResp, error) {
	params := msg.Id
	fmt.Println("recv grpc client GetData request, params", params)
	reply := fmt.Sprintf("reply_%v", params)
	// time.Sleep(3 * time.Second)

	testResp := &common.TestResp{
		I:    7,
		Data: "any data",
	}

	respData, _ := anypb.New(testResp)
	return &bizdemo.DemoResp{
		Result: reply,
		Data:   respData,
	}, nil
}

func (c *BizDemo) Test(ctx context.Context, msg *common.TestReq) (*common.TestResp, error) {
	i := msg.I
	data := msg.Data
	fmt.Printf("recv grpc client request,i=%v, data=%v\n", i, data)
	replyData := fmt.Sprintf("reply_%v", data)
	// time.Sleep(3 * time.Second)

	if i == 777 {
		panic("test panic recover")
	}

	return &common.TestResp{
		I:    i,
		Data: replyData,
	}, nil
}
