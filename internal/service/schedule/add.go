package schedule

import "context"

func (s *Service) AddEvent(ctx context.Context) (string, error) {
	return s.scheduleRepository.AddEvent(ctx)
}
