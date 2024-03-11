package handlers

import (
	booking "event-schedule/internal/app/service/event"
	"event-schedule/internal/app/service/user"
)

type Implementation struct {
	Booking *booking.Service
	User    *user.Service
}

func NewImplementation(booking *booking.Service, user *user.Service) *Implementation {
	return &Implementation{
		Booking: booking,
		User:    user,
	}
}
