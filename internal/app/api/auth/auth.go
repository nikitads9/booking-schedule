package auth

import (
	userRepo "booking-schedule/internal/app/repository/user"
	"booking-schedule/internal/app/service/user"
	"net/http"

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

func GetErrorCode(err error) int {
	switch err {
	case user.ErrBadLogin:
		return http.StatusUnauthorized
	case user.ErrBadPasswd:
		return http.StatusUnauthorized
	case userRepo.ErrNotFound:
		return http.StatusNotFound
	case userRepo.ErrAlreadyExists:
		return http.StatusBadRequest
	case userRepo.ErrDuplicate:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
