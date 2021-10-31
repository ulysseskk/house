package unit_price

import (
	"context"
	"testing"

	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/connector"
)

func TestScanForFillUnitPrice(t *testing.T) {
	config.SetGlobalConfig(&config.Config{
		OCR: &config.OCRConfig{
			Host: "http://39.100.142.29:19706",
		},
		Redis: &config.RedisConfig{
			Host: "192.168.50.106",
			Port: 6379,
		},
		Mysql: &config.MysqlConfig{
			Host:     "192.168.50.106",
			Port:     3306,
			DbName:   "er_shou_fang",
			User:     "root",
			Password: "Khs19940718!",
		}})
	connector.InitOcrConnector()
	connector.InitGlobalRedisClient()
	connector.InitMysql()
	ScanForFillUnitPrice(context.Background())
}
