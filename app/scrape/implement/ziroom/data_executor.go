package ziroom

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ulysseskk/house/app/common/log"

	"github.com/ulysseskk/house/app/common/dal"

	"github.com/ulysseskk/house/app/common/errors"

	"github.com/ulysseskk/house/app/common/model/db"
)

type executeFunc func(ctx context.Context, detail *db.RentDetail) error

var executeChain = []executeFunc{parseZiroomData, tryFillXiaoqu, upsertRoomData}

func parseZiroomData(ctx context.Context, detail *db.RentDetail) error {
	if detail.HouseId == "" {
		return errors.InvalidDataError
	}
	detail.Url = fmt.Sprintf("https:%s", detail.Url)
	// 先将名称中获取到的信息切分
	splits1 := strings.Split(detail.Name, "·")
	if len(splits1) < 2 {
		return nil
	}
	detail.Type = splits1[0]
	regex1, err := regexp.Compile("[0-9]居室-")
	if err != nil {
		return err
	}
	result1 := regex1.FindString(splits1[1])
	if result1 != "" {
		detail.RoomType = strings.TrimSuffix(result1, "-")
		splits2 := strings.Split(splits1[1], result1)
		detail.XiaoquName = splits2[0]
		detail.Direction = splits2[1]
	}
	if detail.RentPriceHistory != nil {
		detail.RentPriceHistory.HouseId = detail.HouseId
	}
	return nil
}

func tryFillXiaoqu(ctx context.Context, detail *db.RentDetail) error {
	if detail.HouseId == "" {
		return errors.InvalidDataError
	}
	xiaoqu, err := dal.GlobalDBAccess.Xiaoqu.GetXiaoquByName(context.Background(), detail.XiaoquName)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("根据名称获取小区错误")
		return err
	}
	if xiaoqu == nil {
		if detail.XiaoquName == "" {
			return nil
		}
		// 如果不存在，创建小区
		xiaoqu := &db.Xiaoqu{
			Name: detail.XiaoquName,
		}
		err := xiaoqu.Create(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("创建小区错误")
			return nil
		}
		detail.XiaoquId = xiaoqu.ID
		return nil
	}
	detail.XiaoquId = xiaoqu.ID
	return nil
}

func upsertRoomData(ctx context.Context, detail *db.RentDetail) error {
	return detail.UpsertRendDetail(ctx)
}
