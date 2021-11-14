package config

type MysqlConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`
	DbName   string `json:"db_name" yaml:"db_name" mapstructure:"db_name"`
	User     string `json:"user" yaml:"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
}
