package dal

type DBAccess struct {
	Xiaoqu           XiaoQuDBAccess
	XiaoQuPosition   XiaoQuPositionDBAccess
	RentDetail       RentDetailDBAccess
	RentPriceHistory RentPriceHistoryDBAccess
}

var GlobalDBAccess = &DBAccess{
	Xiaoqu:           XiaoQuDBAccess{},
	XiaoQuPosition:   XiaoQuPositionDBAccess{},
	RentDetail:       RentDetailDBAccess{},
	RentPriceHistory: RentPriceHistoryDBAccess{},
}
