package handlers

import (
	booking "booking-schedule/internal/app/service/booking"
	"booking-schedule/internal/app/service/user"
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
