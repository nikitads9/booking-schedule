package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type Scheduler struct {
	CheckPeriodSec int64 `yaml:"check_period_sec"`
	BookingTTL     int64 `yaml:"booking_ttl_days"`
}

type RabbitProducer struct {
	DSN       string `yaml:"dsn"`
	QueueName string `yaml:"queue_name"`
}

type SchedulerConfig struct {
	Scheduler      *Scheduler      `yaml:"scheduler"`
	Database       *Database       `yaml:"database"`
	RabbitProducer *RabbitProducer `yaml:"rabbit_producer"`
	Env            string          `yaml:"env"`
}

func ReadSchedulerConfig(path string) (*SchedulerConfig, error) {
	config := &SchedulerConfig{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetSchedulerConfig ...
func (s *SchedulerConfig) GetSchedulerConfig() *Scheduler {
	return s.Scheduler
}

// GetRabbitProducerConfig ...
func (s *SchedulerConfig) GetRabbitProducerConfig() *RabbitProducer {
	return s.RabbitProducer
}

// GetEnv ...
func (s *SchedulerConfig) GetEnv() string {
	return s.Env
}

func (s *SchedulerConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password={password} host=%s port=%s sslmode=%s", s.Database.User, s.Database.Name, s.Database.Host, s.Database.Port, s.Database.Ssl)
	dbDsn = strings.ReplaceAll(dbDsn, dbPassEscSeq, password)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = s.Database.MaxOpenedConnections

	return poolConfig, nil
}
