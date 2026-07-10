package handler

import (
	"encoding/json"
	"log"

	mediadto "go-api/infrastructure/media"
	miniodto "go-api/infrastructure/minio"
	"go-api/usecase/media"

	"github.com/gofiber/fiber/v3"
)

type MinIOHandler struct {
	createMediaUseCase       *media.CreateMediaUseCase
	generateThumbnailUseCase *media.GenerateThumbnailUseCase
}

func NewMinIOHandler(
	createMediaUseCase *media.CreateMediaUseCase,
	generateThumbnailUseCase *media.GenerateThumbnailUseCase,
) *MinIOHandler {
	return &MinIOHandler{
		createMediaUseCase:       createMediaUseCase,
		generateThumbnailUseCase: generateThumbnailUseCase,
	}
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
		userID, err := mediadto.UserIDFromKey(record.S3.Object.Key)
		if err != nil {
			log.Printf("MinIO webhook: invalid object key %q: %v", record.S3.Object.Key, err)
			continue
		}

		fileKey, err := mediadto.FileKeyFromObjectKey(record.S3.Object.Key)
		if err != nil {
			log.Printf("MinIO webhook: invalid object key %q: %v", record.S3.Object.Key, err)
			continue
		}

		media, err := h.createMediaUseCase.Execute(c.Context(), userID, fileKey, record.S3.Object.ContentType, record.S3.Object.Size)
		if err != nil {
			log.Printf("MinIO webhook: failed to create media for key %q: %v", fileKey, err)
			continue
		}

		err = h.generateThumbnailUseCase.Execute(c.Context(), userID, media.ID)
		if err != nil {
			log.Printf("MinIO webhook: failed to generate thumbnail for key %q: %v", fileKey, err)
			continue
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
