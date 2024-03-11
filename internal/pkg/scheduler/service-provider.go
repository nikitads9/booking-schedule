package scheduler

import (
	bookingRepository "booking-schedule/internal/app/repository/booking"
	schedulerService "booking-schedule/internal/app/service/scheduler"
	"booking-schedule/internal/config"
	"booking-schedule/internal/pkg/db"
	"booking-schedule/internal/pkg/rabbit"
	"context"
	"log"
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

	log               *slog.Logger
	rabbitProducer    rabbit.Producer
	bookingRepository bookingRepository.Repository
	schedulerService  *schedulerService.Service
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
			s.GetLogger().Error("could not get db config: %s", err)
			os.Exit(1)
		}
		dbc, err := db.NewClient(ctx, cfg)
		if err != nil {
			s.GetLogger().Error("could not connect to db: %s", err)
			os.Exit(1)
		}
		s.db = dbc
	}

	return s.db
}

func (s *serviceProvider) GetConfig() *config.SchedulerConfig {
	if s.config == nil {
		cfg, err := config.ReadSchedulerConfig(s.configPath)
		if err != nil {
			log.Fatalf("could not get scheduler config: %s", err)
		}

		s.config = cfg
	}

	return s.config
}

func (s *serviceProvider) GetBookingRepository(ctx context.Context) bookingRepository.Repository {
	if s.bookingRepository == nil {
		s.bookingRepository = bookingRepository.NewBookingRepository(s.GetDB(ctx), s.GetLogger())
		return s.bookingRepository
	}

	return s.bookingRepository
}

func (s *serviceProvider) GetSchedulerService(ctx context.Context) *schedulerService.Service {
	if s.schedulerService == nil {
		s.schedulerService = schedulerService.NewSchedulerService(
			s.GetBookingRepository(ctx),
			s.GetLogger(),
			s.GetRabbitProducer(),
			time.Duration(s.GetConfig().GetSchedulerConfig().CheckPeriodSec)*time.Second,
			time.Duration(s.GetConfig().GetSchedulerConfig().BookingTTL)*time.Hour*24)
	}

	return s.schedulerService
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

		s.log.With(slog.String("env", env))
	}

	return s.log
}

// GetRabbitProducer ...
func (s *serviceProvider) GetRabbitProducer() rabbit.Producer {
	if s.rabbitProducer == nil {
		rp, err := rabbit.NewProducer(s.GetConfig().GetRabbitProducerConfig())
		if err != nil {
			s.log.Error("could not connect to rabbit producer: %s", err)
			os.Exit(1)
		}
		s.rabbitProducer = rp
	}

	return s.rabbitProducer
}
