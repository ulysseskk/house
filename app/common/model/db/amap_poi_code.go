package db

type AmapPoiCode struct {
	ID             int32  `db:"id" json:"id"`
	TypeCode       int32  `db:"type_code" json:"type_code"`
	BigCategory    string `db:"big_category" json:"big_category"`
	MidCategory    string `db:"mid_category" json:"mid_category"`
	TinyCategory   string `db:"tiny_category" json:"tiny_category"`
	BigCategoryEn  string `db:"big_category_en" json:"big_category_en"`
	MidCategoryEn  string `db:"mid_category_en" json:"mid_category_en"`
	TinyCategoryEn string `db:"tiny_category_en" json:"tiny_category_en"`
}

func (a AmapPoiCode) TableName() string {
	return "amap_poi_code"
}
