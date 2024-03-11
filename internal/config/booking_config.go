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
	password     = "bookings_pass"
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

type JWT struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

type BookingConfig struct {
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"database"`
	Jwt      *JWT      `yaml:"jwt"`
	Env      string    `yaml:"env"`
}

func ReadBookingConfig(path string) (*BookingConfig, error) {
	config := &BookingConfig{}

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

// GetServerConfig ...
func (e *BookingConfig) GetServerConfig() *Server {
	return e.Server
}

// GetJWTConfig
func (e *BookingConfig) GetJWTConfig() *JWT {
	return e.Jwt
}

// GetEnv ...
func (e *BookingConfig) GetEnv() string {
	return e.Env
}

func (e *BookingConfig) GetDBConfig() (*pgxpool.Config, error) {
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

func (c *BookingConfig) GetAddress() (string, error) {
	address := c.GetServerConfig().Host + c.GetServerConfig().Port
	//TODO: regex check
	return address, nil
}
