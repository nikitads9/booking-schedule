package user

import (
	"context"
	"errors"
	"event-schedule/internal/app/repository/user"
	"event-schedule/internal/app/service/user/security"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

// Login performs the login process using the provided user credentials.
// It retrieves the user from the user repository and generates a JWT token.
// If successful, it returns the generated token.
// If the user cannot be found, it returns ErrBadLogin.
// If there is any other error, it returns a wrapped error.
func (s *Service) SignIn(ctx context.Context, nickname string, pass string) (token string, err error) {
	const op = "user.service.SignIn"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	retrievedUser, err := s.userRepository.GetUserByNickname(ctx, nickname)
	if err != nil {
		log.Error("failed to get user by nickname: ", err)
		if errors.Is(err, user.ErrNotFound) {
			return "", ErrBadLogin
		}
		return "", err
	}

	if ok := security.CheckPasswordHash(pass, *retrievedUser.Password); !ok {
		log.Error("password check failed: ", err)
		return "", ErrBadPasswd
	}

	return s.jwtService.GenerateToken(ctx, retrievedUser.ID)
}
