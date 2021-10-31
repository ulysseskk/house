package dal

import (
	"context"

	"gorm.io/gorm"

	"github.com/abyss414/house/app/common/connector"
	"github.com/abyss414/house/app/common/model/db"
)

type XiaoQuDBAccess struct {
}

func (dba *XiaoQuDBAccess) ListByPage(ctx context.Context, offset, pageSize int) ([]*db.Xiaoqu, error) {
	result := []*db.Xiaoqu{}
	if err := connector.GetMysqlConnector(ctx).Model(&db.Xiaoqu{}).Offset(offset).Limit(pageSize).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (dba *XiaoQuDBAccess) GetXiaoquByName(ctx context.Context, name string) (*db.Xiaoqu, error) {
	result := &db.Xiaoqu{}
	if err := connector.GetMysqlConnector(ctx).Model(&db.Xiaoqu{}).Where("name = ?", name).First(result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (dba *XiaoQuDBAccess) Update(ctx context.Context, xiaoqu *db.Xiaoqu) error {
	if err := connector.GetMysqlConnector(ctx).Save(xiaoqu).Error; err != nil {
		return err
	}
	return nil
}

func (dba *XiaoQuDBAccess) ListNoPosition(ctx context.Context) ([]*db.Xiaoqu, error) {
	result := []*db.Xiaoqu{}
	if err := connector.GetMysqlConnector(ctx).Model(&db.Xiaoqu{}).Where("lon_str is null").Where("lat_str is null").Order("id asc").Find(&result).Error; err != nil {
		return nil, err
	}
	emptyResult := []*db.Xiaoqu{}
	if err := connector.GetMysqlConnector(ctx).Model(&db.Xiaoqu{}).Where("lon_str = ''").Where("lat_str  = ''").Order("id asc").Find(&emptyResult).Error; err != nil {
		return nil, err
	}
	for i, _ := range emptyResult {
		result = append(result, emptyResult[i])
	}
	return result, nil
}
