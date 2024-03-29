package user

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"

	t "booking-schedule/internal/app/repository/table"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetUserByNickname(ctx context.Context, nickName string) (*model.User, error) {
	const op = "users.repository.GetUser"

	requestID := middleware.GetReqID(ctx)

	log := r.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := r.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	builder := sq.Select("*").
		From(t.UserTable).
		Where(sq.Eq{t.TelegramNickname: nickName}).
		PlaceholderFormat(sq.Dollar)

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

	var res = new(model.User)
	err = r.client.DB().GetContext(ctx, res, q, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if errors.As(err, pgNoConnection) {
			log.Error("no connection to database host", sl.Err(err))
			return nil, ErrNoConnection
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("user with this id not found", sl.Err(err))
			return nil, ErrNotFound
		}
		log.Error("query execution error", sl.Err(err))
		return nil, ErrQuery
	}

	span.AddEvent("query successfully executed and response scanned")

	return res, nil
}
