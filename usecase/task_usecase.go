package usecase

import (
	"context"
	"time"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
)

type taskUseCase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUseCase(taskRepository domain.TaskRepository, timeout time.Duration) domain.TaskUseCase {
	return &taskUseCase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

func (tu *taskUseCase) Create(c context.Context, task *domain.Task) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.Create(ctx, task)
}

func (tu *taskUseCase) FetchByUserID(c context.Context, userID *domain.ID) ([]*domain.Task, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()
	return tu.taskRepository.FetchByUserID(ctx, userID)
}
