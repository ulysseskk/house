package old_db

type HouseInfo struct {
	Version       int32  `db:"version" json:"version"`
	ID            string `db:"id" json:"id"`
	Title         string `db:"title" json:"title"`
	Href          string `db:"href" json:"href"`
	GoodTags      string `db:"good_tags" json:"good_tags"`
	Xiaoqu        string `db:"xiaoqu" json:"xiaoqu"`
	XiaoquURL     string `db:"xiaoqu_url" json:"xiaoqu_url"`
	Region        string `db:"region" json:"region"`
	RegionURL     string `db:"region_url" json:"region_url"`
	InnerMetaName string `db:"inner_meta_name" json:"inner_meta_name"`
	HouseInfo     string `db:"house_info" json:"house_info"`
	Tags          string `db:"tags" json:"tags"`
	TotalPrice    string `db:"total_price" json:"total_price"`
	UnitPrice     string `db:"unit_price" json:"unit_price"`
	CreatedAt     int64  `db:"created_at" json:"created_at"`
	Consumed      int32  `db:"consumed" json:"consumed"`
}

func (h HouseInfo) TableName() string {
	return "house_info"
}
