package scheduler

import (
	"event-schedule/internal/app/repository/event"
	"event-schedule/internal/pkg/rabbit"
	"log/slog"
	"time"
)

type Service struct {
	eventRepository event.Repository
	log             *slog.Logger
	rabbitProducer  rabbit.Producer
	checkPeriod     time.Duration
	eventTTL        time.Duration
}

func NewSchedulerService(eventRepository event.Repository, log *slog.Logger, rabbitProducer rabbit.Producer, checkPeriod time.Duration, eventTTL time.Duration) *Service {
	return &Service{
		eventRepository: eventRepository,
		log:             log,
		rabbitProducer:  rabbitProducer,
		checkPeriod:     checkPeriod,
		eventTTL:        eventTTL * 60 * 60 * 24,
	}
}
