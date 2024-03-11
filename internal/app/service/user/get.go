package user

import (
	"booking-schedule/internal/app/model"
	"context"
)

func (s *Service) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return s.userRepository.GetUser(ctx, userID)
}
