package rabbitmq

import (
	"context"
	"fmt"
	"log"

	"go-api/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataMessageHandler interface {
	HandleMessage(ctx context.Context, delivery *amqp.Delivery) error
}

type MetadataWorker struct {
	env     *config.Config
	handler MetadataMessageHandler

	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewMetadataWorker(
	env *config.Config,
	handler MetadataMessageHandler,
) *MetadataWorker {
	return &MetadataWorker{
		env:     env,
		handler: handler,
	}
}

func (c *MetadataWorker) Start() error {
	conn, err := dialWithRetry(c.env.RabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}

	c.channel = channel

	if err := c.channel.ExchangeDeclare(
		c.env.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf(
			"failed to declare exchange %q: %w",
			c.env.ExchangeName,
			err,
		)
	}

	queue, err := c.channel.QueueDeclare(
		c.env.MetadataQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to declare queue %q: %w",
			c.env.MetadataQueueName,
			err,
		)
	}

	messages, err := c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to consume queue %q: %w",
			queue.Name,
			err,
		)
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Printf(
		"[*] Waiting for messages on exchange %q",
		c.env.ExchangeName,
	)

	for message := range messages {
		if err := c.handler.HandleMessage(context.Background(), &message); err != nil {
			log.Printf(
				"rejected message (routing key: %q): %v, body: %s",
				message.RoutingKey,
				err,
				message.Body,
			)

			if nackErr := message.Nack(false, false); nackErr != nil {
				log.Printf("failed to nack message: %v", nackErr)
			}

			continue
		}

		if ackErr := message.Ack(false); ackErr != nil {
			log.Printf("failed to ack message: %v", ackErr)
		}
	}

	return nil
}

func (c *MetadataWorker) Stop() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
