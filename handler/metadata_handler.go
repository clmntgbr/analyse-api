package handler

import (
	"context"

	"go-api/infrastructure/config"
	rabbitmqDTO "go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	metadatauc "go-api/usecase/metadata"
	"go-api/usecase/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataHandler struct {
	env                         *config.Config
	parser                      *security.WorkerParser
	securityValidator           *security.WorkerSecurityValidator
	analyzeMediaMetadataUseCase *metadatauc.AnalyzeMediaMetadataUseCase
	createSignalUseCase         *signal.CreateSignalUseCase
}

func NewMetadataHandler(
	env *config.Config,
	parser *security.WorkerParser,
	securityValidator *security.WorkerSecurityValidator,
	analyzeMediaMetadataUseCase *metadatauc.AnalyzeMediaMetadataUseCase,
	createSignalUseCase *signal.CreateSignalUseCase,
) *MetadataHandler {
	return &MetadataHandler{
		env:                         env,
		parser:                      parser,
		securityValidator:           securityValidator,
		analyzeMediaMetadataUseCase: analyzeMediaMetadataUseCase,
		createSignalUseCase:         createSignalUseCase,
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

	result, err := h.analyzeMediaMetadataUseCase.Execute(ctx, payload.Message.UserID, payload.Message.MediaKey)
	if err != nil {
		return err
	}
	_, err = h.createSignalUseCase.Execute(ctx, payload.Message.MediaID, "metadata", result.Signal.Score, result.Signal.Confidence, result.Signal.Details)
	if err != nil {
		return err
	}

	return nil
}
