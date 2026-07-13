package handler

import (
	"encoding/json"
	"log"
	"strings"

	"go-api/domain/enum"
	mediadto "go-api/infrastructure/media"
	miniodto "go-api/infrastructure/minio"
	"go-api/usecase/media"

	"github.com/gofiber/fiber/v3"
)

type MinIOHandler struct {
	mediaBucket              string
	createMediaUseCase       *media.CreateMediaUseCase
	generateThumbnailUseCase *media.GenerateThumbnailUseCase
	updateMediaStatusUseCase *media.UpdateMediaStatusUseCase
	findMediaMetadataUseCase *media.FindMediaMetadataUseCase
}

func NewMinIOHandler(
	mediaBucket string,
	createMediaUseCase *media.CreateMediaUseCase,
	generateThumbnailUseCase *media.GenerateThumbnailUseCase,
	updateMediaStatusUseCase *media.UpdateMediaStatusUseCase,
	findMediaMetadataUseCase *media.FindMediaMetadataUseCase,
) *MinIOHandler {
	return &MinIOHandler{
		mediaBucket:              mediaBucket,
		createMediaUseCase:       createMediaUseCase,
		generateThumbnailUseCase: generateThumbnailUseCase,
		updateMediaStatusUseCase: updateMediaStatusUseCase,
		findMediaMetadataUseCase: findMediaMetadataUseCase,
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
		if record.S3.Bucket.Name != h.mediaBucket {
			continue
		}

		decodedKey, err := mediadto.DecodeObjectKey(record.S3.Object.Key)
		if err != nil {
			log.Printf("MinIO webhook: invalid object key %q: %v", record.S3.Object.Key, err)
			continue
		}

		if strings.Contains(decodedKey, "/thumbnails/") {
			continue
		}

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

		createdMedia, err := h.createMediaUseCase.Execute(c.Context(), userID, fileKey, record.S3.Object.ContentType, record.S3.Object.Size)
		if err != nil {
			log.Printf("MinIO webhook: failed to create media for key %q: %v", fileKey, err)
			continue
		}

		err = h.generateThumbnailUseCase.Execute(c.Context(), userID, createdMedia.ID)
		if err != nil {
			log.Printf("MinIO webhook: failed to generate thumbnail for key %q: %v", fileKey, err)
			continue
		}

		err = h.updateMediaStatusUseCase.Execute(c.Context(), createdMedia.ID, enum.MediaStatusUploaded)
		if err != nil {
			log.Printf("MinIO webhook: failed to update media status for key %q: %v", fileKey, err)
			continue
		}

		err = h.findMediaMetadataUseCase.Execute(c.Context(), createdMedia.ID)
		if err != nil {
			log.Printf("MinIO webhook: failed to find media metadata for key %q: %v", fileKey, err)
			continue
		}
	}

	return c.SendStatus(fiber.StatusOK)
}
