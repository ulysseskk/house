package db

type Tag struct {
	ID      int32  `db:"id" json:"id"`
	TagName string `db:"tag_name" json:"tag_name"`
	TagType string `db:"tag_type" json:"tag_type"`
}

func (t Tag) TableName() string {
	return "tag"
}
