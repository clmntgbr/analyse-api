package handler

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type MinIOHandler struct{}

func NewMinIOHandler() *MinIOHandler {
	return &MinIOHandler{}
}

type minIOObjectCreatedEvent struct {
	Records []struct {
		EventName string `json:"eventName"`
		S3        struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key  string `json:"key"`
				Size int64  `json:"size"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func (h *MinIOHandler) ObjectCreated(c fiber.Ctx) error {
	payload := c.Body()

	var event minIOObjectCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("MinIO webhook: invalid payload: %v body=%s", err, string(payload))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid payload",
		})
	}

	for _, record := range event.Records {
		log.Printf(
			"MinIO webhook received: event=%s bucket=%s key=%s size=%d",
			record.EventName,
			record.S3.Bucket.Name,
			record.S3.Object.Key,
			record.S3.Object.Size,
		)
	}

	if len(event.Records) == 0 {
		log.Printf("MinIO webhook received with no records: body=%s", strings.TrimSpace(string(payload)))
	}

	return c.SendStatus(fiber.StatusOK)
}
