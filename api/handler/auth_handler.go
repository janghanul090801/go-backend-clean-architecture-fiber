package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"golang.org/x/crypto/bcrypt"
)

func Login(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var request domain.LoginRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		user, err := service.GetUserByEmail(ctx, request.Email)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(domain.ErrorResponse{Message: "User not found with the given email"})
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
			return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "Invalid credentials"})
		}

		accessToken, err := service.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		refreshToken, err := service.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}

func RefreshToken(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var request domain.RefreshTokenRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		id, err := service.ExtractIDFromRefreshToken(request.RefreshToken, config.E.RefreshTokenSecret)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
		}

		user, err := service.GetUserByID(ctx, id)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Message: "User not found"})
		}

		accessToken, err := service.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		refreshToken, err := service.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}

func Signup(service domain.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		var request domain.SignupRequest

		err := c.Bind().Body(&request)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		_, err = service.GetUserByEmail(ctx, request.Email)
		if err == nil {
			return c.Status(http.StatusConflict).JSON(domain.ErrorResponse{Message: "User already exists with the given email"})
		}

		user, err := service.Create(ctx, request.Name, request.Email, request.Password)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		accessToken, err := service.CreateAccessToken(user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		refreshToken, err := service.CreateRefreshToken(user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		response := domain.AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		return c.Status(http.StatusOK).JSON(response)
	}
}
