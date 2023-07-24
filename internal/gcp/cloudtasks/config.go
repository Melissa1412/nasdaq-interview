package cloudtasks

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	TaskQueues map[string]TaskQueue `mapstructure:"cloudTasks" validate:"dive"`
}

type TaskQueue struct {
	Name                    string  `json:"name" validate:"required"`
	Region                  string  `json:"region" validate:"required"`
	ProjectId               string  `json:"projectId" validate:"required"`
	MinBackoff              string  `json:"minBackoff"`
	MaxBackoff              string  `json:"maxBackoff"`
	MaxConcurrentDispatches int  	`json:"maxConcurrentDispatches"`
	MaxDispatchesPerSecond  float64 `json:"maxDispatchesPerSecond"`
	ClientName              string
}

func GetTaskConfig(viperConfig *viper.Viper, clientName string) (*Config, error) {
	if viperConfig == nil {
		return nil, nil
	}

	var taskConfig Config
	err := viperConfig.Unmarshal(&taskConfig)
	if err != nil {
		return nil, err
	}

	for name, task := range taskConfig.TaskQueues {
		task.Name = name
		task.ClientName = clientName

		taskConfig.TaskQueues[name] = task
	}
	return &taskConfig, nil
}

func ValidateConfig(config *Config) error {
	v := validator.New()
	if err := v.Struct(config); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("%s validate failed on the %s rule", err.Namespace(), err.Tag())
		}
	}
	return nil
}
