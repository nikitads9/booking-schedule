package handlers

import (
	"event-schedule/internal/service/schedule"
)

type Implementation struct {
	Service *schedule.Service
}

func NewImplementation(service *schedule.Service) *Implementation {
	return &Implementation{
		Service: service,
	}
}
