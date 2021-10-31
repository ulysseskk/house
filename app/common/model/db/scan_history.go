package db

type ScanHistory struct {
	ID        int32  `db:"id" json:"id"`
	HouseID   string `db:"house_id" json:"house_id"`
	Result    int32  `db:"result" json:"result"`
	Remark    string `db:"remark" json:"remark"`
	StartTime int64  `db:"start_time" json:"start_time"`
	EndTime   int64  `db:"end_time" json:"end_time"`
}

func (s ScanHistory) TableName() string {
	return "scan_history"
}
