package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"

	"golang.org/x/crypto/bcrypt"
)

type LoginHanlder struct {
	loginUsecase domain.LoginUsecase
}

func NewLoginHandler(usecase domain.LoginUsecase) *LoginHanlder {
	return &LoginHanlder{
		loginUsecase: usecase,
	}
}

func (h *LoginHanlder) Login(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var request domain.LoginRequest

	err := c.Bind().Body(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	user, err := h.loginUsecase.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(domain.ErrorResponse{Message: "User not found with the given email"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "Invalid credentials"})
	}

	accessToken, err := h.loginUsecase.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshToken, err := h.loginUsecase.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(http.StatusOK).JSON(loginResponse)
}
