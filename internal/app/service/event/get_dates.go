package event

import (
	"context"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/app/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	const op = "events.service.GetVacantDates"

	s.log = s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	res, err := s.eventRepository.GetVacantDates(ctx, suiteID)
	if err != nil {
		return nil, err
	}

	res = convert.ToFreeIntervals(res)
	//TODO: ошибка и проверка на ошибку
	return res, nil
}
