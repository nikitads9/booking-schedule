package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// TODO: move to env
const (
	dbPassEscSeq = "{password}"
	password     = "events_pass"
)

type Database struct {
	Host                 string `yaml:"host"`
	Port                 string `yaml:"port"`
	User                 string `yaml:"user"`
	Name                 string `yaml:"database"`
	Ssl                  string `yaml:"ssl"`
	MaxOpenedConnections int32  `yaml:"max_opened_connections"`
}

type Server struct {
	Env         string        `yaml:"env"`
	Host        string        `yaml:"host"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
}

func Read(path string) (*Config, error) {
	config := &Config{}

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

func (c *Config) GetDBConfig() (string, error) {
	DbDsn := fmt.Sprintf("user=%s dbname=%s password={password} host=%s port=%s sslmode=%s", c.Database.User, c.Database.Name, c.Database.Host, c.Database.Port, c.Database.Ssl)
	DbDsn = strings.ReplaceAll(DbDsn, dbPassEscSeq, password)

	return DbDsn, nil
}

func (c *Config) GetAddress() (string, error) {
	address := c.Server.Host + c.Server.Port
	//TODO: regex check
	return address, nil
}
