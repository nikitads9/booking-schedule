package scheduler

import (
	"booking-schedule/internal/app/repository/booking"
	"booking-schedule/internal/pkg/rabbit"
	"log/slog"
	"time"
)

type Service struct {
	bookingRepository booking.Repository
	log               *slog.Logger
	rabbitProducer    rabbit.Producer
	checkPeriod       time.Duration
	bookingTTL        time.Duration
}

func NewSchedulerService(bookingRepository booking.Repository, log *slog.Logger, rabbitProducer rabbit.Producer, checkPeriod time.Duration, bookingTTL time.Duration) *Service {
	return &Service{
		bookingRepository: bookingRepository,
		log:               log,
		rabbitProducer:    rabbitProducer,
		checkPeriod:       checkPeriod,
		bookingTTL:        bookingTTL,
	}
}
