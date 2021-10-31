package old_db

type HouseTag struct {
	HouseHistoryID int32 `db:"house_history_id" json:"house_history_id"`
	TagID          int32 `db:"tag_id" json:"tag_id"`
}

func (h HouseTag) TableName() string {
	return "house_tag"
}
