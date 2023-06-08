package service

import (
	"context"
	"fmt"
	"proto/bizdemo"
	"proto/common"
	"server/base/util"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
)

type BizDemo struct {
	//	protodemo.UnimplementedDemoServiceServer
}

func (c *BizDemo) GetData(ctx context.Context, msg *bizdemo.DemoReq) (*bizdemo.DemoResp, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("metadata[token]: ", md["token"])
	} else {
		fmt.Println("fetch metadata failed")
	}
	paramId := msg.Id
	fmt.Println("recv grpc client GetData request, params", paramId)
	reply := fmt.Sprintf("reply:%v", paramId)
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
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("metadata[token]: ", md["token"])
	} else {
		fmt.Println("fetch metadata failed")
	}

	authInfo := util.GetAuthInfo(ctx)
	fmt.Println("auth info from ctx:", authInfo.Uid)

	i := msg.I
	data := msg.Data
	fmt.Printf("recv grpc client request,i=%v, data=%v\n", i, data)
	replyData := fmt.Sprintf("reply:%v", data)
	time.Sleep(3 * time.Second)

	return &common.TestResp{
		I:    i,
		Data: replyData,
	}, nil
}
