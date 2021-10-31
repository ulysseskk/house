package old_db

type HouseStructure struct {
	ID             int32   `db:"id" json:"id"`
	HouseDiagramID int32   `db:"house_diagram_id" json:"house_diagram_id"`
	Room           string  `db:"room" json:"room"`
	Area           float64 `db:"area" json:"area"`
	Forward        int32   `db:"forward" json:"forward"`
	WindowType     string  `db:"window_type" json:"window_type"`
}

func (h HouseStructure) TableName() string {
	return "house_structure"
}
