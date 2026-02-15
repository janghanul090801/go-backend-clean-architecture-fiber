package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/middleware"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

func NewTaskRouter(app fiber.Router, service domain.TaskUseCase) {

	// protected
	protected := app.Group("protected")
	protected.Use(middleware.JwtMiddleware)
	protected.Get("/", handler.FetchTask(service))
	protected.Post("/", handler.CreateTask(service))
}
