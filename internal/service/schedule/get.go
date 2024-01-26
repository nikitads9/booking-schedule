package schedule

import "context"

func (s *Service) GetDayEvents(ctx context.Context, userID string) (string, error) {
	return s.scheduleRepository.GetEvents(ctx, userID)
}
