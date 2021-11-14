package config

type ScrapperConfig struct {
	EnableLianjia bool `json:"enable_lianjia" mapstructure:"enable_lianjia"`
	EnableZiroom  bool `json:"enable_ziroom" mapstructure:"enable_ziroom"`
	Enable5I5j    bool `json:"enable_5i5j" mapstructure:"enable_5i5j"`
	EnableXiaoqu  bool `json:"enable_xiaoqu" mapstructure:"enable_xiaoqu"`
}
