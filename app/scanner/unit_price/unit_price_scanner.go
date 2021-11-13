package unit_price

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ulysseskk/house/app/common/log"

	"github.com/ulysseskk/house/app/common/connector"
	"github.com/ulysseskk/house/app/common/model/db"
)

func ScanForFillUnitPrice(ctx context.Context) error {
	for {
		// 总数
		totalCount := int64(0)
		err := connector.GetMysqlConnector(ctx).Model(&db.ErShouFangPriceHistory{}).Where("unit_price = 0").Where("unit_price_text != ?", "").Count(&totalCount).Error
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("获取需要填充的数量失败")
			return err
		}
		// 开始分页
		pageNo := int64(1)
		pageCount := int64(20)
		nowOffset := int64(0)
		for nowOffset < totalCount {
			result := []*db.ErShouFangPriceHistory{}
			err := connector.GetMysqlConnector(ctx).Where("unit_price = 0").Where("unit_price_text != ?", "").Offset(int((pageNo - 1) * pageCount)).Limit(int(pageCount)).Find(&result).Error
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("分页获取Price history失败")
				return err
			}
			for i, history := range result {
				nowOffset++
				// 先分开元
				splits := strings.Split(history.UnitPriceText, "元")
				priceNumSplits := strings.Split(splits[0], ",")
				priceText := strings.Join(priceNumSplits, "")
				priceUint, err := strconv.ParseUint(priceText, 10, 64)
				if err != nil {
					log.WithContext(ctx).WithError(err).WithField("price_history", history).Errorf("计算价格失败。原数据为%s", priceText)
					continue
				}
				result[i].UnitPrice = priceUint
				// 更新
				result[i].Update(ctx)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}
