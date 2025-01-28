package post_board_apilogger

import (
	"errors"
	"fmt"
	"sync"

	"github.com/a179346/robert-go-monorepo/pkg/rabbitmqlogger"
	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetApiLogger() (*rabbitmqlogger.RabbitMQLogger, error) {
	loggingConfig := post_board_config.GetLoggingConfig()
	if !loggingConfig.Enable {
		return nil, nil
	}

	rabbitMQConfig := post_board_config.GetRabbitMQConfig()

	var conn *amqp.Connection
	err := errors.New("amqp connection have not been initialized")

	getAmqpConnection := func() (*amqp.Connection, error) {
		if err == nil && !conn.IsClosed() {
			return conn, nil
		}

		var mu sync.Mutex
		mu.Lock()
		defer mu.Unlock()

		if err == nil && !conn.IsClosed() {
			return conn, nil
		}
		conn, err = amqp.Dial(rabbitMQConfig.Url)
		if err != nil {
			err = fmt.Errorf("amqp.Dial error: %w", err)
		}
		return conn, err
	}

	return rabbitmqlogger.New(getAmqpConnection, loggingConfig.TargetExchange), nil
}
