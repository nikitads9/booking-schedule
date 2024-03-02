package event

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	"event-schedule/internal/pkg/db"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	AddEvent(ctx context.Context, mod *model.EventInfo) (uuid.UUID, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error)
	GetEvents(ctx context.Context, startDate time.Time, endDate time.Time, userID int64) ([]*model.EventInfo, error)
	UpdateEvent(ctx context.Context, mod *model.EventInfo) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error)
	GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error)
	GetEventListByDate(ctx context.Context, start time.Time, end time.Time) ([]*model.EventInfo, error)
	DeleteEventsBeforeDate(ctx context.Context, end time.Time) error
	CheckAvailibility(ctx context.Context, mod *model.EventInfo) (*model.Availibility, error)
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
