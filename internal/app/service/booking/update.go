package booking

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/logger/sl"
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TODO: сделать единую модель дляupdate и add
func (s *Service) UpdateBooking(ctx context.Context, mod *model.BookingInfo) error {
	const op = "service.booking.UpdateBooking"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.bookingRepository.CheckAvailibility(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("could not check availibility", sl.Err(errTx))
			return errTx
		}

		span.AddEvent("availibility checked")

		if !availibility.Availible && !availibility.OccupiedByClient {
			span.RecordError(ErrNotAvailible)
			span.SetStatus(codes.Error, ErrNotAvailible.Error())
			log.Error("the requested period is not vacant", sl.Err(ErrNotAvailible))
			return ErrNotAvailible
		}

		errTx = s.bookingRepository.UpdateBooking(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("the update booking operation failed", sl.Err(errTx))
			return errTx
		}

		span.AddEvent("transaction successful")

		return nil
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("transaction failed", sl.Err(err))
		if errors.As(err, pgNoConnection) {
			return ErrNoConnection
		}
		return err
	}

	return nil
}
