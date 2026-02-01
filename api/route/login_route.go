package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/handler"
)

func NewLoginRouter(group fiber.Router, controller *handler.LoginHanlder) {
	group.Post("/", controller.Login)
}
