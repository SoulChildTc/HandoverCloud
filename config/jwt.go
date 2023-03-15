package config

import "time"

type Jwt struct {
	Secret string        `yaml:"secret" mapstructure:"secret"`
	Ttl    time.Duration `yaml:"ttl" mapstructure:"ttl"`
}
