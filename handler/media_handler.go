package handler

import (
	"go-api/handler/context"
	mediadto "go-api/infrastructure/media"
	"go-api/infrastructure/paginate"
	"go-api/presenter"
	"go-api/usecase/media"

	"github.com/gofiber/fiber/v3"
)

type MediaHandler struct {
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase
	getMediasUseCase                  *media.GetMediasUseCase
}

func NewMediaHandler(
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase,
	getMediasUseCase *media.GetMediasUseCase,
) *MediaHandler {
	return &MediaHandler{
		generatePresignedUploadUrlUseCase: generatePresignedUploadUrlUseCase,
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
