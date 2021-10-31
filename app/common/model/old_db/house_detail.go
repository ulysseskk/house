package old_db

type HouseDetail struct {
	HouseID       string `db:"house_id" json:"house_id"`
	LatestVersion int32  `db:"latest_version" json:"latest_version"`
}

func (h HouseDetail) TableName() string {
	return "house_detail"
}
