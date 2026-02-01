package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"net/http"
)

type RefreshTokenController struct {
	refreshTokenUsecase domain.RefreshTokenUsecase
}

func NewRefreshTokenController(usecase domain.RefreshTokenUsecase) *RefreshTokenController {
	return &RefreshTokenController{
		refreshTokenUsecase: usecase,
	}
}

func (rtc *RefreshTokenController) RefreshToken(c *fiber.Ctx) error {
	ctx := c.Context()

	var request domain.RefreshTokenRequest

	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	id, err := rtc.refreshTokenUsecase.ExtractIDFromToken(request.RefreshToken, config.E.RefreshTokenSecret)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
	}

	user, err := rtc.refreshTokenUsecase.GetUserByID(ctx, id)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
	}

	accessToken, err := rtc.refreshTokenUsecase.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshToken, err := rtc.refreshTokenUsecase.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(http.StatusOK).JSON(refreshTokenResponse)
}
