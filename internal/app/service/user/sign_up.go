package user

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/app/service/user/security"
	"context"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/codes"
)

func (s *Service) SignUp(ctx context.Context, user *model.User) (string, error) {
	const op = "user.service.SignUp"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	hashedPassword, err := security.HashPassword(user.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to hash password", err)
		return "", ErrHashFailed
	}

	span.AddEvent("password hash created")
	user.Password = hashedPassword

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to create user", err)
		return "", err
	}

	span.AddEvent("created user")

	jwtToken, err := s.jwtService.GenerateToken(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to generate token", err)
		return "", err
	}

	span.AddEvent("signed token acquired")

	return jwtToken, nil
}
