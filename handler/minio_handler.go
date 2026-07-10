package handler

import (
	"encoding/json"
	"log"

	miniodto "go-api/infrastructure/minio"

	"github.com/gofiber/fiber/v3"
)

type MinIOHandler struct{}

func NewMinIOHandler() *MinIOHandler {
	return &MinIOHandler{}
}

func (h *MinIOHandler) ObjectCreated(c fiber.Ctx) error {
	payload := c.Body()

	var event miniodto.ObjectCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("MinIO webhook: invalid payload: %v body=%s", err, string(payload))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	for _, record := range event.Records {
		log.Printf(
			"MinIO webhook received: key=%s",
			record.S3.Object.Key,
		)
	}

	return c.SendStatus(fiber.StatusOK)
}
