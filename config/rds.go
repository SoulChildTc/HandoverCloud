package config

type Database struct {
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Charset         string `yaml:"charset" mapstructure:"charset"`
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Database        string `yaml:"database" mapstructure:"database"`
	MaxOpenConns    int    `yaml:"maxOpenConns" mapstructure:"maxopenconns"`
	MaxIdleConns    int    `yaml:"maxIdleConns" mapstructure:"maxidleconns"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime" mapstructure:"connmaxidletime"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime" mapstructure:"connmaxlifetime"`
	LogLevel        string `yaml:"logLevel" mapstructure:"loglevel"`
	Driver          string `yaml:"driver" mapstructure:"driver"`
	ReportCaller    bool   `yaml:"reportCaller" mapstructure:"reportcaller"`
	Path            string `yaml:"path" mapstructure:"path"`
}
