package scheduler

import (
	"context"
	"encoding/json"
	"event-schedule/internal/app/model"
	"log/slog"
	"time"
)

func (s *Service) Run(ctx context.Context) {
	const op = "scheduler.service.Run"

	log := s.log.With(
		slog.String("op", op),
	)
	log.Info("scheduler initiated")

	ticker := time.NewTicker(s.checkPeriod)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.handleEvents(ctx)
			if err != nil {
				log.Error("failed to handle events:", err)
			}
		}
	}

}

func (s *Service) handleEvents(ctx context.Context) error {
	const op = "scheduler.service.handleEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("started handling")

	events, err := s.getEvents(ctx)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		log.Debug("No events.")
		return nil
	}

	err = s.sendEvent(events)
	if err != nil {
		return err
	}

	err = s.cleanUpOldEvents(ctx)
	if err != nil {
		return err
	}

	log.Debug("successfully handled events")

	return nil
}

func (s *Service) getEvents(ctx context.Context) ([]*model.EventInfo, error) {
	const op = "scheduler.service.getEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	end := time.Now()
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), 0, 0, end.Location())
	start := end.Add(-s.checkPeriod)

	events, err := s.eventRepository.GetEventListByDate(ctx, start, end)
	if err != nil {
		log.Error("failed to get list by date", err)
		return nil, err
	}

	return events, nil
}

func (s *Service) cleanUpOldEvents(ctx context.Context) error {
	const op = "scheduler.service.cleanUpOldEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.eventRepository.DeleteEventsBeforeDate(ctx, time.Now().Add(-s.eventTTL))
	if err != nil {
		log.Error("failed to clean up old events", err)
		return err
	}

	return nil
}

func (s *Service) sendEvent(events []*model.EventInfo) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}

	return s.rabbitProducer.Publish(data)
}
