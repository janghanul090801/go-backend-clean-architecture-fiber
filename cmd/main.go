package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/controller"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/api/route"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/bootstrap"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/config"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/infra/repository"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/usecase"
	"log"
	"net/http"
	"time"
)

func main() {

	app := bootstrap.App()

	api := app.App.Group("/api")

	client := app.Client
	defer app.CloseDBConnection()

	timeout := time.Duration(config.E.ContextTimeout) * time.Second

	// repository
	userRepository := repository.NewUserRepository(client)
	taskRepository := repository.NewTaskRepository(client)

	// usecase
	loginUsecase := usecase.NewLoginUsecase(userRepository, timeout)
	profileUsecase := usecase.NewProfileUsecase(userRepository, timeout)
	refreshTokenUsecase := usecase.NewRefreshTokenUsecase(userRepository, timeout)
	signupUsecase := usecase.NewSignupUsecase(userRepository, timeout)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, timeout)

	// controller
	loginController := controller.NewLoginController(loginUsecase)
	profileController := controller.NewProfileController(profileUsecase)
	refreshTokenController := controller.NewRefreshTokenController(refreshTokenUsecase)
	signupController := controller.NewSignupController(signupUsecase)
	taskController := controller.NewTaskController(taskUsecase)

	// router
	route.NewLoginRouter(api.Group("/login"), loginController)
	route.NewProfileRouter(api.Group("/profile"), profileController)
	route.NewRefreshTokenRouter(api.Group("/refresh"), refreshTokenController)
	route.NewSignupRouter(api.Group("/signup"), signupController)
	route.NewTaskRouter(api.Group("/task"), taskController)

	app.App.All("*", func(c *fiber.Ctx) error {
		notFoundErr := fmt.Errorf(
			"route '%s' does not exist in this API",
			c.OriginalURL(),
		)

		return c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"status": "error",
			"error":  notFoundErr.Error(),
		})
	})

	log.Fatal(app.App.Listen(config.E.ServerAddress))
}
