package ziroom

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/abyss414/house/app/common/model/db"

	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/connector"
)

func TestGetAllInnerAreaLink(t *testing.T) {
	scrapper := NewZiroomExecutor()
	scrapper.GetAllAreaLink(context.Background(), "")

}

func TestGetPage(t *testing.T) {
	config.SetGlobalConfig(&config.Config{OCR: &config.OCRConfig{Host: "http://39.100.142.29:19706"}})
	connector.InitOcrConnector()
	scrapper := NewZiroomExecutor()
	nextPage, err := scrapper.executeForPage(context.Background(), "https://www.ziroom.com/z/d23008614-b18335647-p3/", []executeFunc{
		func(ctx context.Context, detail *db.RentDetail) error {
			jsonByte, err := json.Marshal(detail)
			if err != nil {
				return nil
			}
			fmt.Println(string(jsonByte))
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(nextPage)
}

func TestBuildCache(t *testing.T) {
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
	details := []*db.RentDetail{}
	err := connector.GetMysqlConnector(context.Background()).Find(&details).Error
	if err != nil {
		panic(err)
	}
	key := "house_id_detail_id"
	for _, detail := range details {
		client, _ := connector.GetRedisClient()
		defer client.Close()
		exist, err := client.HExists(key, detail.HouseId)
		if err != nil {
			panic(err)
		}
		if !exist {
			_, err := client.HSet(key, detail.HouseId, strconv.FormatUint(detail.Id, 10))
			if err != nil {
				panic(err)
			}
		}
	}
}
