package sender

import (
	sender "booking-schedule/internal/app/service/sender"
	"booking-schedule/internal/config"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/pkg/rabbit"
	"context"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type serviceProvider struct {
	configType string
	configPath string
	config     *config.SenderConfig

	log            *slog.Logger
	rabbitConsumer rabbit.Consumer
	senderService  *sender.Service
}

func newServiceProvider(configType string, configPath string) *serviceProvider {
	return &serviceProvider{
		configType: configType,
		configPath: configPath,
	}
}

func (s *serviceProvider) GetConfig() *config.SenderConfig {
	if s.config == nil {
		if s.configType == "env" {
			cfg, err := config.ReadSenderConfigEnv()
			if err != nil {
				log.Fatalf("could not get sender config from env: %s", err)
			}

			s.config = cfg
		} else {
			cfg, err := config.ReadSenderConfigFile(s.configPath)
			if err != nil {
				log.Fatalf("could not get sender config from file: %s", err)
			}

			s.config = cfg
		}

	}

	return s.config
}

func (s *serviceProvider) GetSenderService(ctx context.Context) *sender.Service {
	if s.senderService == nil {
		s.senderService = sender.NewSenderService(
			s.GetLogger(),
			s.GetRabbitConsumer())
	}

	return s.senderService
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

// GetRabbitProducer ...
func (s *serviceProvider) GetRabbitConsumer() rabbit.Consumer {
	if s.rabbitConsumer == nil {
		rc, err := rabbit.NewConsumer(s.GetConfig().GetRabbitConsumerConfig())
		if err != nil {
			s.log.Error("could not connect to rabbit consumer", sl.Err(err))
			os.Exit(1)
		}
		s.rabbitConsumer = rc
	}

	return s.rabbitConsumer
}
