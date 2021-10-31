package db

type PositionCount struct {
	ID             int32 `db:"id" json:"id"`
	XiaoquID       int32 `db:"xiaoqu_id" json:"xiaoqu_id"`
	DistanceWithin int32 `db:"distance_within" json:"distance_within"`
	PoiType        int32 `db:"poi_type" json:"poi_type"`
	Count          int32 `db:"count" json:"count"`
}

func (p PositionCount) TableName() string {
	return "position_count"
}
