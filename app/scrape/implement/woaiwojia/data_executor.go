package woaiwojia

import (
	"context"

	"github.com/abyss414/house/app/common/log"

	"github.com/abyss414/house/app/common/dal"
	"github.com/abyss414/house/app/common/errors"

	"github.com/abyss414/house/app/common/model/db"
)

type executeFunc func(ctx context.Context, detail *db.ErShouFangDetail) error

var executeChain = []executeFunc{TryFillXiaoqu, UpsertRecord}

func TryFillXiaoqu(ctx context.Context, detail *db.ErShouFangDetail) error {
	if detail.HouseID == "" {
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
		detail.XiaoquID = xiaoqu.ID
		return nil
	}
	detail.XiaoquID = xiaoqu.ID
	return nil
}

func UpsertRecord(ctx context.Context, detail *db.ErShouFangDetail) error {
	return detail.Upsert(ctx)
}
