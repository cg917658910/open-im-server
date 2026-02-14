# RPC 工程框架（gRPC + cmd + infra builder）

## 1. 对标本仓库的参考实现

- 统一启动器：`pkg/common/startrpc/start.go`
- 命令入口：`pkg/common/cmd/*.go`、`cmd/openim-rpc/openim-rpc-*/main.go`
- DB 构建：`pkg/dbbuild`
- MQ 构建：`pkg/mqbuild`
- 错误码：`pkg/common/servererrs`

## 2. 工程目录（建议）

```text
cmd/openim-rpc/openim-rpc-order/main.go
pkg/common/cmd/order.go
internal/rpc/order/
  config.go
  server.go
  repo.go
  service.go
  callback.go
  notification.go
  sync.go
```

## 3. cmd 层模板（生产可用）

> 要点：cmd 只做配置装配、启动参数和调用统一 Start，不写业务逻辑。

```go
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
	ctx    context.Context
	config *order.Config
}

func NewOrderRpcCmd() *OrderRpcCmd {
	var c order.Config
	ret := &OrderRpcCmd{config: &c}
	ret.RootCmd = NewRootCmd(program.GetProcessName(), WithConfigMap(map[string]any{
		config.OpenIMRPCUserCfgFileName: &c.RpcConfig,
		config.RedisConfigFileName:      &c.RedisConfig,
		config.MongodbConfigFileName:    &c.MongoConfig,
		config.KafkaConfigFileName:      &c.KafkaConfig,
		config.DiscoveryConfigFilename:  &c.Discovery,
		config.ShareFileName:            &c.Share,
	}))
	ret.ctx = context.WithValue(context.Background(), "version", version.Version)
	ret.Command.RunE = func(cmd *cobra.Command, args []string) error { return ret.runE() }
	return ret
}

func (o *OrderRpcCmd) runE() error {
	return startrpc.Start(
		o.ctx,
		&o.config.Discovery,
		&o.config.RpcConfig.CircuitBreaker,
		&o.config.RpcConfig.RateLimiter,
		&o.config.RpcConfig.Prometheus,
		o.config.RpcConfig.RPC.ListenIP,
		o.config.RpcConfig.RPC.RegisterIP,
		o.config.RpcConfig.RPC.AutoSetPorts,
		o.config.RpcConfig.RPC.Ports,
		o.Index(),
		o.config.Discovery.RpcService.User,
		nil,
		o.config,
		nil,
		nil,
		order.Start,
	)
}
```

## 4. 业务层模板（DB + MQ + 错误码）

```go
package order

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/mqbuild"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func Start(ctx context.Context, cfg *Config, _ any, _ any) error {
	db := dbbuild.NewBuilder(&cfg.MongoConfig, &cfg.RedisConfig)
	_, err := db.Mongo(ctx)
	if err != nil {
		return errs.WrapMsg(err, "init mongo failed")
	}

	mq := mqbuild.NewBuilder(&cfg.KafkaConfig)
	producer, err := mq.GetTopicProducer(ctx, "order-topic")
	if err != nil {
		return errs.WrapMsg(err, "init mq producer failed")
	}
	_ = producer

	log.ZInfo(ctx, "order rpc started")
	// 业务错误映射示例：
	_ = servererrs.DatabaseError
	return nil
}
```

## 5. 工程约束（强烈建议）

1. **配置优先**：端口、topic、超时、开关全部配置化。
2. **依赖注入**：repo/service 初始化集中在 Start，不在 handler 临时创建。
3. **错误分层**：内部错误（errs.Wrap）+ 对外错误码（servererrs）。
4. **异步解耦**：通知/回调走 MQ，不阻塞主 RPC 链路。
5. **统一观测**：QPS、P99、错误率、慢查询、重试次数。

## 6. 上线检查清单

- [ ] 支持 SIGTERM 优雅退出。
- [ ] DB/MQ 初始化失败快速失败，不假启动。
- [ ] grpc 接口有超时和参数校验。
- [ ] 熔断/限流参数可调整并验证生效。
- [ ] 回调与消息消费可重放、幂等。


## 7. 脚手架文件（可直接复制）

已提供：

- `framework/rpc/scaffold/cmd/openim-rpc-order/main.go.tpl`
- `framework/rpc/scaffold/pkg/common/cmd/order.go.tpl`
- `framework/rpc/scaffold/internal/rpc/order/config.go.tpl`
- `framework/rpc/scaffold/internal/rpc/order/start.go.tpl`
- `framework/rpc/scaffold/config/openim-rpc-order.yml.tpl`
