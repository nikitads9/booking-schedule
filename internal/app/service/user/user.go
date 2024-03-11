package user

import (
	"booking-schedule/internal/app/repository/user"
	"booking-schedule/internal/app/service/jwt"
	"errors"
	"log/slog"
)

type Service struct {
	userRepository user.Repository
	jwtService     jwt.Service
	log            *slog.Logger
}

var (
	ErrBadLogin   = errors.New("incorrect nickname")
	ErrBadPasswd  = errors.New("incorrect password")
	ErrHashFailed = errors.New("failed to hash password")
)

func NewUserService(userRepository user.Repository, jwtService jwt.Service, log *slog.Logger) *Service {
	return &Service{
		userRepository: userRepository,
		jwtService:     jwtService,
		log:            log,
	}
}
