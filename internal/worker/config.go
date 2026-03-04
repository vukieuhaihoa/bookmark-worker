package worker

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	QueueName   string `envconfig:"QUEUE_NAME" default:"bookmark_import_queue"`
	ServiceName string `envconfig:"SERVICE_NAME" default:"bookmark-worker"`
	InstanceID  string `envconfig:"INSTANCE_ID" default:""`
	LOG_LEVEL   string `envconfig:"LOG_LEVEL" default:"debug"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	return cfg, nil
}
