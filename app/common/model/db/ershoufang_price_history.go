package db

import (
	"context"

	"github.com/abyss414/house/app/common/log"

	"github.com/abyss414/house/app/common/connector"
)

type ErShouFangPriceHistory struct {
	Id                 uint64 `json:"id"`
	HouseId            string `json:"house_id"`
	TotalPrice         uint64 `db:"total_price" json:"total_price"`
	UnitPrice          uint64 `db:"unit_price" json:"unit_price"`
	UnitPriceText      string `json:"unit_price_text"`
	Status             uint64 `db:"status" json:"status"`
	ScrapeAt           int64  `json:"scrape_at"`
	ScrapeDayTimestamp uint64 `json:"scrape_day_timestamp"`
}

func (history *ErShouFangPriceHistory) Insert(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Create(history).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("插入ErShouFangPriceHistory失败")
		return err
	}
	return nil
}

func (history *ErShouFangPriceHistory) Update(ctx context.Context) error {
	if err := connector.GetMysqlConnector(ctx).Save(history).Error; err != nil {
		log.GlobalLogger().WithError(err).Errorf("更新ErShouFangPriceHistory失败")
		return err
	}
	return nil
}

func (history *ErShouFangPriceHistory) GetTotalPrice() uint64 {
	return history.TotalPrice
}
func (history *ErShouFangPriceHistory) SetTotalPrice(price uint64) {
	history.TotalPrice = price
}
func (history *ErShouFangPriceHistory) UnitPriceString() string {
	return history.UnitPriceText
}
func (history *ErShouFangPriceHistory) GetUnitPrice() uint64 {
	return history.UnitPrice
}
func (history *ErShouFangPriceHistory) SetUnitPrice(unitPrice uint64) {
	history.UnitPrice = unitPrice
}
