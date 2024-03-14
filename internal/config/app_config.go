package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/exaring/otelpgx"
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

type Tracer struct {
	EndpointURL  string  `yaml:"endpoint_url"`
	SamplingRate float64 `yaml:"sampling_rate"`
}

type AppConfig struct {
	Server   *Server   `yaml:"server"`
	Database *Database `yaml:"database"`
	Jwt      *JWT      `yaml:"jwt"`
	Env      string    `yaml:"env"`
	Tracer   *Tracer   `yaml:"tracer"`
}

func ReadAppConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}

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
func (e *AppConfig) GetServerConfig() *Server {
	return e.Server
}

// GetJWTConfig
func (e *AppConfig) GetJWTConfig() *JWT {
	return e.Jwt
}

// GetTracerConfig
func (e *AppConfig) GetTracerConfig() *Tracer {
	return e.Tracer
}

// GetEnv ...
func (e *AppConfig) GetEnv() string {
	return e.Env
}

func (e *AppConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password={password} host=%s port=%s sslmode=%s", e.Database.User, e.Database.Name, e.Database.Host, e.Database.Port, e.Database.Ssl)
	dbDsn = strings.ReplaceAll(dbDsn, dbPassEscSeq, password)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = e.Database.MaxOpenedConnections

	return poolConfig, nil
}

func (c *AppConfig) GetAddress() (string, error) {
	address := c.GetServerConfig().Host + c.GetServerConfig().Port
	//TODO: regex check
	return address, nil
}
