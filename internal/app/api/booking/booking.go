package booking

import (
	"booking-schedule/internal/app/service/booking"

	"go.opentelemetry.io/otel/trace"
)

type Implementation struct {
	booking *booking.Service
	tracer  trace.Tracer
}

func NewImplementation(booking *booking.Service, tracer trace.Tracer) *Implementation {
	return &Implementation{
		booking: booking,
		tracer:  tracer,
	}
}
