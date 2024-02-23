package handlers

import (
	schedule "event-schedule/internal/app/service/event"
)

type Implementation struct {
	Service *schedule.Service
}

func NewImplementation(service *schedule.Service) *Implementation {
	return &Implementation{
		Service: service,
	}
}
