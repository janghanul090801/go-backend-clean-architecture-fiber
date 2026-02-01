package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/middleware"
)

func NewRefreshTokenRouter(group fiber.Router, controller *handler.RefreshTokenHandler) {

	// protected
	protected := group.Group("protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Post("/", controller.RefreshToken)
}
