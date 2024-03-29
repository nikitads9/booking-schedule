package booking

import (
	"context"
	"errors"
	"log/slog"

	"booking-schedule/internal/app/model"
	t "booking-schedule/internal/app/repository/table"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TODO: проверить а что будет еси комнаты забронированы гетерогенно: разными клиентами и приходит запрос на обновление одним из них, накладывающийся на второго (умозрительно вроде норм)
func (r *repository) CheckAvailibility(ctx context.Context, mod *model.BookingInfo) (*model.Availibility, error) {
	const op = "repository.booking.CheckAvailibility"

	requestID := middleware.GetReqID(ctx)

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := r.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	subQuery := sq.Select("1").From(t.BookingTable).Where(sq.And{
		sq.And{
			sq.Or{
				sq.And{sq.Eq{t.SuiteID: mod.SuiteID},
					sq.Or{sq.Eq{t.UserID: mod.UserID}},
				},
				sq.Eq{t.ID: mod.ID},
			},
			sq.Or{
				sq.And{
					sq.GtOrEq{t.StartDate: mod.StartDate},
					sq.LtOrEq{t.StartDate: mod.EndDate},
				},
				sq.And{
					sq.GtOrEq{t.EndDate: mod.StartDate},
					sq.LtOrEq{t.EndDate: mod.EndDate},
				},
			},
		},
	}).
		Prefix("(SELECT EXISTS (").
		Suffix(")) as occupied_by_client").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sq.Select("1").From(t.BookingTable).Where(sq.And{
		sq.And{
			sq.Eq{t.SuiteID: mod.SuiteID},
			sq.Or{
				sq.And{
					sq.GtOrEq{t.StartDate: mod.StartDate},
					sq.LtOrEq{t.StartDate: mod.EndDate},
				},
				sq.And{
					sq.GtOrEq{t.EndDate: mod.StartDate},
					sq.LtOrEq{t.EndDate: mod.EndDate},
				},
			},
		},
	}).
		Prefix("SELECT NOT EXISTS (").
		Suffix(") as availible,").
		SuffixExpr(subQuery).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to build query", sl.Err(err))
		return nil, ErrQueryBuild
	}

	span.AddEvent("query built")

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res = new(model.Availibility)
	err = r.client.DB().GetContext(ctx, res, q, args...)
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
