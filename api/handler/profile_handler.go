package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

type ProfileHandler struct {
	profileUsecase domain.ProfileUsecase
}

func NewProfileHandler(usecase domain.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{
		profileUsecase: usecase,
	}
}

func (h *ProfileHandler) Fetch(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Locals("id").(*domain.ID)

	profile, err := h.profileUsecase.GetProfileByID(ctx, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(profile)
}
