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
	"github.com/go-chi/chi/middleware"
)

/*
SELECT suite_id, start_date, end_date
FROM events
WHERE suite_id = YOUR_SUITE_ID
AND NOT (

	(end_date <= '2024-02-22T17:43:00-03:00') OR
	(start_date >= '2024-03-22T17:43:00-03:00')

)
ORDER BY start_date;
*/
func (r *repository) GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	const op = "events.repository.GetVacantDates"

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	now := "'" + time.Now().Format("01-02-2006") + "'"
	month := "'" + time.Now().Add(720*time.Hour).Format("01-02-2006") + "'"

	builder := sq.Select(t.StartDate+` as start`, t.EndDate+` as end`).
		From(t.EventTable).
		Where(sq.And{
			sq.Eq{t.SuiteID: suiteID},
			sq.And{
				sq.Gt{t.EndDate: now},
				sq.Lt{t.StartDate: month},
			},
		}).
		OrderBy(t.StartDate).
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

	var res []*model.Interval
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
