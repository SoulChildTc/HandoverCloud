package config

type Log struct {
	Path         string `yaml:"path" mapstructure:"path"`
	Level        string `yaml:"level" mapstructure:"level"`
	Console      bool   `yaml:"console" mapstructure:"console"`
	CloseFileLog bool   `yaml:"closeFileLog" mapstructure:"closeFileLog"`
	Rotate       Rotate `yaml:"rotate" mapstructure:"rotate"`
}

type Rotate struct {
	Enable     bool `yaml:"enable" mapstructure:"enable"`
	MaxSize    int  `yaml:"maxSize" mapstructure:"maxsize"`
	MaxBackups int  `yaml:"maxBackups" mapstructure:"maxbackups"`
	MaxAge     int  `yaml:"maxAge" mapstructure:"maxage"`
	Compress   bool `yaml:"compress" mapstructure:"compress"`
	Localtime  bool `yaml:"localtime" mapstructure:"localtime"`
}
