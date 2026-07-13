package rabbitmq

import (
	"context"
	"sync"

	"go-api/infrastructure/config"
)

type lazyPublisher struct {
	env *config.Config

	mu        sync.Mutex
	publisher Publisher
}

func NewLazyPublisherFromEnv(env *config.Config) Publisher {
	return &lazyPublisher{env: env}
}

func (p *lazyPublisher) Publish(ctx context.Context, queueName string, message any) error {
	publisher, err := p.getPublisher()
	if err != nil {
		return err
	}

	return publisher.Publish(ctx, queueName, message)
}

func (p *lazyPublisher) getPublisher() (Publisher, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.publisher != nil {
		return p.publisher, nil
	}

	publisher, err := NewPublisherFromEnv(p.env)
	if err != nil {
		return nil, err
	}

	p.publisher = publisher
	return p.publisher, nil
}
