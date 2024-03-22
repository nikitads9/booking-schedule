package booking

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	AddBooking(ctx context.Context, mod *model.BookingInfo) (uuid.UUID, error)
	GetBooking(ctx context.Context, bookingID uuid.UUID, userID int64) (*model.BookingInfo, error)
	GetBookings(ctx context.Context, startDate time.Time, endDate time.Time, userID int64) ([]*model.BookingInfo, error)
	UpdateBooking(ctx context.Context, mod *model.BookingInfo) error
	DeleteBooking(ctx context.Context, bookingID uuid.UUID, userID int64) error
	GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error)
	GetBusyDates(ctx context.Context, suiteID int64) ([]*model.Interval, error)
	GetBookingListByDate(ctx context.Context, start time.Time, end time.Time) ([]*model.BookingInfo, error)
	DeleteBookingsBeforeDate(ctx context.Context, end time.Time) error
	CheckAvailibility(ctx context.Context, mod *model.BookingInfo) (*model.Availibility, error)
}

var (
	ErrNotFound       = errors.New("no booking with this id")
	ErrNoRowsAffected = errors.New("no database entries affected by this operation")
	ErrUnauthorized   = errors.New("no user associated with this token")

	ErrQuery        = errors.New("failed to execute query")
	ErrQueryBuild   = errors.New("failed to build query")
	ErrPgxScan      = errors.New("failed to read database response")
	ErrNoConnection = errors.New("could not connect to database")
	ErrUuid         = errors.New("failed to generate uuid")

	pgNoConnection = new(*pgconn.ConnectError)
	ErrNoSuchUser  = &pgconn.PgError{
		Severity:       "ERROR",
		Code:           "23503",
		Message:        "violates foreign key constraint",
		ConstraintName: "fk_users"}
)

type repository struct {
	client db.Client
	log    *slog.Logger
	tracer trace.Tracer
}

func NewBookingRepository(client db.Client, log *slog.Logger, tracer trace.Tracer) Repository {
	return &repository{
		client: client,
		log:    log,
		tracer: tracer,
	}
}
