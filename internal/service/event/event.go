package event

import (
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/repository/event"
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
	ErrNoModel = errors.New("received no model")
	//ErrAvailibility  = errors.New("could not scan for availibility")
	ErrNotAvailible = errors.New("this period is not availible for booking")
	//ErrAddition      = errors.New("failed to add event")
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
