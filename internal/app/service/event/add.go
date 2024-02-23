package event

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"github.com/gofrs/uuid"
)

func (s *Service) AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error) {
	const op = "events.service.AddEvent"

	s.log = s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	if mod == nil {
		s.log.Error(ErrNoModel.Error())
		return uuid.Nil, ErrNoModel
	}

	var id uuid.UUID

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		availibility, errTx := s.eventRepository.CheckAvailibility(ctx, mod)
		if errTx != nil {
			s.log.Error("could not check availibility", errTx)
			return errTx
		}

		if !availibility.Availible {
			s.log.Error("the requested period is not vacant", errTx)
			return ErrNotAvailible
		}

		id, errTx = s.eventRepository.AddEvent(ctx, mod)
		if errTx != nil {
			s.log.Error("the add event operation failed", errTx)
			return errTx
		}

		return nil
	})

	if err != nil {
		s.log.Error("transaction failed", err)
		if errors.As(err, pgNoConnection) {
			return uuid.Nil, ErrNoTransaction
		}
		return uuid.Nil, err
	}

	return id, nil
}
