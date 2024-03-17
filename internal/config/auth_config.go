package config

import (
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthServer struct {
	Host        string        `yaml:"host" env:"AUTH_HOST" env-default:"0.0.0.0"`
	Port        string        `yaml:"port" env:"AUTH_PORT" env-default:"5000"`
	Timeout     time.Duration `yaml:"timeout" env:"AUTH_TIMEOUT" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"AUTH_IDLE_TIMEOUT" env-default:"30s"`
}

type AuthConfig struct {
	Env      string     `yaml:"env" env:"env" env-default:"dev"`
	Server   AuthServer `yaml:"server"`
	Database Database   `yaml:"database"`
	Jwt      JWT        `yaml:"jwt"`
	Tracer   Tracer     `yaml:"tracer"`
}

func ReadAuthConfigFile(path string) (*AuthConfig, error) {
	config := &AuthConfig{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadAuthConfigEnv() (*AuthConfig, error) {
	config := &AuthConfig{}

	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetServerConfig ...
func (a *AuthConfig) GetServerConfig() *AuthServer {
	return &a.Server
}

// GetJWTConfig
func (a *AuthConfig) GetJWTConfig() *JWT {
	return &a.Jwt
}

// GetTracerConfig
func (a *AuthConfig) GetTracerConfig() *Tracer {
	return &a.Tracer
}

// GetEnv ...
func (a *AuthConfig) GetEnv() string {
	return a.Env
}

func (a *AuthConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", a.Database.User, a.Database.Name, a.Database.Password, a.Database.Host, a.Database.Port, a.Database.Ssl)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = a.Database.MaxOpenedConnections

	return poolConfig, nil
}

func (c *AuthConfig) GetAddress() (string, error) {
	address := c.GetServerConfig().Host + ":" + c.GetServerConfig().Port
	//TODO: regex check
	return address, nil
}
