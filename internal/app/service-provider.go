package app

import (
	"context"
	"event-schedule/internal/api"
	"event-schedule/internal/client/db"
	"event-schedule/internal/config"
	scheduleRepository "event-schedule/internal/repository/schedule"
	scheduleService "event-schedule/internal/service/schedule"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type serviceProvider struct {
	db db.Client
	//txManager  db.TxManager
	configPath string
	config     *config.Config

	server *http.Server
	log    *slog.Logger

	/* 	scheduleRepository scheduleRepository.Repository
	   	scheduleService    *scheduleService.Service */

	impl *api.Implementation
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

func (s *serviceProvider) GetScheduleRepository(ctx context.Context) scheduleRepository.Repository {
	if s.impl.Service.scheduleRepository == nil {
		s.impl.Service.scheduleRepository = scheduleRepository.NewScheduleRepository(s.GetDB(ctx))
		return s.impl.Service.scheduleRepository
	}

	return s.impl.Service.scheduleRepository
}

func (s *serviceProvider) GetScheduleService(ctx context.Context) *scheduleService.Service {
	if s.scheduleService == nil {
		scheduleRepository := s.GetScheduleRepository(ctx)
		s.scheduleService = scheduleService.NewScheduleService(scheduleRepository)
	}

	return s.scheduleService
}

func (s *serviceProvider) Getimpl(ctx context.Context) *api.Implementation {
	if s.impl == nil {
		s.impl = api.NewImplementation(s.GetScheduleService(ctx))
	}

	return s.impl
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

/* func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDB(ctx).DB())
	}

	return s.txManager
}
*/
