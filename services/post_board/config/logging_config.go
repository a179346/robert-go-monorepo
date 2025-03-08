package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type LoggingConfig struct {
	Enable         bool
	TargetExchange string

	ElasticSearchAddress     string
	ElasticSearchIndexPrefix string
	ConsumerSourceQueue      string
	ConsumerConcurrency      int
}

var loggingConfig LoggingConfig

func init() {
	loggingConfig.Enable = envhelper.GetBool("LOGGING_ENABLE", true)
	loggingConfig.TargetExchange = envhelper.GetString("LOGGING_TARGET_EXCHANGE", "logging-exchange")

	loggingConfig.ElasticSearchAddress = envhelper.GetString("LOGGING_ELASTICSEARCH_ADDRESS", "http://127.0.0.1:9200")
	loggingConfig.ElasticSearchIndexPrefix = envhelper.GetString("LOGGING_ELASTICSEARCH_INDEX_PREFIX", "post-board-api-")
	loggingConfig.ConsumerSourceQueue = envhelper.GetString("LOGGING_CONSUMER_SOURCE_QUEUE", "logging-queue")
	loggingConfig.ConsumerConcurrency = envhelper.GetInt("LOGGING_CONSUMER_CONCURRENCY", 2)
}

func GetLoggingConfig() LoggingConfig {
	return loggingConfig
}
