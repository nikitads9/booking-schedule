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
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) GetBookings(ctx context.Context, startDate time.Time, endDate time.Time, userID int64) ([]*model.BookingInfo, error) {
	const op = "bookings.repository.GetBookings"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Select(t.ID, t.SuiteID, t.StartDate, t.EndDate, t.NotifyAt, t.CreatedAt, t.UpdatedAt, t.UserID).
		From(t.BookingTable).
		Where(sq.And{
			sq.Eq{t.UserID: userID},
			sq.Or{
				sq.And{
					sq.GtOrEq{t.StartDate: startDate},
					sq.LtOrEq{t.StartDate: endDate},
				},
				sq.And{
					sq.GtOrEq{t.EndDate: startDate},
					sq.LtOrEq{t.EndDate: endDate},
				},
			},
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.BookingInfo
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			log.Error("bookings associated with this user not found within this period", err)
			return nil, ErrNotFound
		}
		log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
