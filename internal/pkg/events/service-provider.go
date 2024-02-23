package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"event-schedule/internal/app/api/handlers"
	eventRepository "event-schedule/internal/app/repository/event"
	eventService "event-schedule/internal/app/service/event"
	"event-schedule/internal/config"
	"event-schedule/internal/pkg/db"
	"event-schedule/internal/pkg/db/transaction"
)

type serviceProvider struct {
	db         db.Client
	txManager  db.TxManager
	configPath string
	config     *config.Config

	server *http.Server
	log    *slog.Logger

	eventRepository eventRepository.Repository
	eventService    *eventService.Service

	eventImpl *handlers.Implementation
}

func newServiceProvider(configPath string) *serviceProvider {
	return &serviceProvider{
		configPath: configPath,
	}
}

func (s *serviceProvider) GetDB(ctx context.Context) db.Client {
	if s.db == nil {
		cfg, err := s.GetConfig().GetDBConfig()
		if err != nil {
			s.log.Error("could not get config err: %s", err)
		}
		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			s.log.Error("coud not connect to db err: %s", err)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.Config {
	if s.config == nil {
		cfg, err := config.Read(s.configPath)
		if err != nil {
			log.Fatalf("could not get config err: %s", err.Error())
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetEventRepository(ctx context.Context, log *slog.Logger) eventRepository.Repository {
	if s.eventRepository == nil {
		s.eventRepository = eventRepository.NewEventRepository(s.GetDB(ctx), log)
		return s.eventRepository
	}

	return s.eventRepository
}

func (s *serviceProvider) GetEventService(ctx context.Context, log *slog.Logger) *eventService.Service {
	if s.eventService == nil {
		eventRepository := s.GetEventRepository(ctx, log)
		s.eventService = eventService.NewEventService(eventRepository, log, s.TxManager(ctx))
	}

	return s.eventService
}

func (s *serviceProvider) GetEventImpl(ctx context.Context) *handlers.Implementation {
	if s.eventImpl == nil {
		s.eventImpl = handlers.NewImplementation(s.GetEventService(ctx, s.setupLogger()))
	}

	return s.eventImpl
}

func (s *serviceProvider) getServer(router http.Handler) *http.Server {
	if s.server == nil {
		address, err := s.GetConfig().GetAddress()
		if err != nil {
			s.log.Error("could not get server address err: %s", err)
		}
		s.server = &http.Server{
			Addr:         address,
			Handler:      router,
			ReadTimeout:  s.GetConfig().Server.Timeout,
			WriteTimeout: s.GetConfig().Server.Timeout,
			IdleTimeout:  s.GetConfig().Server.IdleTimeout,
		}
	}

	return s.server
}

func (s *serviceProvider) setupLogger() *slog.Logger {
	env := s.GetConfig().Server.Env
	switch env {
	case envLocal:
		s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	s.log.With(slog.String("env", s.GetConfig().Server.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	return s.log
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDB(ctx).DB())
	}

	return s.txManager
}
