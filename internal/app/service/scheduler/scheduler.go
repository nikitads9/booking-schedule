package scheduler

import (
	"booking-schedule/internal/app/repository/booking"
	"booking-schedule/internal/pkg/rabbit"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	bookingRepository booking.Repository
	log               *slog.Logger
	tracer            trace.Tracer
	rabbitProducer    rabbit.Producer
	checkPeriod       time.Duration
	bookingTTL        time.Duration
}

func NewSchedulerService(bookingRepository booking.Repository, log *slog.Logger, tracer trace.Tracer, rabbitProducer rabbit.Producer, checkPeriod time.Duration, bookingTTL time.Duration) *Service {
	return &Service{
		bookingRepository: bookingRepository,
		log:               log,
		tracer:            tracer,
		rabbitProducer:    rabbitProducer,
		checkPeriod:       checkPeriod,
		bookingTTL:        bookingTTL,
	}
}
