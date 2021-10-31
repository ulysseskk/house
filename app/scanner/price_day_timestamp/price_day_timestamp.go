package price_day_timestamp

import (
	"context"
	"time"

	"github.com/abyss414/house/app/common/connector"
	"github.com/abyss414/house/app/common/log"
	"github.com/abyss414/house/app/common/model/db"
)

func ScanForPriceDayTimestamp(ctx context.Context) error {
	for {
		// 总数
		totalCount := int64(0)
		err := connector.GetMysqlConnector(ctx).Model(&db.ErShouFangPriceHistory{}).Where("scrape_day_timestamp = 0 or scrape_day_timestamp is null").Count(&totalCount).Error
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("获取需要计算的数量失败")
			return err
		}
		// 开始分页
		pageNo := int64(1)
		pageCount := int64(20)
		nowOffset := int64(0)
		for nowOffset < totalCount {
			result := []*db.ErShouFangPriceHistory{}
			err := connector.GetMysqlConnector(ctx).Where("scrape_day_timestamp = 0 or scrape_day_timestamp is null").Offset(int((pageNo - 1) * pageCount)).Limit(int(pageCount)).Find(&result).Error
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("分页获取Price history失败")
				return err
			}
			for i, history := range result {
				nowOffset++
				// 获取当天0点时间戳
				timeObj := time.Unix(history.ScrapeAt, 0)
				dayTime := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 0, 0, 0, 0, timeObj.Location())
				history.ScrapeDayTimestamp = uint64(dayTime.Unix())
				// 更新
				result[i].Update(ctx)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
