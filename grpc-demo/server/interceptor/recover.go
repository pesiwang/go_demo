package interceptor

import (
	"context"
	"fmt"
	"proto/common"
	"runtime"

	"google.golang.org/grpc"
)

func Recover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			printStack()
			resp = &common.BaseResp{
				Code: 5,
				Msg:  "server inner error",
			}
		}
	}()
	fmt.Println("Recover interceptor: before handler")
	resp, err = handler(ctx, req)
	fmt.Println("Recover interceptor: after handler")

	return
}

func printStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	stackMsg := string(buf[:n])
	fmt.Println(stackMsg)

	return stackMsg
}
