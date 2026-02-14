package rpc

import (
	"context"
	"net"
	netrpc "net/rpc"

	"github.com/openimsdk/open-im-server/v3/framework/clean/pkg/common"
)

func StartServer(ctx context.Context, addr string, loggerPrefix string) error {
	logger := common.NewLogger(loggerPrefix)
	server := netrpc.NewServer()
	if err := server.RegisterName("DemoService", &DemoService{}); err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Printf("rpc server listening at %s", addr)

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				logger.Printf("rpc server stopped")
				return nil
			default:
				logger.Printf("rpc accept error: %v", err)
				continue
			}
		}
		go server.ServeConn(conn)
	}
}
