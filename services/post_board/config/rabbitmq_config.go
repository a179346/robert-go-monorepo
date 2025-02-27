package post_board_config

import (
	"fmt"

	"github.com/a179346/robert-go-monorepo/pkg/env_helper"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Url      string
}

var rabbitMQConfig RabbitMQConfig

func init() {
	rabbitMQConfig.Host = env_helper.GetString("RABBITMQ_HOST", "localhost")
	rabbitMQConfig.Port = env_helper.GetInt("RABBITMQ_PORT", 5672)
	rabbitMQConfig.User = env_helper.GetString("RABBITMQ_USER", "post-board-user")
	rabbitMQConfig.Password = env_helper.GetString("RABBITMQ_PASSWORD", "mymqpass")

	rabbitMQConfig.Url = fmt.Sprintf(
		"amqp://%v:%v@%v:%v/",
		rabbitMQConfig.User,
		rabbitMQConfig.Password,
		rabbitMQConfig.Host,
		rabbitMQConfig.Port,
	)
}

func GetRabbitMQConfig() RabbitMQConfig {
	return rabbitMQConfig
}
