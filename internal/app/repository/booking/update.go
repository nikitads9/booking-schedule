package booking

import (
	"booking-schedule/internal/app/model"
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) UpdateBooking(ctx context.Context, mod *model.BookingInfo) error {
	const op = "bookings.repository.UpdateBooking"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Update(t.BookingTable).
		Set(t.UpdatedAt, time.Now()).
		Set("start_date", mod.StartDate).
		Set("end_date", mod.EndDate).
		Set("suite_id", mod.SuiteID).
		Where(sq.And{
			sq.Eq{t.ID: mod.ID},
			sq.Eq{t.UserID: mod.UserID},
		}).
		PlaceholderFormat(sq.Dollar)

	if mod.NotifyAt != 0 {
		builder = builder.Set("notify_at", mod.NotifyAt)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	result, err := r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return ErrNoConnection
		}
		log.Error("query execution error", err)
		return ErrQuery
	}

	if result.RowsAffected() == 0 {
		log.Error("unsuccessful update", ErrNoRowsAffected)
		return ErrNotFound
	}

	return nil
}
