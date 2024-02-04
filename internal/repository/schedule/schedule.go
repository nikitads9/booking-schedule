package schedule

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"

	"github.com/gofrs/uuid"
)

type Repository interface {
	AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error)
	GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error)
	UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
}

var ErrNotFound = errors.New("there is no event with this id")
var ErrFailed = errors.New("the operation failed")

type repository struct {
	client db.Client
}

func NewScheduleRepository(client db.Client) Repository {
	return &repository{
		client: client,
	}
}
