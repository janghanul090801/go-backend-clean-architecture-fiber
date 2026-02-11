package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

type TaskHandler struct {
	taskUsecase domain.TaskUsecase
}

func NewTaskHandler(usecase domain.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: usecase,
	}
}

func (h *TaskHandler) Create(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	var task domain.Task

	err := c.Bind().Body(&task)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	userID := c.Locals("id").(*domain.ID)

	println("live")

	task.UserID = *userID

	_, err = h.taskUsecase.Create(ctx, &task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	println("live")

	return c.Status(http.StatusOK).JSON(domain.SuccessResponse{
		Message: "Task created successfully",
	})
}

func (h *TaskHandler) Fetch(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Locals("id").(*domain.ID)

	tasks, err := h.taskUsecase.FetchByUserID(ctx, userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Message: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(tasks)
}
