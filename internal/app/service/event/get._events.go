package event

import (
	"context"
	"event-schedule/internal/app/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error) {
	const op = "events.service.GetEvents"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	if mod == nil {
		log.Error(ErrNoModel.Error())
		return nil, ErrNoModel
	}

	return s.eventRepository.GetEvents(ctx, mod)
}
