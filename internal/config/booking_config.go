package config

import (
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingServer struct {
	Host        string        `yaml:"host" env:"BOOKINGS_HOST" env-default:"0.0.0.0"`
	Port        string        `yaml:"port" env:"BOOKINGS_PORT" env-default:"3000"`
	Timeout     time.Duration `yaml:"timeout" env:"BOOKINGS_TIMEOUT" env-default:"6s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"BOOKINGS_IDLE_TIMEOUT" env-default:"30s"`
}

type Database struct {
	Host                 string `yaml:"host" env:"DB_HOST" env-default:"db"`
	Port                 string `yaml:"port" env:"DB_PORT" env-default:"5433"`
	Name                 string `yaml:"database" env:"DB_NAME" env-default:"bookings_db"`
	User                 string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password             string `yaml:"password" env:"DB_PASSWORD" env-default:"bookings_pass"`
	Ssl                  string `yaml:"ssl" env:"DB_SSL" env-default:"disable"`
	MaxOpenedConnections int32  `yaml:"max_opened_connections" env:"DB_MAX_CONN" env-default:"10"`
}

type JWT struct {
	Secret     string        `yaml:"secret" env:"JWT_SIGNING_KEY" env-default:"verysecretivejwt"`
	Expiration time.Duration `yaml:"expiration" env:"JWT_EXPIRATION" env-default:"2160h"`
}

type Tracer struct {
	EndpointURL  string  `yaml:"endpoint_url" env:"TRACER_URL" env-default:"http://otelcol:4318"`
	SamplingRate float64 `yaml:"sampling_rate" env:"TRACER_SAMPLING_RATE" env-default:"1.0"`
}

type BookingConfig struct {
	Env      string        `yaml:"env" env:"env" env-default:"dev"`
	Server   BookingServer `yaml:"server"`
	Database Database      `yaml:"database"`
	Jwt      JWT           `yaml:"jwt"`
	Tracer   Tracer        `yaml:"tracer"`
}

func ReadBookingConfigFile(path string) (*BookingConfig, error) {
	config := &BookingConfig{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadBookingConfigEnv() (*BookingConfig, error) {
	config := &BookingConfig{}

	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetServerConfig ...
func (b *BookingConfig) GetServerConfig() *BookingServer {
	return &b.Server
}

// GetJWTConfig
func (b *BookingConfig) GetJWTConfig() *JWT {
	return &b.Jwt
}

// GetTracerConfig
func (b *BookingConfig) GetTracerConfig() *Tracer {
	return &b.Tracer
}

// GetEnv ...
func (b *BookingConfig) GetEnv() string {
	return b.Env
}

func (b *BookingConfig) GetDBConfig() (*pgxpool.Config, error) {
	dbDsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s", b.Database.User, b.Database.Name, b.Database.Password, b.Database.Host, b.Database.Port, b.Database.Ssl)

	poolConfig, err := pgxpool.ParseConfig(dbDsn)
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	poolConfig.MaxConns = b.Database.MaxOpenedConnections

	return poolConfig, nil
}

func (c *BookingConfig) GetAddress() (string, error) {
	address := c.GetServerConfig().Host + ":" + c.GetServerConfig().Port
	//TODO: regex check
	return address, nil
}
