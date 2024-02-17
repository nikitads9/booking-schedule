package event

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	t "event-schedule/internal/repository/table"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error) {
	const op = "events.repository.GetEvents"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Select("id", t.SuiteID, t.StartDate, t.EndDate, t.NotifyAt, t.CreatedAt, t.UpdatedAt).
		From(t.EventTable).
		Where(sq.And{
			sq.Eq{t.OwnerID: mod.UserID},
			sq.GtOrEq{t.StartDate: mod.StartDate},
			sq.LtOrEq{t.EndDate: mod.EndDate},
		}).
		PlaceholderFormat(sq.Dollar)

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
		if pgxscan.NotFound(err) {
			r.log.Error("events associated with this user not found within this period", err)
			return nil, ErrNotFound
		}
		r.log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
