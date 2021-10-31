package statistic

import (
	"context"
	"time"

	"github.com/abyss414/house/app/common/connector"
)

type ScrapeStat struct {
	Platform             string         `json:"platform"`
	Type                 string         `json:"type"`
	StartTime            time.Time      `json:"start_time"`
	EndTime              time.Time      `json:"end_time"`
	HouseCount           int            `json:"house_count"`
	ErrCount             int            `json:"err_count"`
	ErrReasonTop         map[string]int `json:"err_reason_top" gorm:"-"`
	TotalPageCount       int            `json:"total_page_count"`
	HouseListPageCount   int            `json:"house_list_page_count"`
	HouseDetailPageCount int            `json:"house_detail_page_count"`
}

func (stat *ScrapeStat) Clone() ScrapeStat {
	return ScrapeStat{
		Platform:             stat.Platform,
		Type:                 stat.Type,
		StartTime:            stat.StartTime,
		EndTime:              stat.EndTime,
		HouseCount:           stat.HouseCount,
		ErrCount:             stat.ErrCount,
		ErrReasonTop:         stat.ErrReasonTop,
		TotalPageCount:       stat.TotalPageCount,
		HouseListPageCount:   stat.HouseListPageCount,
		HouseDetailPageCount: stat.HouseDetailPageCount,
	}
}

func (stat ScrapeStat) Insert(ctx context.Context) error {
	return connector.GetMysqlConnector(ctx).Create(&stat).Error
}

type ScrapeStatus struct {
	Platform         string       `json:"platform"`
	StartAt          time.Time    `json:"start_at"`
	CurrentOperation string       `json:"current_operation"`
	RecentOperations []*Operation `json:"recent_operations"`
}

type Operation struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Type   string    `json:"type"`
	Detail string    `json:"detail"`
}
