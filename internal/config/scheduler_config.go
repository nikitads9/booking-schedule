package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type RabbitProducer struct {
	DSN       string `yaml:"dsn"`
	QueueName string `yaml:"queue_name"`
}

type Scheduler struct {
	CheckPeriodSec int64 `yaml:"check_period_sec"`
}

type SchedulerConfig struct {
	Scheduler      *Scheduler      `yaml:"scheduler"`
	Database       *Database       `yaml:"database"`
	RabbitProducer *RabbitProducer `yaml:"rabbit_producer"`
	Logger         *Logger         `yaml:"logger"`
}

func ReadSchedulerConfig(path string) (*SchedulerConfig, error) {
	config := &SchedulerConfig{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
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

// GetLoggerConfig ...
func (s *SchedulerConfig) GetLoggerConfig() *Logger {
	return s.Logger
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
