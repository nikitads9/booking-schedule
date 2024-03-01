package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type RabbitConsumer struct {
	DSN       string `yaml:"dsn"`
	QueueName string `yaml:"queue_name"`
}

type SenderConfig struct {
	RabbitConsumer *RabbitConsumer `yaml:"rabbit_consumer"`
	Logger         *Logger         `yaml:"logger"`
}

func ReadSenderConfig(path string) (*SenderConfig, error) {
	config := &SenderConfig{}

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

// GetRabbitConsumerConfig ...
func (s *SenderConfig) GetRabbitConsumerConfig() *RabbitConsumer {
	return s.RabbitConsumer
}

// GetLoggerConfig ...
func (s *SenderConfig) GetLoggerConfig() *Logger {
	return s.Logger
}
