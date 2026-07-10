package handler

import (
	"errors"
	"go-api/handler/context"
	mediadto "go-api/infrastructure/media"
	"go-api/infrastructure/paginate"
	"go-api/infrastructure/storage"
	"go-api/presenter"
	"go-api/usecase/media"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MediaHandler struct {
	storage                           *storage.MinIOStorage
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase
	getMediaUseCase                   *media.GetMediaUseCase
	getMediasUseCase                  *media.GetMediasUseCase
}

func NewMediaHandler(
	storage *storage.MinIOStorage,
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase,
	getMediaUseCase *media.GetMediaUseCase,
	getMediasUseCase *media.GetMediasUseCase,
) *MediaHandler {
	return &MediaHandler{
		storage:                           storage,
		generatePresignedUploadUrlUseCase: generatePresignedUploadUrlUseCase,
		getMediaUseCase:                   getMediaUseCase,
		getMediasUseCase:                  getMediasUseCase,
	}
}

func (h *MediaHandler) GeneratePresignedUploadUrl(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var request mediadto.PresignUploadInput
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}

	url, err := h.generatePresignedUploadUrlUseCase.Execute(c.Context(), user.ID, request)
	if err != nil {
		if errors.Is(err, media.ErrUnsupportedMediaType) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Unsupported media type",
				"errors":  err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	return c.JSON(presenter.NewGeneratePresignedUploadUrlDetailResponse(url))
}

func (h *MediaHandler) GetMedias(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"errors":  err.Error(),
		})
	}

	var query paginate.PaginateQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"errors":  err.Error(),
		})
	}
	query.Normalize()

	medias, total, err := h.getMediasUseCase.Execute(c.Context(), user.ID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"errors":  err.Error(),
		})
	}

	return c.JSON(paginate.NewPaginateResponse(presenter.NewMediaListResponses(medias), int(total), query))
}

func (h *MediaHandler) GetThumbnail(c fiber.Ctx) error {
	user, err := context.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	mediaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Printf("MediaHandler: failed to parse media ID: %v", err)
		return c.SendStatus(fiber.StatusNotFound)
	}

	media, err := h.getMediaUseCase.Execute(c.Context(), user.ID, mediaID)
	if err != nil {
		log.Printf("MediaHandler: failed to get media: %v", err)
		return c.SendStatus(fiber.StatusNotFound)
	}

	log.Printf("MediaHandler: media found: %+v", media)

	if media.Thumbnail == "" {
		log.Printf("MediaHandler: thumbnail not found for media ID: %v", mediaID)
		return c.SendStatus(fiber.StatusNotFound)
	}

	reader, err := h.storage.GetThumbnail(c.Context(), mediadto.NewThumbnailObjectKey(user.ID, media.ID))
	if err != nil {
		log.Printf("MediaHandler: failed to get thumbnail: %v", err)
		return c.SendStatus(fiber.StatusNotFound)
	}
	defer reader.Close()

	log.Printf("MediaHandler: thumbnail found for media ID: %v", mediaID)

	c.Set("Content-Type", "image/jpeg")
	c.Set("Cache-Control", "public, max-age=86400")

	return c.SendStream(reader)
}
