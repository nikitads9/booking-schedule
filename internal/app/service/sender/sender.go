package sender

import (
	"booking-schedule/internal/pkg/rabbit"
	"log/slog"
)

type Service struct {
	log            *slog.Logger
	rabbitConsumer rabbit.Consumer
}

func NewSenderService(log *slog.Logger, rabbitConsumer rabbit.Consumer) *Service {
	return &Service{
		log:            log,
		rabbitConsumer: rabbitConsumer,
	}
}
