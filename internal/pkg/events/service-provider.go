package events

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"

	"event-schedule/internal/app/api/handlers"
	eventRepository "event-schedule/internal/app/repository/event"
	userRepository "event-schedule/internal/app/repository/user"
	eventService "event-schedule/internal/app/service/event"
	"event-schedule/internal/app/service/jwt"
	userService "event-schedule/internal/app/service/user"
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

	userRepository userRepository.Repository
	userService    *userService.Service
	jwtService     jwt.Service

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
			log.Fatalf("could not get events-api config: %s", err)
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

func (s *serviceProvider) GetUserRepository(ctx context.Context) userRepository.Repository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewUserRepository(s.GetDB(ctx), s.GetLogger())
		return s.userRepository
	}

	return s.userRepository
}

func (s *serviceProvider) GetEventService(ctx context.Context) *eventService.Service {
	if s.eventService == nil {
		eventRepository := s.GetEventRepository(ctx)
		s.eventService = eventService.NewEventService(eventRepository, s.GetJWTService(), s.GetLogger(), s.TxManager(ctx))
	}

	return s.eventService
}

func (s *serviceProvider) GetUserService(ctx context.Context) *userService.Service {
	if s.userService == nil {
		userRepository := s.GetUserRepository(ctx)
		s.userService = userService.NewUserService(userRepository, s.GetJWTService(), s.GetLogger())
	}

	return s.userService
}

func (s *serviceProvider) GetJWTService() jwt.Service {
	if s.jwtService == nil {
		s.jwtService = jwt.NewJWTService(s.GetConfig().GetJWTConfig().Secret, s.GetConfig().GetJWTConfig().Expiration, s.GetLogger())
	}

	return s.jwtService
}

func (s *serviceProvider) GetEventImpl(ctx context.Context) *handlers.Implementation {
	if s.eventImpl == nil {
		s.eventImpl = handlers.NewImplementation(s.GetEventService(ctx), s.GetUserService(ctx))
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
			TLSConfig: &tls.Config{
				MinVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
			},
		}
	}

	return s.server
}

func (s *serviceProvider) GetLogger() *slog.Logger {
	if s.log == nil {
		env := s.GetConfig().GetEnv()
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
