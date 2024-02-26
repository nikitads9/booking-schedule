package scheduler

import (
	"context"
	eventRepository "event-schedule/internal/app/repository/event"
	schedulerService "event-schedule/internal/app/service/scheduler"
	"event-schedule/internal/config"
	"event-schedule/internal/pkg/db"
	"event-schedule/internal/pkg/rabbit"
	"log/slog"
	"os"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	db         db.Client
	configPath string
	config     *config.SchedulerConfig

	log              *slog.Logger
	rabbitProducer   rabbit.Producer
	eventRepository  eventRepository.Repository
	schedulerService *schedulerService.Service
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
			s.log.Error("could not connect to db err: %s", err)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.SchedulerConfig {
	if s.config == nil {
		cfg, err := config.ReadSchedulerConfig(s.configPath)
		if err != nil {
			s.log.Error("could not get config: %s", err)
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

func (s *serviceProvider) GetSchedulerService(ctx context.Context) *schedulerService.Service {
	if s.schedulerService == nil {
		eventRepository := s.GetEventRepository(ctx)
		s.schedulerService = schedulerService.NewSchedulerService(eventRepository, s.GetLogger(), s.GetRabbitProducer(), time.Duration(s.GetConfig().GetSchedulerConfig().CheckPeriodSec)*time.Second)
	}

	return s.schedulerService
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

// GetRabbitProducer ...
func (s *serviceProvider) GetRabbitProducer() rabbit.Producer {
	if s.rabbitProducer == nil {
		rp, err := rabbit.NewProducer(s.GetConfig().GetRabbitProducerConfig())
		if err != nil {
			s.log.Error("could not connect to rabbit producer err: %s", err)
			os.Exit(1)
		}
		s.rabbitProducer = rp
	}

	return s.rabbitProducer
}
