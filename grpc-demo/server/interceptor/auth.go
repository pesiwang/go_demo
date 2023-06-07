package interceptor

import (
	"context"
	"fmt"
	"server/base/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("Auth interceptor: before handler")
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("metadata[token]: ", md["token"])
		ctx = util.SetAuthInfo(ctx, &util.AuthInfo{Uid: 123456})
		fmt.Println("set authinfo.Uid to 123456")
	} else {
		fmt.Println("fetch metadata failed")
	}
	resp, err = handler(ctx, req)
	fmt.Println("Auth interceptor: after handler")

	return
}
