package event

import (
	"context"
	"event-schedule/internal/app/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) GetVacantRooms(ctx context.Context, mod *model.Interval) ([]*model.Suite, error) {
	const op = "events.service.GetVacantRooms"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	if mod == nil {
		log.Error(ErrNoModel.Error())
		return nil, ErrNoModel
	}

	return s.eventRepository.GetVacantRooms(ctx, mod)
}
