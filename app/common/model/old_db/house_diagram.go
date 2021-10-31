package old_db

type HouseDiagram struct {
	ID             int32  `db:"id" json:"id"`
	Diagram        string `db:"diagram" json:"diagram"`
	HouseHistoryID int32  `db:"house_history_id" json:"house_history_id"`
	HouseID        string `db:"house_id" json:"house_id"`
	Md5sum         string `db:"md5sum" json:"md5sum"`
	Path           string `db:"path" json:"path"`
}

func (h HouseDiagram) TableName() string {
	return "house_diagram"
}
