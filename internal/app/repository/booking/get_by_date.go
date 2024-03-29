package booking

import (
	"booking-schedule/internal/app/model"
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

func (r *repository) GetBookingListByDate(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.BookingInfo, error) {
	op := "repository.booking.GetBookingListByDate"

	log := r.log.With(slog.String("op", op))

	ctx, span := r.tracer.Start(ctx, op)
	defer span.End()

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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to build a query", sl.Err(err))
		return nil, ErrQueryBuild
	}

	span.AddEvent("query built")

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.BookingInfo
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", sl.Err(err))
			return nil, ErrNoConnection
		}
		log.Error("query execution error", sl.Err(err))
		return nil, ErrQuery
	}

	span.AddEvent("query successfully executed")

	return res, nil
}
