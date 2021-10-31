package db

type Position struct {
	ID       int32  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	TypeID   int32  `db:"type_id" json:"type_id"`
	ParentID int32  `db:"parent_id" json:"parent_id"`
	Gis      string `db:"gis" json:"gis"`
}

func (p Position) TableName() string {
	return "position"
}
