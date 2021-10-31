package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/abyss414/house/app/scrape/implement/woaiwojia"

	"github.com/abyss414/house/app/scrape/implement/ziroom"

	"github.com/abyss414/house/app/common/constant"
	"github.com/abyss414/house/app/scrape/implement/lianjia"

	"github.com/abyss414/house/app/common/log"
	"github.com/abyss414/house/app/common/model/statistic"
	"github.com/abyss414/house/app/common/tracing"
)

func NewScraper(city string, conf *Conf, platform string) *Scrapper {
	// 配置检查
	var finalConf *Conf
	if conf == nil {
		finalConf = defaultHouseInfoScrapperConf
	} else {
		finalConf = &Conf{}
		if conf.TimeWaitBetweenArea != 0 {
			finalConf.TimeWaitBetweenArea = conf.TimeWaitBetweenArea
		} else {
			finalConf.TimeWaitBetweenArea = defaultHouseInfoScrapperConf.TimeWaitBetweenArea
		}
		if conf.TimeWaitBetweenPage != 0 {
			finalConf.TimeWaitBetweenPage = conf.TimeWaitBetweenPage
		} else {
			finalConf.TimeWaitBetweenPage = defaultHouseInfoScrapperConf.TimeWaitBetweenPage
		}
	}
	ctx := context.WithValue(context.Background(), "platform", platform)
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	scraper := &Scrapper{
		city:         city,
		Platform:     platform,
		conf:         finalConf,
		ctx:          cancelCtx,
		executor:     initExecutorForPlatform(platform),
		cancelFunc:   cancelFunc,
		runtimeError: nil,
		finished:     false,
	}
	return scraper
}

var defaultHouseInfoScrapperConf = &Conf{
	TimeWaitBetweenArea: 1 * time.Second,
	TimeWaitBetweenPage: 4 * time.Second,
}

type Conf struct {
	TimeWaitBetweenArea time.Duration
	TimeWaitBetweenPage time.Duration
}

type Scrapper struct {
	city         string
	Platform     string
	conf         *Conf
	ctx          context.Context
	cancelFunc   context.CancelFunc
	runtimeError error
	finished     bool
	executor     PlatformExecutor
	startTime    time.Time
	endTime      time.Time
}

func (scraper *Scrapper) Status() statistic.ScrapeStatus {
	result := scraper.executor.Status()
	result.StartAt = scraper.startTime
	return result
}
func (scraper *Scrapper) Report() statistic.ScrapeStat {
	result := scraper.executor.Report()
	result.Platform = scraper.Platform
	result.StartTime = scraper.startTime
	result.EndTime = scraper.endTime
	return result
}

func (scraper *Scrapper) Start() {
	baseCtx := tracing.StartTracingFromCtx(scraper.ctx, fmt.Sprintf("%s_scrapper_main_thread"))
	log.WithContext(baseCtx).Infof("启动Scrapper")

	go func() {
		scraper.startTime = time.Now()
		err := scraper.scanForCity(baseCtx, scraper.city)
		if err != nil {
			scraper.runtimeError = err
			log.GlobalLogger().WithError(err).Errorf("Lianjia Scrapper因错误退出")
		}
		scraper.finished = true
		scraper.cancelFunc()
		scraper.endTime = time.Now()
	}()
}
func (scraper *Scrapper) Error() error {
	return scraper.runtimeError
}

func (scraper *Scrapper) Stop() {

}
func (scraper *Scrapper) Finished() bool {
	return scraper.finished
}
func (scraper *Scrapper) Process() int {
	return 0
}

func (s *Scrapper) scanForCity(parentCtx context.Context, cityName string) error {
	ctx := tracing.StartTracingFromCtx(parentCtx, fmt.Sprintf("%s_scan_for_city_%s", s.Platform, cityName))
	// 第一层，扫描所有的外层链接
	outerCtx := tracing.StartTracingFromCtx(ctx, "获取外层URL")
	outerLinks, err := s.executor.GetAllAreaLink(outerCtx, cityName)
	if err != nil {
		log.WithContext(outerCtx).WithError(err).Errorf("扫描首层链接失败。")
		return err
	}
	// 第二层，扫描所有的内层链接
	innerCtx := tracing.StartTracingFromCtx(parentCtx, "获取内层URL")
	for name, link := range outerLinks {
		log.WithContext(innerCtx).Infof("开始扫描%s地区以获取下层链接", name)
		innerAreaLink, err := s.executor.GetInnerAreaLink(innerCtx, name, link)
		if err != nil {
			log.WithContext(innerCtx).WithError(err).Errorf("扫描二层链接失败。")
			return err
		}
		log.WithContext(innerCtx).Infof("扫描%s地区下层链接完毕，扫描到%d个链接", name, len(innerAreaLink))
		// 第三层，扫描每一页，直到没有下一页链接
		pageCtx := tracing.StartTracingFromCtx(parentCtx, "获取房源")
		for innerAreaName, innerAreaLink := range innerAreaLink {
			log.WithContext(pageCtx).Infof("开始扫描%s地区以获取房源信息", innerAreaName)
			nextPage := innerAreaLink
			for nextPage != "" {
				time.Sleep(s.conf.TimeWaitBetweenPage)
				currentPage := nextPage
				log.WithContext(pageCtx).Infof("开始扫描%s链接", currentPage)
				nextPage, err = s.executor.ExecuteForPage(pageCtx, nextPage)
				if err != nil {
					log.WithContext(pageCtx).WithError(err).Errorf("扫描房源页面%s失败。", currentPage)
					continue
				}
				log.WithContext(pageCtx).Infof("扫描%s链接成功", currentPage)
			}
			time.Sleep(s.conf.TimeWaitBetweenArea)
		}
	}
	return nil
}

type PlatformExecutor interface {
	SetPlatform(platform string)
	Platform() string
	GetAllAreaLink(ctx context.Context, cityName string) (map[string]string, error)
	GetInnerAreaLink(ctx context.Context, areaName, areaLink string) (map[string]string, error)
	ExecuteForPage(ctx context.Context, pageLink string) (string, error)
	ScrapeStatistic
}

func initExecutorForPlatform(platform string) PlatformExecutor {
	switch platform {
	case constant.PlatformLianjia:
		return lianjia.NewLianjiaExecutor()
	case constant.PlatformZiroom:
		return ziroom.NewZiroomExecutor()
	case constant.Platform5i5j:
		return woaiwojia.New5i5jExecutor()
	}
	return nil
}
