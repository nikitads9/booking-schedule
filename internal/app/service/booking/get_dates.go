package booking

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"context"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) GetVacantDates(ctx context.Context, suiteID int64) ([]*api.Interval, error) {
	const op = "service.booking.GetVacantDates"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	dates, err := s.bookingRepository.GetBusyDates(ctx, suiteID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("could not get busy dates", sl.Err(err))
		return nil, err
	}

	span.AddEvent("acquired busy dates", trace.WithAttributes(attribute.Int("quantity", len(dates))))

	vacant := convert.ToVacantDates(dates)

	span.AddEvent("converted to vacant dates", trace.WithAttributes(attribute.Int("quantity", len(vacant))))
	log.Info("converted to vacant dates", slog.Int("quantity: ", len(vacant)))

	return vacant, nil
}
