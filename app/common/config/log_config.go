package config

type LogConfig struct {
	Method   string        `json:"method" yaml:"method"`
	FilePath string        `json:"file_path" yaml:"file_path" mapstructure:"file_path"`
	Syslog   *SyslogConfig `json:"syslog"`
}

type SyslogConfig struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
}
