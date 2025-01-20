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

func initRabbitMQConfig() {
	rabbitMQConfig.Host = env_helper.GetStringEnv("RABBITMQ_HOST", "localhost")
	rabbitMQConfig.Port = env_helper.GetIntEnv("RABBITMQ_PORT", 5672)
	rabbitMQConfig.User = env_helper.GetStringEnv("RABBITMQ_USER", "post-board-user")
	rabbitMQConfig.Password = env_helper.GetStringEnv("RABBITMQ_PASSWORD", "mymqpass")

	rabbitMQConfig.Url = fmt.Sprintf(
		"amqp://%v:%v@%v:%v/",
		rabbitMQConfig.User,
		rabbitMQConfig.Password,
		rabbitMQConfig.Host,
		rabbitMQConfig.Port,
	)
}

func GetRabbitMQConfig() RabbitMQConfig {
	initAll()
	return rabbitMQConfig
}
