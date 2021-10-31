package db

type PositionMetaData struct {
	Name       string `db:"name" json:"name"`
	DataRole   string `db:"data_role" json:"data_role"`
	Platform   string `db:"platform" json:"platform"`
	Href       string `db:"href" json:"href"`
	Level      int32  `db:"level" json:"level"`
	ParentName string `db:"parent_name" json:"parent_name"`
	City       string `db:"city" json:"city"`
	Title      string `db:"title" json:"title"`
}

func (p PositionMetaData) TableName() string {
	return "position_meta_data"
}
