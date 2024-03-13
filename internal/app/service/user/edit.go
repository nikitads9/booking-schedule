package user

import (
	"booking-schedule/internal/app/model"
	"context"
)

func (s *Service) EditUser(ctx context.Context, user *model.UpdateUserInfo) error {
	return s.userRepository.EditUser(ctx, user)
}
