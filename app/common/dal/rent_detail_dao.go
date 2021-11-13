package dal

import (
	"context"

	"github.com/ulysseskk/house/app/common/connector"
	"github.com/ulysseskk/house/app/common/model/db"
	"gorm.io/gorm"
)

type RentDetailDBAccess struct {
}

func (RentDetailDBAccess) GetRentDetailByHouseId(ctx context.Context, houseId string) (*db.RentDetail, error) {
	result := &db.RentDetail{}
	if err := connector.GetMysqlConnector(ctx).Where("house_id = ?", houseId).First(result).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
