package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

// TODO: move to env
const (
	dbPassEscSeq = "{password}"
	password     = "events_pass"
)

type Server struct {
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Host                 string `yaml:"host"`
	Port                 string `yaml:"port"`
	User                 string `yaml:"user"`
	Name                 string `yaml:"database"`
	Ssl                  string `yaml:"ssl"`
	MaxOpenedConnections int32  `yaml:"max_opened_connections"`
}

type Logger struct {
	Env string `yaml:"env"`
}

type EventConfig struct {
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"database"`
	Logger   *Logger   `yaml:"logger"`
}

func ReadEventConfig(path string) (*EventConfig, error) {
	config := &EventConfig{}

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

// GetServerConfig ...
func (e *EventConfig) GetServerConfig() *Server {
	return e.Server
}

// GetLoggerConfig ...
func (e *EventConfig) GetLoggerConfig() *Logger {
	return e.Logger
}

func (e *EventConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password={password} host=%s port=%s sslmode=%s", e.Database.User, e.Database.Name, e.Database.Host, e.Database.Port, e.Database.Ssl)
	dbDsn = strings.ReplaceAll(dbDsn, dbPassEscSeq, password)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = e.Database.MaxOpenedConnections

	return poolConfig, nil
}

func (c *EventConfig) GetAddress() (string, error) {
	address := c.GetServerConfig().Host + c.GetServerConfig().Port
	//TODO: regex check
	return address, nil
}
