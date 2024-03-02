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
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
)

func (r *repository) GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error) {
	const op = "events.repository.GetVacantRooms"

	log := r.log.With(
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
				sq.Lt{"e." + t.StartDate: startDate},
				sq.Gt{"e." + t.EndDate: endDate},
			},
				sq.And{
					sq.Lt{"e." + t.StartDate: endDate},
					sq.Gt{"e." + t.EndDate: startDate},
				}},
		},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("failed to build subquery", err)
		return nil, ErrQueryBuild
	}

	builder = builder.Where("NOT EXISTS ("+subQuery+") OR NOT EXISTS (SELECT DISTINCT "+t.SuiteID+" FROM "+t.EventTable+")", subQueryArgs...)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.Suite
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			log.Error("no vacant rooms within this period", err)
			return nil, ErrNotFound
		}
		log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
