package config

import (
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Scheduler struct {
	CheckPeriodSec int64 `yaml:"check_period_sec" env:"SCHEDULER_PERIOD" env-default:"60"`
	BookingTTL     int64 `yaml:"booking_ttl_days" env:"BOOKING_TTL" env-default:"365"`
}

type RabbitProducer struct {
	DSN       string `yaml:"dsn" env:"AMQP_DSN" env-default:"amqp://guest:guest@queue:5672/bookings"`
	QueueName string `yaml:"queue_name" env:"AMQP_QUEUE" env-default:"bookings"`
}

type SchedulerConfig struct {
	Env            string         `yaml:"env" env:"env" env-default:"dev"`
	Scheduler      Scheduler      `yaml:"scheduler"`
	Database       Database       `yaml:"database"`
	RabbitProducer RabbitProducer `yaml:"rabbit_producer"`
	Tracer         Tracer         `yaml:"tracer"`
}

func ReadSchedulerConfigFile(path string) (*SchedulerConfig, error) {
	config := &SchedulerConfig{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadSchedulerConfigEnv() (*SchedulerConfig, error) {
	config := &SchedulerConfig{}

	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetSchedulerConfig ...
func (s *SchedulerConfig) GetSchedulerConfig() *Scheduler {
	return &s.Scheduler
}

// GetRabbitProducerConfig ...
func (s *SchedulerConfig) GetRabbitProducerConfig() *RabbitProducer {
	return &s.RabbitProducer
}

// GetTracerConfig
func (e *SchedulerConfig) GetTracerConfig() *Tracer {
	return &e.Tracer
}

// GetEnv ...
func (s *SchedulerConfig) GetEnv() string {
	return s.Env
}

func (s *SchedulerConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", s.Database.User, s.Database.Name, s.Database.Password, s.Database.Host, s.Database.Port, s.Database.Ssl)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = s.Database.MaxOpenedConnections

	return poolConfig, nil
}
