package handler

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataHandler struct {
}

func NewMetadataHandler() *MetadataHandler {
	return &MetadataHandler{}
}

func (h *MetadataHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	fmt.Println("Received metadata message:", string(message.Body))
	return nil
}
