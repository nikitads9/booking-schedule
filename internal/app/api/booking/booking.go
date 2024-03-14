package booking

import (
	booking "booking-schedule/internal/app/service/booking"
	"booking-schedule/internal/app/service/user"

	"go.opentelemetry.io/otel/trace"
)

type Implementation struct {
	booking *booking.Service
	user    *user.Service
	tracer  trace.Tracer
}

func NewImplementation(booking *booking.Service, user *user.Service, tracer trace.Tracer) *Implementation {
	return &Implementation{
		booking: booking,
		user:    user,
		tracer:  tracer,
	}
}
