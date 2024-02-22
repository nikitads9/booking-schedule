package event

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error)
	GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error)
	UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetVacantRooms(ctx context.Context, mod *model.Interval) ([]*model.Suite, error)
	GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error)
	CheckAvailibility(ctx context.Context, suiteID int64, start time.Time, end time.Time, userID int64) (*model.Availibility, error)
}

var (
	ErrNotFound       = errors.New("no event with this id")
	ErrQuery          = errors.New("failed to execute query")
	ErrQueryBuild     = errors.New("failed to build query")
	ErrNoRowsAffected = errors.New("no database entries affected by this operation")
	ErrParseDuration  = errors.New("failed to parse duration")
	ErrPgxScan        = errors.New("failed to read database response")
	ErrNoConnection   = errors.New("could not connect to database")
	ErrNoDates        = errors.New("no vacant dates for this room within month")
	ErrUuid           = errors.New("failed to generate uuid")
	pgNoConnection    = new(*pgconn.ConnectError)
)

type repository struct {
	client db.Client
	log    *slog.Logger
}

func NewEventRepository(client db.Client, log *slog.Logger) Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
