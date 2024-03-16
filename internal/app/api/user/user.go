package user

import (
	"booking-schedule/internal/app/service/user"

	"go.opentelemetry.io/otel/trace"
)

type Implementation struct {
	user   *user.Service
	tracer trace.Tracer
}

func NewImplementation(user *user.Service, tracer trace.Tracer) *Implementation {
	return &Implementation{
		user:   user,
		tracer: tracer,
	}
}
