package handler

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HeuristicHandler struct{}

func NewHeuristicHandler() *HeuristicHandler {
	return &HeuristicHandler{}
}

func (h *HeuristicHandler) HandleMessage(_ context.Context, message *amqp.Delivery) error {
	log.Printf("heuristic message received: %s", message.Body)
	return nil
}
