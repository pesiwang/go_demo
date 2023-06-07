package interceptor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func ServerLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("server log interceptor: before handler")
	start := time.Now()
	resp, err = handler(ctx, req)
	timeCost := time.Since(start)
	fmt.Printf("server log interceptor: after handler, time cost %v milliseconds\n", timeCost.Milliseconds())

	return
}
