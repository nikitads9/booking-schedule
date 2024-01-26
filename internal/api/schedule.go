package api

import (
	"event-schedule/internal/service/schedule"
)

type Implementation struct {
	*schedule.Service
}

func NewImplementation(scheduleService *schedule.Service) *Implementation {
	return &Implementation{
		scheduleService,
	}
}
