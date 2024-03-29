package booking

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/logger/sl"
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) AddBooking(ctx context.Context, mod *model.BookingInfo) (uuid.UUID, error) {
	const op = "service.bookings.AddBooking"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	var id uuid.UUID

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.bookingRepository.CheckAvailibility(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("could not check availibility", sl.Err(errTx))
			return errTx
		}
		span.AddEvent("availibility checked")

		if !availibility.Availible {
			span.RecordError(ErrNotAvailible)
			span.SetStatus(codes.Error, ErrNotAvailible.Error())
			log.Error("the requested period is not vacant", sl.Err(ErrNotAvailible))
			return ErrNotAvailible
		}

		id, errTx = s.bookingRepository.AddBooking(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("the add booking operation failed", sl.Err(errTx))
			return errTx
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("transaction failed", sl.Err(err))
		if errors.As(err, pgNoConnection) {
			return uuid.Nil, ErrNoConnection
		}
		return uuid.Nil, err
	}

	span.AddEvent("transaction successful")

	return id, nil
}
