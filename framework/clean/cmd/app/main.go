package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/openimsdk/open-im-server/v3/framework/clean/internal/httpapi"
	"github.com/openimsdk/open-im-server/v3/framework/clean/internal/rpc"
	"github.com/openimsdk/open-im-server/v3/framework/clean/internal/socket"
	"github.com/openimsdk/open-im-server/v3/framework/clean/pkg/common"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := common.LoadConfig()
	broker := common.NewBroker()
	defer broker.Close()
	store := httpapi.NewInMemoryStore()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return rpc.StartServer(ctx, cfg.RPCAddr, "[rpc]") })
	g.Go(func() error { return socket.StartServer(ctx, cfg.SocketAddr, "[socket]", broker) })
	g.Go(func() error { return httpapi.StartServer(ctx, cfg.HTTPAddr, cfg.RPCAddr, "[http]", broker, store) })

	if err := g.Wait(); err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
