package config

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

type Config struct {
	Mysql    *MysqlConfig    `json:"mysql" yaml:"mysql"`
	Log      *LogConfig      `json:"log" yaml:"log"`
	Amap     *AmapConfig     `json:"amap" yaml:"amap"`
	OCR      *OCRConfig      `json:"ocr" yaml:"ocr"`
	Redis    *RedisConfig    `json:"redis" yaml:"redis"`
	Scrapper *ScrapperConfig `json:"scrapper" yaml:"scrapper"`
}

var conf *Config

func InitConfig() error {
	c := &config.AppConfig{
		AppID:          os.Getenv("APOLLO_ENV_ID"),
		Cluster:        os.Getenv("APOLLO_CLUSTER"),
		IP:             os.Getenv("APOLLO_ADDRESS"),
		NamespaceName:  os.Getenv("APOLLO_NAMESPACE"),
		IsBackupConfig: true,
		Secret:         os.Getenv("APOLLO_SECRET"),
	}
	apolloLogger := logrus.New()
	apolloLogger.SetLevel(logrus.DebugLevel)
	apolloLogger.SetOutput(os.Stdout)
	agollo.SetLogger(apolloLogger)
	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		return err
	}
	cache := client.GetConfigCache(c.NamespaceName)
	propertiesFileString := ""
	cache.Range(func(key, value interface{}) bool {
		propertiesFileString = fmt.Sprintf("%s\n%+v=%+v", propertiesFileString, key, value)
		return true
	})
	fmt.Println(fmt.Sprintf("读取到配置%s", propertiesFileString))
	viper.SetConfigType("properties")
	err = viper.ReadConfig(bytes.NewBufferString(propertiesFileString))
	if err != nil {
		return err
	}
	configObj := &Config{}
	err = viper.Unmarshal(configObj)
	if err != nil {
		return err
	}
	conf = configObj
	return nil
}

func GlobalConfig() *Config {
	return conf
}

func SetGlobalConfig(manualConfig *Config) {
	conf = manualConfig
}
