package booking

import (
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) DeleteBookingsBeforeDate(ctx context.Context, date time.Time) error {
	const op = "bookings.repository.DeleteBookingsBeforeDate"

	log := r.log.With(
		slog.String("op", op),
	)

	builder := sq.Delete(t.BookingTable).
		Where(sq.Lt{t.EndDate: date}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		log.Error("query execution error", err)
		return ErrQuery
	}

	return nil
}
