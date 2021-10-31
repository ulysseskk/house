package db

type XiaoquPosition struct {
	ID       int32   `db:"id" json:"id"`
	XiaoquID int32   `db:"xiaoqu_id" json:"xiaoqu_id"`
	PoiType  string  `db:"poi_type" json:"poi_type"`
	Name     string  `db:"name" json:"name"`
	Lon      float64 `db:"lon" json:"lon"`
	Lat      float64 `db:"lat" json:"lat"`
}

func (x XiaoquPosition) TableName() string {
	return "xiaoqu_position"
}
