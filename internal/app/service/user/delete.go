package user

import (
	"context"
)

func (s *Service) DeleteUser(ctx context.Context, userID int64) error {
	return s.userRepository.DeleteUser(ctx, userID)
}
