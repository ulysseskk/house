package db

type PositionDistance struct {
	ID       int32   `db:"id" json:"id"`
	XiaoquID int32   `db:"xiaoqu_id" json:"xiaoqu_id"`
	Name     string  `db:"name" json:"name"`
	PoiType  int32   `db:"poi_type" json:"poi_type"`
	Distance float64 `db:"distance" json:"distance"`
}

func (p PositionDistance) TableName() string {
	return "position_distance"
}
