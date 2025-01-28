package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

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
	loggingConfig.Enable = env_helper.GetBool("LOGGING_ENABLE", true)
	loggingConfig.TargetExchange = env_helper.GetString("LOGGING_TARGET_EXCHANGE", "logging-exchange")

	loggingConfig.ElasticSearchAddress = env_helper.GetString("LOGGING_ELASTICSEARCH_ADDRESS", "http://127.0.0.1:9200")
	loggingConfig.ElasticSearchIndexPrefix = env_helper.GetString("LOGGING_ELASTICSEARCH_INDEX_PREFIX", "post-board-api-")
	loggingConfig.ConsumerSourceQueue = env_helper.GetString("LOGGING_CONSUMER_SOURCE_QUEUE", "logging-queue")
	loggingConfig.ConsumerConcurrency = env_helper.GetInt("LOGGING_CONSUMER_CONCURRENCY", 2)
}

func GetLoggingConfig() LoggingConfig {
	return loggingConfig
}
