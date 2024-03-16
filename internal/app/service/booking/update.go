package booking

import (
	"booking-schedule/internal/app/model"
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/codes"
)

// TODO: сделать единую модель дляupdate и add
func (s *Service) UpdateBooking(ctx context.Context, mod *model.BookingInfo) error {
	const op = "service.booking.UpdateBooking"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.bookingRepository.CheckAvailibility(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("could not check availibility", errTx)
			return errTx
		}

		span.AddEvent("availibility checked")

		if !availibility.Availible && !availibility.OccupiedByClient {
			span.RecordError(ErrNotAvailible)
			span.SetStatus(codes.Error, ErrNotAvailible.Error())
			log.Error("the requested period is not vacant", ErrNotAvailible)
			return ErrNotAvailible
		}

		errTx = s.bookingRepository.UpdateBooking(ctx, mod)
		if errTx != nil {
			span.RecordError(errTx)
			span.SetStatus(codes.Error, errTx.Error())
			log.Error("the update booking operation failed", errTx)
			return errTx
		}

		span.AddEvent("transaction successful")

		return nil
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("transaction failed", err)
		if errors.As(err, pgNoConnection) {
			return ErrNoConnection
		}
		return err
	}

	return nil
}
