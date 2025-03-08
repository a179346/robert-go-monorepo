package post_board_config

import (
	"fmt"

	"github.com/a179346/robert-go-monorepo/pkg/envhelper"
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
	rabbitMQConfig.Host = envhelper.GetString("RABBITMQ_HOST", "localhost")
	rabbitMQConfig.Port = envhelper.GetInt("RABBITMQ_PORT", 5672)
	rabbitMQConfig.User = envhelper.GetString("RABBITMQ_USER", "post-board-user")
	rabbitMQConfig.Password = envhelper.GetString("RABBITMQ_PASSWORD", "mymqpass")

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
