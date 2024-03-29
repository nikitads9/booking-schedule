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
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (r *repository) GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error) {
	const op = "repository.booking.GetVacantRooms"

	requestID := middleware.GetReqID(ctx)

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := r.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	builder := sq.Select(t.SuiteTable+".id AS "+t.SuiteID, t.Name, t.Capacity).
		Distinct().
		From(t.SuiteTable).
		PlaceholderFormat(sq.Dollar)
	subQuery, subQueryArgs, err := sq.Select("1").
		From(t.BookingTable + " AS e").
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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to build subquery", sl.Err(err))
		return nil, ErrQueryBuild
	}

	builder = builder.Where("NOT EXISTS ("+subQuery+") OR NOT EXISTS (SELECT DISTINCT "+t.SuiteID+" FROM "+t.BookingTable+")", subQueryArgs...)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build a query", sl.Err(err))
		return nil, ErrQueryBuild
	}

	span.AddEvent("query built")

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.Suite
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", sl.Err(err))
			return nil, ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			log.Error("no vacant rooms within this period", sl.Err(err))
			return nil, ErrNotFound
		}
		log.Error("query execution error", sl.Err(err))
		return nil, ErrQuery
	}

	span.AddEvent("query successfully executed")

	return res, nil
}
