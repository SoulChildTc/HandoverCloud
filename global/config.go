package global

import (
	"github.com/spf13/viper"
	"soul/config"
)

var (
	V      *viper.Viper
	Config config.Configuration
)
