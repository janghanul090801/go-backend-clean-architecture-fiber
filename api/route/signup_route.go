package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
)

func NewSignupRouter(group fiber.Router, controller *handler.SignupHandler) {
	group.Post("/", controller.Signup)
}
