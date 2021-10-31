package db

type Region struct {
	ID       int32  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	ParentID int32  `db:"parent_id" json:"parent_id"`
	Type     string `db:"type" json:"type"`
	URL      string `db:"url" json:"url"`
}

func (r Region) TableName() string {
	return "region"
}
