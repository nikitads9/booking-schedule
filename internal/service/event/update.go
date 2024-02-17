package event

import (
	"context"
	"errors"
	"event-schedule/internal/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error {
	const op = "events.service.UpdateEvent"

	s.log = s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	if mod == nil {
		s.log.Error(ErrNoModel.Error())
		return ErrNoModel
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.eventRepository.CheckAvailibility(ctx, mod.SuiteID, mod.StartDate, mod.EndDate, mod.UserID)
		if errTx != nil {
			s.log.Error("could not check availibility", errTx)
			return errTx
		}

		if !availibility.Availible && !availibility.OccupiedByOwner {
			s.log.Error("the requested period is not vacant", errTx)
			return ErrNotAvailible
		}

		errTx = s.eventRepository.UpdateEvent(ctx, mod)
		if errTx != nil {
			s.log.Error("the update event operation failed", errTx)
			return errTx
		}

		return nil
	})

	if err != nil {
		s.log.Error("transaction failed", err)
		if errors.As(err, pgNoConnection) {
			return ErrNoTransaction
		}
		return err
	}

	return nil
}
