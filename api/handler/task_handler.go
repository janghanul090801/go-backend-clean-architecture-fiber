package handler

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

func CreateTask(service domain.TaskUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()
		var task domain.Task
		var errInfo domain.Error

		err := c.Bind().Body(&task)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
		}

		userID := c.Locals("id").(*domain.ID)

		_, err = service.Create(ctx, &task, userID)
		if err != nil {
			if ok := errors.As(err, &errInfo); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		return c.Status(http.StatusOK).JSON(domain.SuccessResponse{
			Message: "Task created successfully",
		})
	}
}

func FetchTask(service domain.TaskUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.RequestCtx()

		var errInfo domain.Error

		userID := c.Locals("id").(*domain.ID)

		tasks, err := service.FetchByUserID(ctx, userID)
		if err != nil {
			if ok := errors.As(err, &errInfo); ok {
				return c.Status(errInfo.StatusCode).JSON(domain.ErrorResponse{Message: err.Error()})
			}
		}

		return c.Status(http.StatusOK).JSON(tasks)
	}
}
