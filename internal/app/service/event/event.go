package event

import (
	"errors"
	"event-schedule/internal/app/repository/event"
	"event-schedule/internal/app/service/jwt"
	"event-schedule/internal/pkg/db"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	eventRepository event.Repository
	jwtService      jwt.Service
	log             *slog.Logger
	txManager       db.TxManager
}

var (
	ErrNoModel      = errors.New("received no model")
	ErrNotAvailible = errors.New("this period is not availible for booking")
	ErrNoConnection = errors.New("can't begin transaction, no connection to database")
	pgNoConnection  = new(*pgconn.ConnectError)
)

func NewEventService(eventRepository event.Repository, jwtService jwt.Service, log *slog.Logger, txManager db.TxManager) *Service {
	return &Service{
		eventRepository: eventRepository,
		jwtService:      jwtService,
		log:             log,
		txManager:       txManager,
	}
}
