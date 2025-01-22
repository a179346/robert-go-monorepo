package post_board_applogger

import (
	"fmt"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/pkg/rabbitmqlogger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetAppLogger() (*rabbitmqlogger.RabbitMQLogger, error) {
	loggingConfig := post_board_config.GetLoggerConfig()
	if !loggingConfig.Enable {
		return nil, nil
	}

	rabbitMQConfig := post_board_config.GetRabbitMQConfig()

	conn, err := amqp.Dial(rabbitMQConfig.Url)
	if err != nil {
		return nil, fmt.Errorf("amqp.Dial error: %w", err)
	}

	return rabbitmqlogger.New(conn, loggingConfig.TargetExchange), nil
}
