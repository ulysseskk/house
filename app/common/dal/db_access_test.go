package dal

import (
	"context"
	"testing"

	"github.com/abyss414/house/app/common/model/statistic"

	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/connector"
	"github.com/abyss414/house/app/common/model/db"
)

func TestMigrate(t *testing.T) {
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
	connector.InitMysql()
	err := connector.GetMysqlConnector(context.Background()).AutoMigrate(&db.AmapPoiCode{},
		&db.ErShouFangDetail{},
		&db.ErShouFangPriceHistory{},
		&db.Position{},
		&db.PositionCount{},
		&db.PositionDistance{},
		&db.PositionMetaData{},
		&db.Region{},
		&db.RentDetail{},
		&db.RentPriceHistory{},
		&db.ScanHistory{},
		&db.SystemCache{},
		&db.Tag{},
		&db.Xiaoqu{},
		&db.XiaoquPosition{},
		&statistic.ScrapeStat{})
	if err != nil {
		panic(err)
	}
}
