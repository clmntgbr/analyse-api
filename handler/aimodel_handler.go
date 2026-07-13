package handler

import (
	"context"
	"log"

	rabbitmqDTO "go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	pipelineuc "go-api/usecase/pipeline"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AiModelHandler struct {
	worker *StageWorkerHandler
}

func NewAiModelHandler(
	parser *security.WorkerParser,
	securityValidator *security.WorkerSecurityValidator,
	dispatcher *pipelineuc.Dispatcher,
) *AiModelHandler {
	process := func(_ context.Context, message rabbitmqDTO.AnalyzeMessage) error {
		log.Printf("ai_model worker processing media %s", message.MediaID)
		return nil
	}

	return &AiModelHandler{
		worker: NewStageWorkerHandler(
			"ai_model",
			parser,
			securityValidator,
			dispatcher,
			process,
		),
	}
}

func (h *AiModelHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	log.Printf("ai_model worker received: %s", message.Body)
	return h.worker.HandleMessage(ctx, message)
}
