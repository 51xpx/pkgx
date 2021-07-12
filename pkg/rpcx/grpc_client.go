package rpcx

import (
	"context"
	"fmt"
	"time"

	"github.com/51xpx/pkgx/pkg/jaeger_trace"
	"github.com/gin-gonic/gin"
	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewGrpcClientConn(serviceAddress string, c *gin.Context) *grpc.ClientConn {

	var conn *grpc.ClientConn
	var err error
	var jaegerTraceIsOpen int

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	if jaegerTraceIsOpen == 1 {

		tracer, _ := c.Get("Tracer")
		parentSpanContext, _ := c.Get("ParentSpanContext")

		conn, err = grpc.DialContext(
			ctx,
			serviceAddress,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(
				grpc_middeware.ChainUnaryClient(
					jaeger_trace.ClientInterceptor(tracer.(opentracing.Tracer), parentSpanContext.(opentracing.SpanContext)),
					ClientInterceptor(), // log
				),
			),
		)
	} else {
		conn, err = grpc.DialContext(
			ctx,
			serviceAddress,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(
				grpc_middeware.ChainUnaryClient(
					ClientInterceptor(), // log
				),
			),
		)
	}

	if err != nil {
		fmt.Println(serviceAddress, "grpc conn err:", err)
	}
	return conn
}
