package order

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/mqbuild"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"
)

func Start(ctx context.Context, cfg *Config, _ discovery.SvcDiscoveryRegistry, _ grpc.ServiceRegistrar) error {
	db := dbbuild.NewBuilder(&cfg.MongoConfig, &cfg.RedisConfig)
	if _, err := db.Mongo(ctx); err != nil {
		return errs.WrapMsg(err, "init mongo failed")
	}
	if _, err := db.Redis(ctx); err != nil {
		return errs.WrapMsg(err, "init redis failed")
	}
	mq := mqbuild.NewBuilder(&cfg.KafkaConfig)
	if _, err := mq.GetTopicProducer(ctx, cfg.KafkaConfig.ToMongoTopic); err != nil {
		return errs.WrapMsg(err, "init producer failed")
	}
	<-ctx.Done()
	return ctx.Err()
}
