package repository

import (
	"context"

	"github.com/janghanul090801/go-backend-clean-architecture-fiber/domain"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/task"
	"github.com/janghanul090801/go-backend-clean-architecture-fiber/ent/user"
)

type taskRepository struct {
	client *ent.Client
}

func NewTaskRepository(client *ent.Client) domain.TaskRepository {
	return &taskRepository{
		client: client,
	}
}

func (r *taskRepository) Create(c context.Context, task *domain.Task) (*domain.Task, error) {

	t, err := r.client.Task.Create().
		SetTitle(task.Title).
		SetOwnerID(task.UserID).
		Save(c)

	if err != nil {
		return nil, err
	}

	return &domain.Task{
		ID:        t.ID,
		Title:     t.Title,
		CreatedAt: t.CreatedAt,
	}, nil
}

func (r *taskRepository) FetchByUserID(c context.Context, userID *domain.ID) ([]*domain.Task, error) {
	t, err := r.client.Task.
		Query().
		Where(
			task.HasOwnerWith(
				user.IDEQ(
					*userID,
				),
			),
		).
		WithOwner().
		Order(
			ent.Desc(
				task.FieldCreatedAt,
			),
		).
		All(c)

	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(t))
	for i, te := range t {
		tasks[i] = toDomainTask(te)
	}

	return tasks, err
}
