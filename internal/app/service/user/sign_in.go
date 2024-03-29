package user

import (
	"booking-schedule/internal/app/repository/user"
	"booking-schedule/internal/app/service/user/security"
	"booking-schedule/internal/logger/sl"
	"context"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Login performs the login process using the provided user credentials.
// It retrieves the user from the user repository and generates a JWT token.
// If successful, it returns the generated token.
// If the user cannot be found, it returns ErrBadLogin.
// If there is any other error, it returns a wrapped error.
func (s *Service) SignIn(ctx context.Context, nickname string, pass string) (token string, err error) {
	const op = "user.service.SignIn"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	ctx, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	retrievedUser, err := s.userRepository.GetUserByNickname(ctx, nickname)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("failed to get user by nickname", sl.Err(err))
		if errors.Is(err, user.ErrNotFound) {
			span.SetStatus(codes.Error, user.ErrNotFound.Error())
			return "", ErrBadLogin
		}
		return "", err
	}

	span.AddEvent("user retrieved", trace.WithAttributes(attribute.Int64("id", retrievedUser.ID)))

	if ok := security.CheckPasswordHash(pass, retrievedUser.Password); !ok {
		span.RecordError(ErrBadPasswd)
		span.SetStatus(codes.Error, ErrBadPasswd.Error())
		log.Error("password check failed", sl.Err(ErrBadPasswd))
		return "", ErrBadPasswd
	}

	span.AddEvent("password checked")

	return s.jwtService.GenerateToken(ctx, retrievedUser.ID)
}
