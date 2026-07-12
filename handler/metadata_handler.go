package handler

import (
	"context"
	"log"

	"go-api/infrastructure/config"
	rabbitmqDTO "go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataHandler struct {
	env               *config.Config
	parser            *security.WorkerParser
	securityValidator *security.WorkerSecurityValidator
}

func NewMetadataHandler(env *config.Config, parser *security.WorkerParser, securityValidator *security.WorkerSecurityValidator) *MetadataHandler {
	return &MetadataHandler{
		env:               env,
		parser:            parser,
		securityValidator: securityValidator,
	}
}

func (h *MetadataHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	var payload rabbitmqDTO.MessagePayload
	if err := h.parser.ParseAndValidate(message.Body, &payload); err != nil {
		return err
	}

	if err := h.securityValidator.Validate(payload.SecretKey); err != nil {
		return err
	}

	log.Println("🔄 Received message", payload)

	return nil
}
