package rpcx

import (
	"time"

	"google.golang.org/grpc"
)

// GRPC .
type GRPC struct {
	Address string
	Conn    *grpc.ClientConn
}

// NewGRPC new custom grpc
func NewGRPC(address string) (*GRPC, error) {
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithTimeout(2*time.Second),
		grpc.WithBackoffMaxDelay(time.Second),
	)
	if err != nil {
		return nil, err
	}
	return &GRPC{
		Address: address,
		Conn:    conn,
	}, nil
}

// Close close connection
func (g *GRPC) Close() error {
	return g.Conn.Close()
}
