package booking

import (
	"booking-schedule/internal/app/model"
	"context"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) GetBusyDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	const op = "bookings.service.GetVacantDates"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	res, err := s.bookingRepository.GetBusyDates(ctx, suiteID)
	if err != nil {
		log.Error("could not get vacant dates:", err)
		return nil, err
	}

	return res, nil
}
