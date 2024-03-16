package booking

import (
	"booking-schedule/internal/app/model"
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/codes"
)

func (r *repository) GetBooking(ctx context.Context, bookingID uuid.UUID, userID int64) (*model.BookingInfo, error) {
	const op = "repository.booking.GetBooking"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
	ctx, span := r.tracer.Start(ctx, op)
	defer span.End()

	builder := sq.Select(t.ID, t.SuiteID, t.StartDate, t.EndDate, t.NotifyAt, t.CreatedAt, t.UpdatedAt, t.UserID).
		From(t.BookingTable).
		Where(sq.And{
			sq.Eq{t.ID: bookingID},
			sq.Eq{t.UserID: userID},
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	span.AddEvent("query built")

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res = new(model.BookingInfo)
	err = r.client.DB().GetContext(ctx, res, q, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("booking with this id not found", err)
			return nil, ErrNotFound
		}
		log.Error("query execution error", err)
		return nil, ErrQuery
	}

	span.AddEvent("query successfully executed")

	return res, nil
}
