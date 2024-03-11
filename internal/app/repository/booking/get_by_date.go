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
)

func (r *repository) GetBookingListByDate(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.BookingInfo, error) {
	op := "bookings.repository.GetBookingListByDate"
	log := r.log.With(slog.String("op", op))

	builder := sq.Select(t.ID, t.SuiteID, t.StartDate, t.EndDate, t.NotifyAt, t.CreatedAt, t.UpdatedAt, t.UserID).
		From(t.BookingTable).
		Where(sq.Or{
			sq.And{
				sq.Gt{t.StartDate: startDate},
				sq.LtOrEq{t.StartDate: endDate},
			},
			sq.And{
				sq.Gt{t.StartDate + "-" + t.NotifyAt: startDate},
				sq.LtOrEq{t.StartDate + "-" + t.NotifyAt: endDate},
			},
		}).PlaceholderFormat(sq.Dollar)

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
		log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
