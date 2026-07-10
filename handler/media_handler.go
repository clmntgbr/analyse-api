package handler

import (
	mediadto "go-api/infrastructure/media"
	"go-api/handler/context"
	"go-api/presenter"
	"go-api/usecase/media"

	"github.com/gofiber/fiber/v3"
)

type MediaHandler struct {
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase
}

func NewMediaHandler(
	generatePresignedUploadUrlUseCase *media.GeneratePresignedUploadUrlUseCase,
) *MediaHandler {
	return &MediaHandler{
		generatePresignedUploadUrlUseCase: generatePresignedUploadUrlUseCase,
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

	return c.JSON(presenter.NewMediaDetailResponse(url))
}
