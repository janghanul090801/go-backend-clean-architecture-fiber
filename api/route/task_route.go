package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/middleware"
)

func NewTaskRouter(group fiber.Router, controller *handler.TaskHandler) {

	// protected
	protected := group.Group("protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Get("/", controller.Fetch)
	protected.Post("/", controller.Create)
}
