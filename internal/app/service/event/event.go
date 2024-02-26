package event

import (
	"errors"
	"event-schedule/internal/app/repository/event"
	"event-schedule/internal/pkg/db"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
	eventRepository event.Repository
	log             *slog.Logger
	txManager       db.TxManager
}

// TODO: clean up spare errors
var (
	ErrNoModel       = errors.New("received no model")
	ErrNotAvailible  = errors.New("this period is not availible for booking")
	ErrEmptyUpdate   = errors.New("no parameters for update received")
	ErrNoTransaction = errors.New("can't begin transaction, no connection to database")
	pgNoConnection   = new(*pgconn.ConnectError)
)

func NewEventService(eventRepository event.Repository, log *slog.Logger, txManager db.TxManager) *Service {
	return &Service{
		eventRepository: eventRepository,
		log:             log,
		txManager:       txManager,
	}
}
