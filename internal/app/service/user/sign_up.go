package user

import (
	"context"
	"event-schedule/internal/app/model"
	"event-schedule/internal/app/service/user/security"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

func (s *Service) SignUp(ctx context.Context, user *model.User) (string, error) {
	const op = "user.service.SignUp"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	hashedPassword, err := security.HashPassword(*user.Password)
	if err != nil {
		log.Error("failed to hash password", err)
		return "", ErrHashFailed
	}
	user.Password = &hashedPassword

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		log.Error("failed to create user", err)
		return "", err
	}

	jwtToken, err := s.jwtService.GenerateToken(ctx, id)
	if err != nil {
		log.Error("failed to generate token", err)
		return "", err
	}

	return jwtToken, nil
}
