package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
)

func NewRefreshTokenRouter(group fiber.Router, controller *handler.RefreshTokenHandler) {

	// protected
	group.Post("/", controller.RefreshToken)
}
