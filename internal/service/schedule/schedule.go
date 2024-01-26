package schedule

import "event-schedule/internal/repository/schedule"

type Service struct {
	scheduleRepository schedule.Repository
}

func NewScheduleService(scheduleRepository schedule.Repository) *Service {
	return &Service{
		scheduleRepository: scheduleRepository,
	}
}
