package migrations

import (
	"context"
	"database/sql"
	"fmt"

	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/pressly/goose/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	goose.AddMigrationContext(Up1737375716, Down1737375716)
}

const deadLoggingExchange = "dead-logging-exchange"
const deadLoggingQueue = "dead-logging-queue"
const loggingExchange = "logging-exchange"
const fileLoggingQueue = "file-logging-queue"

func Up1737375716(ctx context.Context, tx *sql.Tx) error {
	conn, err := amqp.Dial(post_board_config.GetRabbitMQConfig().Url)
	if err != nil {
		return fmt.Errorf("amqp.Dial error: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel error: %w", err)
	}

	err = ch.ExchangeDeclare(
		deadLoggingExchange,
		"fanout",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDeclare(%v) error: %w", deadLoggingExchange, err)
	}

	_, err = ch.QueueDeclare(
		deadLoggingQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:     amqp.QueueTypeClassic,
			amqp.QueueOverflowArg: amqp.QueueOverflowRejectPublish,
		},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare(%v) error: %w", deadLoggingQueue, err)
	}

	err = ch.QueueBind(
		deadLoggingQueue,
		deadLoggingQueue,
		deadLoggingExchange,
		false,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueBind(%v) error: %w", deadLoggingQueue, err)
	}

	err = ch.ExchangeDeclare(
		loggingExchange,
		"fanout",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDeclare(%v) error: %w", loggingExchange, err)
	}

	_, err = ch.QueueDeclare(
		fileLoggingQueue,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:        amqp.QueueTypeQuorum,
			amqp.QueueOverflowArg:    amqp.QueueOverflowRejectPublish,
			"x-dead-letter-exchange": deadLoggingExchange,
			"x-delivery-limit":       10,
		},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare(%v) error: %w", fileLoggingQueue, err)
	}

	err = ch.QueueBind(
		fileLoggingQueue,
		fileLoggingQueue,
		loggingExchange,
		false,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueBind(%v) error: %w", fileLoggingQueue, err)
	}

	return nil
}

func Down1737375716(ctx context.Context, tx *sql.Tx) error {
	conn, err := amqp.Dial(post_board_config.GetRabbitMQConfig().Url)
	if err != nil {
		return fmt.Errorf("amqp.Dial error: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel error: %w", err)
	}

	err = ch.QueueUnbind(
		fileLoggingQueue,
		fileLoggingQueue,
		loggingExchange,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueUnbind(%v) error: %w", fileLoggingQueue, err)
	}

	err = ch.QueueUnbind(
		deadLoggingQueue,
		deadLoggingQueue,
		deadLoggingExchange,
		amqp.Table{},
	)
	if err != nil {
		return fmt.Errorf("ch.QueueUnbind(%v) error: %w", deadLoggingQueue, err)
	}

	err = ch.ExchangeDelete(
		loggingExchange,
		true,
		false,
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDelete(%v) error: %w", loggingExchange, err)
	}

	_, err = ch.QueueDelete(
		fileLoggingQueue,
		false,
		false,
		false,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDelete(%v) error: %w", fileLoggingQueue, err)
	}

	err = ch.ExchangeDelete(
		deadLoggingExchange,
		true,
		false,
	)
	if err != nil {
		return fmt.Errorf("ch.ExchangeDelete(%v) error: %w", deadLoggingExchange, err)
	}

	_, err = ch.QueueDelete(
		deadLoggingQueue,
		true,
		true,
		false,
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDelete(%v) error: %w", deadLoggingQueue, err)
	}

	return nil
}
