package rpcx

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/51xpx/pkgx/pkg/jsonx"
	"github.com/51xpx/pkgx/pkg/timex"
	"google.golang.org/grpc"
)

var grpcChannel = make(chan string, 100)

func ClientInterceptor() grpc.UnaryClientInterceptor {

	go handleGrpcChannel()

	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// 开始时间
		startTime := timex.GetCurrentMilliUnix()

		err := invoker(ctx, method, req, reply, cc, opts...)

		// 结束时间
		endTime := timex.GetCurrentMilliUnix()

		// 日志格式
		grpcLogMap := make(map[string]interface{})

		grpcLogMap["request_time"] = startTime
		grpcLogMap["request_data"] = req
		grpcLogMap["request_method"] = method

		grpcLogMap["response_data"] = reply
		grpcLogMap["response_error"] = err

		grpcLogMap["cost_time"] = fmt.Sprintf("%vms", endTime-startTime)

		grpcLogJson, _ := jsonx.Marshal(grpcLogMap)

		grpcChannel <- string(grpcLogJson)

		return err
	}
}

func handleGrpcChannel() {
	if f, err := os.OpenFile("app.Config.Server.Name", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		log.Println(err)
	} else {
		for accessLog := range grpcChannel {
			_, _ = f.WriteString(accessLog + "\n")
		}
	}
	return
}
