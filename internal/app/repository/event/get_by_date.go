package event

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	t "event-schedule/internal/app/repository/table"
	"event-schedule/internal/pkg/db"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) GetEventListByDate(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.EventInfo, error) {
	op := "events.repository.GetEventListByDate"
	r.log = r.log.With(slog.String("op", op))

	builder := sq.Select(t.ID, t.SuiteID, t.StartDate, t.EndDate, t.NotifyAt, t.OwnerID).
		From(t.EventTable).
		Where(sq.Or{
			sq.And{
				sq.Gt{t.StartDate: startDate},
				sq.LtOrEq{t.StartDate: endDate},
			},
			sq.And{
				//TODO: maybe brackets are required
				sq.Gt{t.StartDate + "-" + t.NotifyAt: startDate},
				sq.LtOrEq{t.StartDate + "-" + t.NotifyAt: endDate},
			},
		}).PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.EventInfo
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		r.log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
