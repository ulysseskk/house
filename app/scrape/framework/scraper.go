package framework

import (
	"github.com/ulysseskk/house/app/common/model/statistic"
)

type ScrapeStatistic interface {
	Status() statistic.ScrapeStatus
	Report() statistic.ScrapeStat
}

type LifeCycle interface {
	Start()
	Stop()
	Finished() bool
	Process() int
	Error() error
}
