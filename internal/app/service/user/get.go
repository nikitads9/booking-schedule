package user

import (
	"context"
	"event-schedule/internal/app/model"
)

func (s *Service) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	return s.userRepository.GetUser(ctx, userID)
}
