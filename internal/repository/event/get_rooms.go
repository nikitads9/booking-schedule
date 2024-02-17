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

func (r *repository) GetVacantRooms(ctx context.Context, mod *model.Interval) ([]*model.Suite, error) {
	const op = "events.repository.GetVacantRooms"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	builder := sq.Select(t.SuiteTable+".id AS "+t.SuiteID, t.Name, t.Capacity).
		Distinct().
		From(t.SuiteTable).
		PlaceholderFormat(sq.Dollar)
	subQuery, subQueryArgs, err := sq.Select("1").
		From(t.EventTable + " AS e").
		Where(sq.And{
			sq.ConcatExpr("e."+t.SuiteID+"=", t.SuiteTable+".id"),
			sq.Or{sq.And{
				sq.Lt{"e." + t.StartDate: mod.StartDate},
				sq.Gt{"e." + t.EndDate: mod.EndDate},
			},
				sq.And{
					sq.Lt{"e." + t.StartDate: mod.EndDate},
					sq.Gt{"e." + t.EndDate: mod.StartDate},
				}},
		},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.log.Error("failed to build subquery", err)
		return nil, ErrQueryBuild
	}

	builder = builder.Where("NOT EXISTS ("+subQuery+") OR NOT EXISTS (SELECT DISTINCT "+t.SuiteID+" FROM "+t.EventTable+")", subQueryArgs...)

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.Suite
	err = r.client.DB().SelectContext(ctx, res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			r.log.Error("no vacant rooms within this period", err)
			return nil, ErrNotFound
		}
		r.log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
