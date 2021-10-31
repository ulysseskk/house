package db

import (
	"context"
	"strconv"

	"github.com/abyss414/house/app/common/log"

	"github.com/abyss414/house/app/common/connector"
)

var (
	ZiroomHouseIdRentDetailIDRedisKey = "house_id_detail_id"
)

type RentDetail struct {
	Id               uint64            `json:"id" gorm:"column:id" gorm:"primary_key"`
	Name             string            `json:"name"`
	HouseId          string            `json:"house_id"`
	Url              string            `json:"url"`
	Platform         string            `json:"platform"`
	XiaoquId         uint64            `json:"xiaoqu_id"`
	XiaoquName       string            `json:"xiaoqu_name"`
	Type             string            `json:"type"` // 整租、合租
	IndependentBath  bool              `json:"independent_bath"`
	Direction        string            `json:"direction"`
	Area             float64           `json:"area"`
	Floor            int               `json:"floor"`
	TotalFloor       int               `json:"total_floor"`
	RoomType         string            `json:"room_type"` //n居室
	RentPriceHistory *RentPriceHistory `json:"rent_price_history" gorm:"-"`
}

func (detail *RentDetail) UpsertRendDetail(ctx context.Context) error {
	// 查看redis
	log.WithContext(ctx).Errorf("开始获取redis client from pool")
	redisClient, err := connector.GetRedisClient()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("获取redis client from pool错误")
		return err
	}
	defer redisClient.Close()
	log.WithContext(ctx).Errorf("获取redis client from pool成功")
	exist, err := redisClient.HExists(ZiroomHouseIdRentDetailIDRedisKey, detail.HouseId)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("HExists错误, out key %s.Inner key %s", ZiroomHouseIdRentDetailIDRedisKey, detail.HouseId)
		return err
	}
	if exist {
		defer func() {
			if r := recover(); r != nil {
				log.GlobalLogger().WithField("recover", r).Errorf("HGET失败, Key %s, innerKey %s", ZiroomHouseIdRentDetailIDRedisKey, detail.HouseId)
			}
		}()
		id, err := redisClient.HGet(ZiroomHouseIdRentDetailIDRedisKey, detail.HouseId)
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("HGET失败")
			return err
		}
		detail.Id, err = strconv.ParseUint(id, 10, 64)
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("Parse id 失败。从redis中获取到的id为%s", id)
			return err
		}
		err = detail.UpdateRentDetail(ctx)
		if err != nil {
			return err
		}
	} else {
		err := detail.InsertRentDetail(ctx)
		if err != nil {
			return err
		}
		_, err = redisClient.HSet(ZiroomHouseIdRentDetailIDRedisKey, detail.HouseId, strconv.FormatUint(detail.Id, 10))
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("HSET失败")
			return err
		}
	}
	if detail.RentPriceHistory != nil && detail.RentPriceHistory.Price != 0 {
		err := detail.RentPriceHistory.Create(ctx)
		if err != nil {
			return err
		}
	}
	if detail.RentPriceHistory != nil && detail.RentPriceHistory.Price == 0 {
		log.WithContext(ctx).WithError(err).WithField("info", detail).Errorf("房源无价格信息")
	}
	return nil
}

func (detail *RentDetail) InsertRentDetail(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Create(detail).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("插入RentDetail失败")
		return err
	}
	return nil
}

func (detail *RentDetail) UpdateRentDetail(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Where("house_id = ?", detail.HouseId).Save(detail).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("更新RentDetail失败")
		return err
	}
	return nil
}
