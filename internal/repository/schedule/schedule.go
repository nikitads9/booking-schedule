package schedule

import (
	"context"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
)

type Repository interface {
	AddEvent(ctx context.Context) (string, error)
	GetEvents(ctx context.Context, userID int64) ([]*model.Event, error)
	GetEvent(ctx context.Context, eventID string) (*model.Event, error)
	UpdateEvent(ctx context.Context)
	DeleteEvent(ctx context.Context)
}
type repository struct {
	client db.Client
}

func NewScheduleRepository(client db.Client) Repository {
	return &repository{
		client: client,
	}
}
