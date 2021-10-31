package framework

import (
	"github.com/abyss414/house/app/common/model/statistic"
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
