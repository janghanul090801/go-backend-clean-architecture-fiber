package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SignupHandler struct {
	signupUsecase domain.SignupUsecase
}

func NewSignupHandler(usecase domain.SignupUsecase) *SignupHandler {
	return &SignupHandler{
		signupUsecase: usecase,
	}
}

func (h *SignupHandler) Signup(c *fiber.Ctx) error {
	ctx := c.Context()
	var request domain.SignupRequest

	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	_, err = h.signupUsecase.GetUserByEmail(ctx, request.Email)
	if err == nil {
		return c.Status(http.StatusConflict).JSON(domain.ErrorResponse{Message: "User already exists with the given email"})
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	request.Password = string(encryptedPassword)

	user := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}

	err = h.signupUsecase.Create(ctx, &user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	accessToken, err := h.signupUsecase.CreateAccessToken(&user, config.E.AccessTokenSecret, config.E.AccessTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	refreshToken, err := h.signupUsecase.CreateRefreshToken(&user, config.E.RefreshTokenSecret, config.E.RefreshTokenExpiryHour)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	signupResponse := domain.SignupResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return c.Status(http.StatusOK).JSON(signupResponse)
}
