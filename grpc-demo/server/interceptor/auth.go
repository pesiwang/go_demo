package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("Auth interceptor: before handler")
	resp, err = handler(ctx, req)
	fmt.Println("Auth interceptor: after handler")

	return
}
