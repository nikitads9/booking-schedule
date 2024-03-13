package auth

import (
	"booking-schedule/internal/app/service/user"

	"go.opentelemetry.io/otel/trace"
)

type Implementation struct {
	User   *user.Service
	Tracer trace.Tracer
}

func NewImplementation(user *user.Service, tracer trace.Tracer) *Implementation {
	return &Implementation{
		User:   user,
		Tracer: tracer,
	}
}
