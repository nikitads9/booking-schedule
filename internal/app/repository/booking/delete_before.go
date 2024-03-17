package booking

import (
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/codes"
)

func (r *repository) DeleteBookingsBeforeDate(ctx context.Context, date time.Time) error {
	const op = "repository.booking.DeleteBookingsBeforeDate"

	log := r.log.With(
		slog.String("op", op),
	)
	ctx, span := r.tracer.Start(ctx, op)
	defer span.End()

	builder := sq.Delete(t.BookingTable).
		Where(sq.Lt{t.EndDate: date}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to build a query", sl.Err(err))
		return ErrQueryBuild
	}

	span.AddEvent("query built")

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", sl.Err(err))
			return ErrNoConnection
		}
		log.Error("query execution error", sl.Err(err))
		return ErrQuery
	}

	span.AddEvent("query successfully executed")

	return nil
}
