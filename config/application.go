package config

type Configuration struct {
	AppName    string   `yaml:"appName" mapstructure:"appName"`
	Listen     string   `yaml:"listen" mapstructure:"listen"`
	Port       int      `yaml:"port" mapstructure:"port"`
	Log        Log      `yaml:"log" mapstructure:"log"`
	Env        string   `yaml:"env" mapstructure:"env"`
	Config     string   `yaml:"config" mapstructure:"config"`
	Database   Database `yaml:"database" mapstructure:"database"`
	Jwt        Jwt      `yaml:"jwt" mapstructure:"jwt"`
	KubeConfig string   `yaml:"kubeConfig" mapstructure:"kubeConfig"`
	InCluster  bool     `yaml:"inCluster" mapstructure:"inCluster"`
}
