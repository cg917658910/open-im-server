package order

import "github.com/openimsdk/open-im-server/v3/pkg/common/config"

type Config struct {
	Discovery  config.Discovery
	Share      config.Share
	RedisConfig config.Redis
	MongoConfig config.Mongo
	KafkaConfig config.Kafka
	RpcConfig   config.OpenImRPC
}
