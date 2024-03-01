package event

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

// TODO: сделать единую модель дляupdate и add
func (s *Service) UpdateEvent(ctx context.Context, mod *model.Event) error {
	const op = "events.service.UpdateEvent"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	if mod == nil {
		log.Error(ErrNoModel.Error())
		return ErrNoModel
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.eventRepository.CheckAvailibility(ctx, mod)
		if errTx != nil {
			log.Error("could not check availibility", errTx)
			return errTx
		}

		if !availibility.Availible && !availibility.OccupiedByClient {
			log.Error("the requested period is not vacant", errTx)
			return ErrNotAvailible
		}

		errTx = s.eventRepository.UpdateEvent(ctx, mod)
		if errTx != nil {
			log.Error("the update event operation failed", errTx)
			return errTx
		}

		return nil
	})

	if err != nil {
		log.Error("transaction failed", err)
		if errors.As(err, pgNoConnection) {
			return ErrNoTransaction
		}
		return err
	}

	return nil
}
