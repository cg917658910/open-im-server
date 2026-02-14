package cmd

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/order"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/system/program"
	"github.com/spf13/cobra"
)

type OrderRpcCmd struct {
	*RootCmd
	ctx         context.Context
	configMap   map[string]any
	orderConfig *order.Config
}

func NewOrderRpcCmd() *OrderRpcCmd {
	var c order.Config
	ret := &OrderRpcCmd{orderConfig: &c}
	ret.configMap = map[string]any{
		config.DiscoveryConfigFilename: &c.Discovery,
		config.ShareFileName:           &c.Share,
		config.RedisConfigFileName:     &c.RedisConfig,
		config.MongodbConfigFileName:   &c.MongoConfig,
		config.KafkaConfigFileName:     &c.KafkaConfig,
		config.OpenIMRPCMsgCfgFileName: &c.RpcConfig,
	}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(ret.configMap))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return ret.runE()
	}
	return ret
}

func (o *OrderRpcCmd) Exec() error { return o.Execute() }

func (o *OrderRpcCmd) runE() error {
	return startrpc.Start(
		o.ctx,
		&o.orderConfig.Discovery,
		&o.orderConfig.RpcConfig.CircuitBreaker,
		&o.orderConfig.RpcConfig.RateLimiter,
		&o.orderConfig.RpcConfig.Prometheus,
		o.orderConfig.RpcConfig.RPC.ListenIP,
		o.orderConfig.RpcConfig.RPC.RegisterIP,
		o.orderConfig.RpcConfig.RPC.AutoSetPorts,
		o.orderConfig.RpcConfig.RPC.Ports,
		o.Index(),
		o.orderConfig.Discovery.RpcService.Msg,
		nil,
		o.orderConfig,
		nil,
		nil,
		order.Start,
	)
}
