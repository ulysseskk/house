package old_db

type HouseHistory struct {
	ID                    int32   `db:"id" json:"id"`
	HouseID               string  `db:"house_id" json:"house_id"`
	Title                 string  `db:"title" json:"title"`
	XiaoquID              int32   `db:"xiaoqu_id" json:"xiaoqu_id"`
	TotalPrice            int32   `db:"total_price" json:"total_price"`
	UnitPrice             int32   `db:"unit_price" json:"unit_price"`
	RegionID              int32   `db:"region_id" json:"region_id"`
	Forward               int32   `db:"forward" json:"forward"`
	BuildingArea          float64 `db:"building_area" json:"building_area"`
	HouseType             string  `db:"house_type" json:"house_type"`
	Fixture               string  `db:"fixture" json:"fixture"`
	Warm                  string  `db:"warm" json:"warm"`
	Floor                 string  `db:"floor" json:"floor"`
	BuildingStructure     string  `db:"building_structure" json:"building_structure"`
	BuildingType          string  `db:"building_type" json:"building_type"`
	Elevator              int32   `db:"elevator" json:"elevator"`
	OnBoardAt             string  `db:"on_board_at" json:"on_board_at"`
	LastSale              string  `db:"last_sale" json:"last_sale"`
	KeepYears             string  `db:"keep_years" json:"keep_years"`
	Mortgage              string  `db:"mortgage" json:"mortgage"`
	InnerArea             float64 `db:"inner_area" json:"inner_area"`
	Ownership             string  `db:"ownership" json:"ownership"`
	HouseUsage            string  `db:"house_usage" json:"house_usage"`
	PropertyRight         string  `db:"property_right" json:"property_right"`
	ScrapAt               int64   `db:"scrap_at" json:"scrap_at"`
	Status                int32   `db:"status" json:"status"`
	HouseCredentialStatus string  `db:"house_credential_status" json:"house_credential_status"`
	ElevatorRate          string  `db:"elevator_rate" json:"elevator_rate"`
	HouseStructureType    string  `db:"house_structure_type" json:"house_structure_type"`
	KeySellingPoints      string  `db:"key_selling_points" json:"key_selling_points"`
	XiaoquIntro           string  `db:"xiaoqu_intro" json:"xiaoqu_intro"`
	SellingDetail         string  `db:"selling_detail" json:"selling_detail"`
	TexAnalyse            string  `db:"tex_analyse" json:"tex_analyse"`
	FixtureDesc           string  `db:"fixture_desc" json:"fixture_desc"`
	OwnershipMortgage     string  `db:"ownership_mortgage" json:"ownership_mortgage"`
	DiagramID             int32   `db:"diagram_id" json:"diagram_id"`
}

func (h HouseHistory) TableName() string {
	return "house_history"
}
