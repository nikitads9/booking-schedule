package schedule

import (
	"context"
	"event-schedule/internal/client/db"
)

type Repository interface {
	AddEvent(ctx context.Context) (string, error)
	GetEvents(ctx context.Context, userID string) (string, error)
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
