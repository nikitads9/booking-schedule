package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type RabbitConsumer struct {
	DSN       string `yaml:"dsn" env:"AMQP_DSN" env-default:"amqp://guest:guest@queue:5672/bookings"`
	QueueName string `yaml:"queue_name" env:"AMQP_QUEUE" env-default:"bookings"`
}

type SenderConfig struct {
	Env            string         `yaml:"env" env:"env" env-default:"dev"`
	RabbitConsumer RabbitConsumer `yaml:"rabbit_consumer"`
}

func ReadSenderConfigFile(path string) (*SenderConfig, error) {
	config := &SenderConfig{}

	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadSenderConfigEnv() (*SenderConfig, error) {
	config := &SenderConfig{}

	err := cleanenv.ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// GetRabbitConsumerConfig ...
func (s *SenderConfig) GetRabbitConsumerConfig() *RabbitConsumer {
	return &s.RabbitConsumer
}

// GetEnv ...
func (s *SenderConfig) GetEnv() string {
	return s.Env
}
