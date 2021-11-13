package lianjia

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ulysseskk/house/app/common/model/statistic"

	"github.com/ulysseskk/house/app/common/tracing"

	"github.com/ulysseskk/house/app/common/model/db"

	"github.com/ulysseskk/house/app/scrape/metrics"

	"github.com/gocolly/colly/v2"
	"github.com/ulysseskk/house/app/common/log"
)

func NewLianjiaExecutor() *LianjiaExecutor {
	return &LianjiaExecutor{
		platform: "",
		status:   &statistic.ScrapeStatus{},
		stat:     &statistic.ScrapeStat{},
	}
}

type LianjiaExecutor struct {
	platform string
	status   *statistic.ScrapeStatus
	stat     *statistic.ScrapeStat
}

func (e *LianjiaExecutor) SetPlatform(platform string) {
	e.platform = platform
}
func (e *LianjiaExecutor) Platform() string {
	return e.platform
}
func (e *LianjiaExecutor) GetAllAreaLink(parentCtx context.Context, cityName string) (map[string]string, error) {
	rootUrl := "https://bj.lianjia.com/ershoufang/"
	ctx := tracing.StartTracingFromCtx(parentCtx, "获取地区链接列表")
	e.status.CurrentOperation = "获取地区链接列表"
	areaMap := map[string]string{}
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err := c.Head("https://www.lianjia.com/")
	if err != nil {
		return nil, err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	c.OnHTML("div", func(element *colly.HTMLElement) {
		if element.Attr("data-role") == "ershoufang" {
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				area := element.Text
				href := fmt.Sprintf("%s%s", strings.TrimSuffix(rootUrl, "/ershoufang/"), element.Attr("href"))
				log.WithContext(ctx).Infof("")
				areaMap[area] = href
			})
		}
	})
	start := time.Now()
	err = c.Visit(rootUrl)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	timeUseMs := (end.UnixNano() - start.UnixNano()) / 1000000
	metrics.TimeUse.WithLabelValues("lianjia").Observe(float64(timeUseMs))
	return areaMap, nil
}
func (e *LianjiaExecutor) GetInnerAreaLink(ctx context.Context, areaName, areaLink string) (map[string]string, error) {
	return map[string]string{areaName: areaLink}, nil
}

func (e *LianjiaExecutor) ExecuteForPage(parentCtx context.Context, pageUrl string) (nextPageUrl string, err error) {
	return e.executeForPage(parentCtx, pageUrl, executeChain)
}
func (e *LianjiaExecutor) executeForPage(parentCtx context.Context, pageUrl string, callbacks []executeFunc) (nextPageUrl string, err error) {
	ctx := tracing.StartTracingFromCtx(parentCtx, "获取页面房源")
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err = c.Head("https://www.lianjia.com/")
	if err != nil {
		return "", err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	// 查页码
	c.OnHTML("div", func(element *colly.HTMLElement) {
		if element.Attr("comp-module") == "page" {
			dataStr := element.Attr("page-data")
			baseUrl := element.Attr("page-url")
			dataStruct := struct {
				TotalPage int `json:"totalPage"`
				CurPage   int `json:"curPage"`
			}{}
			if dataStr == "" {
				return
			}
			err := json.Unmarshal([]byte(dataStr), &dataStruct)
			if err != nil {
				return
			}
			nextPageNum := dataStruct.CurPage + 1
			if nextPageNum > dataStruct.TotalPage {
				return
			}
			nextPageUrl = fmt.Sprintf("%s%s", strings.Split(pageUrl, "/ershoufang")[0], strings.ReplaceAll(baseUrl, "{page}", strconv.Itoa(nextPageNum)))
		}
	})
	// 查本页
	c.OnHTML("div", func(element *colly.HTMLElement) {
		if element.Attr("class") == "info clear" {
			info := &db.ErShouFangDetail{
				Platform: "lianjia",
				ScrapAt:  time.Now().Unix(),
			}
			priceInfo := &db.ErShouFangPriceHistory{
				Id:         0,
				HouseId:    "",
				TotalPrice: 0,
				UnitPrice:  0,
				Status:     0,
				ScrapeAt:   time.Now().Unix(),
			}
			element.ForEach(".title", func(i int, element *colly.HTMLElement) {
				info.Title = element.Text
				for _, node := range element.DOM.Children().Nodes {
					if node.Data == "a" {
						for _, attribute := range node.Attr {
							if attribute.Key == "href" {
								info.Url = attribute.Val
							}
							if attribute.Key == "data-housecode" {
								info.HouseID = attribute.Val
							}
						}
					}
				}
			})
			element.ForEach(".positionInfo", func(i int, element *colly.HTMLElement) {
				for _, node := range element.DOM.Children().Nodes {
					for _, attribute := range node.Attr {
						if attribute.Val == "region" {
							info.XiaoquName = strings.TrimSpace(node.FirstChild.Data)
						}
					}
				}
			})
			element.ForEach(".houseInfo", func(i int, element *colly.HTMLElement) {
				textInfo := element.Text
				infoStrs := strings.Split(textInfo, "|")
				info.HouseType = strings.TrimSpace(infoStrs[0])
				if info.HouseType == "车位" {
					if len(infoStrs) < 5 {
						return
					}
					info.AreaText = strings.TrimSpace(infoStrs[1])
					info.DirectionText = strings.TrimSpace(infoStrs[2])
					info.Fixture = strings.TrimSpace(infoStrs[3])
					info.BuildingStructure = strings.TrimSpace(infoStrs[4])
				} else {
					if len(infoStrs) < 6 {
						return
					}
					info.AreaText = strings.TrimSpace(infoStrs[1])
					info.DirectionText = strings.TrimSpace(infoStrs[2])
					info.Fixture = strings.TrimSpace(infoStrs[3])
					info.Floor = strings.TrimSpace(infoStrs[4])
					info.BuildingStructure = strings.TrimSpace(infoStrs[5])
				}
			})
			element.ForEach(".priceInfo", func(i int, element *colly.HTMLElement) {
				element.ForEach("span", func(i int, element *colly.HTMLElement) {
					value, exist := element.DOM.Parent().Attr("class")
					if exist {
						// 总价
						if strings.Contains(value, "totalPrice") {
							if num, err := strconv.Atoi(element.Text); err == nil {
								priceInfo.TotalPrice = uint64(num)
							}
						}
						// 单价
						if strings.Contains(value, "unitPrice") {
							priceInfo.UnitPriceText = element.Text
						}
					}

				})
			})
			priceInfo.HouseId = info.HouseID
			info.PriceHistory = priceInfo
			e.stat.HouseCount += 1
			for _, executor := range callbacks {
				err := executor(ctx, info)
				if err != nil {
					log.WithContext(ctx).WithError(err).Errorf("执行回调失败.")
				}
			}
		}
	})
	metrics.VisitCounter.WithLabelValues("ziroom").Add(1)
	start := time.Now()
	err = c.Visit(pageUrl)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Visit %s error", pageUrl)
	}
	end := time.Now()
	timeUseMs := (end.UnixNano() - start.UnixNano()) / 1000000
	metrics.TimeUse.WithLabelValues("lianjia").Observe(float64(timeUseMs))
	e.stat.HouseListPageCount += 1
	return nextPageUrl, nil
}
func (e *LianjiaExecutor) Status() statistic.ScrapeStatus {
	return *e.status
}
func (e *LianjiaExecutor) Report() statistic.ScrapeStat {
	return e.stat.Clone()
}
