package db

import (
	"context"

	"github.com/abyss414/house/app/common/connector"
)

type Xiaoqu struct {
	ID     uint64 `db:"id" json:"id"`
	Name   string `db:"name" json:"name"`
	LonStr string `json:"lon_str"`
	LatStr string `json:"lat_str"`
	//Lat      string `db:"lat" json:"lat"`
	//Lon      string `db:"lon" json:"lon"`
	RegionID int32  `db:"region_id" json:"region_id"`
	URL      string `db:"url" json:"url"`
	//Gis      gormGIS.GeoPoint `json:"gis" db:"gis"`
}

func (x Xiaoqu) TableName() string {
	return "xiaoqu"
}

func (x *Xiaoqu) Create(ctx context.Context) error {
	return connector.GetMysqlConnector(ctx).Create(x).Error
}
