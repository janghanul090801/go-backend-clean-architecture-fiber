package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

type RefreshTokenHandler struct {
	refreshTokenUsecase domain.RefreshTokenUsecase
}

func NewRefreshTokenHandler(usecase domain.RefreshTokenUsecase) *RefreshTokenHandler {
	return &RefreshTokenHandler{
		refreshTokenUsecase: usecase,
	}
}

func (h *RefreshTokenHandler) RefreshToken(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var request domain.RefreshTokenRequest

	err := c.Bind().Body(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	id, err := h.refreshTokenUsecase.ExtractIDFromToken(request.RefreshToken, config.E.RefreshTokenSecret)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
	}

	user, err := h.refreshTokenUsecase.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
	}

	accessToken, err := h.refreshTokenUsecase.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshToken, err := h.refreshTokenUsecase.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(http.StatusOK).JSON(refreshTokenResponse)
}
