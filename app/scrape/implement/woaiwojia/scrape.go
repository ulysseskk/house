package woaiwojia

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ulysseskk/house/app/common/model/statistic"

	"github.com/ulysseskk/house/app/common/model/db"

	"github.com/ulysseskk/house/app/scrape/metrics"

	"github.com/gocolly/colly/v2"
	"github.com/ulysseskk/house/app/common/log"
)

func New5i5jExecutor() *WIWJExecutor {
	return &WIWJExecutor{
		platform: "",
		status:   &statistic.ScrapeStatus{},
		stat:     &statistic.ScrapeStat{},
	}
}

type WIWJExecutor struct {
	platform string
	status   *statistic.ScrapeStatus
	stat     *statistic.ScrapeStat
}

func (e *WIWJExecutor) Status() statistic.ScrapeStatus {
	return *e.status
}
func (e *WIWJExecutor) Report() statistic.ScrapeStat {
	return e.stat.Clone()
}
func (e *WIWJExecutor) SetPlatform(platform string) {
	e.platform = platform
}
func (e *WIWJExecutor) Platform() string {
	return e.platform
}
func (e *WIWJExecutor) GetAllAreaLink(ctx context.Context, cityName string) (map[string]string, error) {
	rootUrl := "https://bj.5i5j.com/ershoufang/"
	e.status.CurrentOperation = "获取地区链接列表"
	areaMap := map[string]string{}
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err := c.Head("https://www.5i5j.com/")
	if err != nil {
		return nil, err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	c.OnRequest(func(request *colly.Request) {
		log.GlobalLogger().WithField("request", request).Infof("访问URL:%+v", request.URL)
	})
	c.OnResponse(func(response *colly.Response) {
		log.GlobalLogger().WithField("request", response.Request).Infof("获取到返回。Status code %d", response.StatusCode)
	})
	c.OnHTML("ul", func(element *colly.HTMLElement) {
		if element.Attr("class") == "new_di_tab sTab" {
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				area := element.ChildText("li")
				href := fmt.Sprintf("%s%s", strings.TrimSuffix(rootUrl, "/ershoufang/"), element.Attr("href"))
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
	metrics.TimeUse.WithLabelValues("woaiwojia").Observe(float64(timeUseMs))
	return areaMap, nil
}
func (e *WIWJExecutor) GetInnerAreaLink(ctx context.Context, areaName, areaLink string) (map[string]string, error) {
	return map[string]string{areaName: areaLink}, nil
}
func (e *WIWJExecutor) ExecuteForPage(ctx context.Context, pageLink string) (string, error) {
	return e.executeForPage(ctx, pageLink, executeChain)
}
func (e *WIWJExecutor) executeForPage(ctx context.Context, pageUrl string, callbacks []executeFunc) (nextPageUrl string, err error) {
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err = c.Head("https://www.5i5j.com/")
	if err != nil {
		return "", err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	// 查页码
	c.OnHTML("div", func(element *colly.HTMLElement) {
		if element.Attr("class") == "pageSty rf" {
			element.ForEach("a", func(i int, element *colly.HTMLElement) {
				if element.Text == "下一页" {
					nextPageUrl = fmt.Sprintf("%s%s", strings.Split(pageUrl, "/ershoufang")[0], element.Attr("href"))
				}
			})
		}
	})
	// 查本页
	c.OnHTML("ul", func(element *colly.HTMLElement) {
		if element.Attr("class") == "pList" {
			element.ForEach("li", func(i int, element *colly.HTMLElement) {
				info := &db.ErShouFangDetail{
					Platform: "woaiwojia",
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
				element.ForEach(".listImg", func(i int, element *colly.HTMLElement) {
					houseIdStruct := &struct {
						HouseidVar string `json:"houseid_var"`
					}{}
					err := json.Unmarshal([]byte(element.Attr("giojson")), houseIdStruct)
					if err != nil {
						return
					}
					info.HouseID = houseIdStruct.HouseidVar
					element.ForEach("a", func(i int, element *colly.HTMLElement) {
						info.Url = fmt.Sprintf("%s%s", strings.Split(pageUrl, "/ershoufang")[0], element.Attr("href"))
					})
					element.ForEach("img", func(i int, element *colly.HTMLElement) {
						info.Title = element.Attr("title")
					})
				})
				element.ForEach(".i_01", func(i int, element *colly.HTMLElement) {
					total := element.DOM.Parent().Text()
					splits := strings.Split(total, "·")
					info.HouseType = strings.TrimSpace(splits[0])
					info.AreaText = strings.TrimSpace(splits[1])
					info.DirectionText = strings.TrimSpace(splits[2])
					info.Floor = strings.TrimSpace(splits[3])
					info.Fixture = strings.TrimSpace(splits[4])
				})
				element.ForEach("a", func(i int, element *colly.HTMLElement) {
					info.XiaoquName = element.Text
				})
				element.ForEach("strong", func(i int, element *colly.HTMLElement) {
					price, err := strconv.Atoi(element.Text)
					if err == nil {
						priceInfo.TotalPrice = uint64(price)
					}
				})
				element.ForEach(".jia", func(i int, element *colly.HTMLElement) {
					element.ForEach("p", func(i int, element *colly.HTMLElement) {
						if element.Attr("class") == "" {
							priceInfo.UnitPriceText = element.Text
						}
					})
				})
				priceInfo.HouseId = info.HouseID
				info.PriceHistory = priceInfo
				e.stat.HouseCount += 1
				for _, callback := range callbacks {
					err := callback(ctx, info)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("执行回调失败.")
					}
				}
			})
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
	metrics.TimeUse.WithLabelValues("woaiwojia").Observe(float64(timeUseMs))
	e.stat.HouseListPageCount += 1
	return nextPageUrl, nil
}
