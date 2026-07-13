package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"go-api/infrastructure/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	PublishMetadataEvent(ctx context.Context, config *config.Config, event MetadataEvent) error
}

type publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) Publisher {
	return &publisher{
		channel: channel,
	}
}

func NewPublisherFromEnv(env *config.Config) (Publisher, error) {
	conn, err := dialWithRetry(env.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ at %s: %w", env.RabbitMQURL, err)
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	return NewPublisher(ch), nil
}

func (p *publisher) PublishMetadataEvent(ctx context.Context, config *config.Config, event MetadataEvent) error {
	message := MessagePayload{
		SecretKey: config.RabbitMQSecretKey,
		Message:   event,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		"",
		config.MetadataQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
