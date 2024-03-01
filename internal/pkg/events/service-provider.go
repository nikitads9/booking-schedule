package events

import (
	"context"
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

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	db         db.Client
	txManager  db.TxManager
	configPath string
	config     *config.EventConfig

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
			s.log.Error("could not get db config: %s", err)
		}
		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			s.log.Error("coud not connect to db: %s", err)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.EventConfig {
	if s.config == nil {
		cfg, err := config.ReadEventConfig(s.configPath)
		if err != nil {
			s.log.Error("coud not get events-api config: %s", err)
			os.Exit(1)
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetEventRepository(ctx context.Context) eventRepository.Repository {
	if s.eventRepository == nil {
		s.eventRepository = eventRepository.NewEventRepository(s.GetDB(ctx), s.GetLogger())
		return s.eventRepository
	}

	return s.eventRepository
}

func (s *serviceProvider) GetEventService(ctx context.Context) *eventService.Service {
	if s.eventService == nil {
		eventRepository := s.GetEventRepository(ctx)
		s.eventService = eventService.NewEventService(eventRepository, s.GetLogger(), s.TxManager(ctx))
	}

	return s.eventService
}

func (s *serviceProvider) GetEventImpl(ctx context.Context) *handlers.Implementation {
	if s.eventImpl == nil {
		s.eventImpl = handlers.NewImplementation(s.GetEventService(ctx))
	}

	return s.eventImpl
}

func (s *serviceProvider) getServer(router http.Handler) *http.Server {
	if s.server == nil {
		address, err := s.GetConfig().GetAddress()
		if err != nil {
			s.log.Error("could not get server address: %s", err)
			return nil
		}
		s.server = &http.Server{
			Addr:         address,
			Handler:      router,
			ReadTimeout:  s.GetConfig().GetServerConfig().Timeout,
			WriteTimeout: s.GetConfig().GetServerConfig().Timeout,
			IdleTimeout:  s.GetConfig().GetServerConfig().IdleTimeout,
		}
	}

	return s.server
}

func (s *serviceProvider) GetLogger() *slog.Logger {
	if s.log == nil {
		//TODO: move env to logger config
		env := s.GetConfig().GetLoggerConfig().Env
		switch env {
		case envLocal:
			s.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envDev:
			s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		case envProd:
			s.log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}

		s.log.With(slog.String("env", env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	}

	return s.log
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDB(ctx).DB())
	}

	return s.txManager
}
