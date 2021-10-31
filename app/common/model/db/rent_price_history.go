package db

import (
	"context"
	"time"

	"github.com/abyss414/house/app/common/log"

	"github.com/abyss414/house/app/common/connector"
)

type RentPriceHistory struct {
	Id        uint64    `json:"id"`
	Price     float64   `json:"price"`
	HouseId   string    `json:"house_id"`
	PriceStr  string    `json:"price_str"`
	Unit      string    `json:"unit"`
	ScrapedAt time.Time `json:"scraped_at"`
}

func (RentPriceHistory) TableName() string {
	return "rent_price_history"
}

func (history *RentPriceHistory) Create(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Model(&RentPriceHistory{}).Create(history).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("插入RentPriceHistory失败")
		return err
	}
	return nil
}
