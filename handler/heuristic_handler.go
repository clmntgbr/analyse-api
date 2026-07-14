package handler

import (
	"context"

	rabbitmqDTO "go-api/infrastructure/messaging/rabbitmq"
	"go-api/infrastructure/messaging/security"
	heuristicuc "go-api/usecase/heuristic"
	pipelineuc "go-api/usecase/pipeline"
	"go-api/usecase/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HeuristicHandler struct {
	parser                        *security.WorkerParser
	securityValidator             *security.WorkerSecurityValidator
	dispatcher                    *pipelineuc.Dispatcher
	analyzeMediaHeuristicsUseCase *heuristicuc.AnalyzeMediaHeuristicsUseCase
	createSignalUseCase           *signal.CreateSignalUseCase
}

func NewHeuristicHandler(
	parser *security.WorkerParser,
	securityValidator *security.WorkerSecurityValidator,
	dispatcher *pipelineuc.Dispatcher,
	analyzeMediaHeuristicsUseCase *heuristicuc.AnalyzeMediaHeuristicsUseCase,
	createSignalUseCase *signal.CreateSignalUseCase,
) *HeuristicHandler {
	return &HeuristicHandler{
		parser:                        parser,
		securityValidator:             securityValidator,
		dispatcher:                    dispatcher,
		analyzeMediaHeuristicsUseCase: analyzeMediaHeuristicsUseCase,
		createSignalUseCase:           createSignalUseCase,
	}
}

func (h *HeuristicHandler) process(ctx context.Context, message rabbitmqDTO.AnalyzeMessage) error {
	result, err := h.analyzeMediaHeuristicsUseCase.Execute(ctx, message.UserID, message.MediaKey)
	if err != nil {
		return err
	}

	_, err = h.createSignalUseCase.Execute(
		ctx,
		message.MediaID,
		"heuristics",
		result.Signal.Score,
		result.Signal.Confidence,
		result.Signal.Details,
	)
	return err
}

func (h *HeuristicHandler) HandleMessage(ctx context.Context, message *amqp.Delivery) error {
	worker := NewStageWorkerHandler(
		"heuristics",
		h.parser,
		h.securityValidator,
		h.dispatcher,
		h.process,
	)

	return worker.HandleMessage(ctx, message)
}
