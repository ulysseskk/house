package ziroom

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ulysseskk/house/app/common/model/statistic"

	"github.com/ulysseskk/house/app/scrape/metrics"

	"github.com/ulysseskk/house/app/common/connector"

	"github.com/go-resty/resty/v2"

	"golang.org/x/net/html"

	"github.com/ulysseskk/house/app/common/model/db"

	"github.com/gocolly/colly/v2"
	"github.com/ulysseskk/house/app/common/log"
)

func NewZiroomExecutor() *ZiroomExecutor {
	return &ZiroomExecutor{
		platform: "",
		status:   &statistic.ScrapeStatus{},
		stat:     &statistic.ScrapeStat{},
	}
}

type ZiroomExecutor struct {
	platform string
	status   *statistic.ScrapeStatus
	stat     *statistic.ScrapeStat
}

func (e *ZiroomExecutor) SetPlatform(platform string) {
	e.platform = platform
}
func (e *ZiroomExecutor) Platform() string {
	return e.platform
}
func (e *ZiroomExecutor) GetAllAreaLink(ctx context.Context, cityName string) (map[string]string, error) {
	rootUrl := "https://www.ziroom.com/z/"
	log.WithContext(ctx).Infof("开始扫描一级地区链接列表")
	areaMap := map[string]string{}
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err := c.Head("https://www.ziroom.com/")
	if err != nil {
		return nil, err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	c.OnHTML("div.Z_filter", func(element *colly.HTMLElement) {
		element.ForEach("li.f-item", func(i int, element *colly.HTMLElement) {
			if element.ChildText("strong.title") == "找房方式" {
				element.ForEach("div.opt", func(i int, element *colly.HTMLElement) {
					element.ForEach(".opt-type", func(i int, element *colly.HTMLElement) {
						if element.ChildText("span.opt-name") == "区域" {
							element.ForEach("a", func(i int, element *colly.HTMLElement) {
								area := element.Text
								href := fmt.Sprintf("https:%s", element.Attr("href"))
								if strings.Contains(href, "javascript") {
									return
								}
								areaMap[area] = href
							})
						}
					})
				})
			}
		})
	})
	c.OnError(func(response *colly.Response, err error) {
		log.WithContext(ctx).WithError(err).Errorf("扫描一级链接请求失败")
	})
	metrics.VisitCounter.WithLabelValues("ziroom").Add(1)
	start := time.Now()
	err = c.Visit(rootUrl)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	timeUseMs := (end.UnixNano() - start.UnixNano()) / 1000000
	metrics.TimeUse.WithLabelValues("ziroom").Observe(float64(timeUseMs))
	return areaMap, nil
}
func (e *ZiroomExecutor) GetInnerAreaLink(ctx context.Context, areaName, url string) (map[string]string, error) {
	e.status.CurrentOperation = fmt.Sprintf("扫描%s地区的下层链接", areaName)
	innerAreaMap := map[string]string{}
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err := c.Head("https://www.ziroom.com/")
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
	c.OnHTML("div.Z_filter", func(element *colly.HTMLElement) {
		element.ForEach("li.f-item", func(i int, element *colly.HTMLElement) {
			if element.ChildText("strong.title") == "找房方式" {
				element.ForEach("div.opt", func(i int, element *colly.HTMLElement) {
					element.ForEach(".opt-type", func(i int, element *colly.HTMLElement) {
						element.ForEach(".grand-child-opt", func(i int, element *colly.HTMLElement) {
							element.ForEach("a.checkbox ", func(i int, element *colly.HTMLElement) {
								href := fmt.Sprintf("https:%s", element.Attr("href"))
								if strings.Contains(href, "javascript") {
									return
								}
								name := element.Text
								innerAreaMap[name] = href
							})
						})
					})
				})
			}
		})
	})
	c.OnError(func(response *colly.Response, err error) {
		fmt.Println(response.Request.URL, ":", err)
	})
	metrics.VisitCounter.WithLabelValues("ziroom").Add(1)
	start := time.Now()
	err = c.Visit(url)
	if err != nil {
		return nil, err
	}
	end := time.Now()
	timeUseMs := (end.UnixNano() - start.UnixNano()) / 1000000
	metrics.TimeUse.WithLabelValues("ziroom").Observe(float64(timeUseMs))
	return innerAreaMap, nil
}
func (e *ZiroomExecutor) ExecuteForPage(ctx context.Context, pageUrl string) (nextPageUrl string, err error) {
	return e.executeForPage(ctx, pageUrl, executeChain)
}
func (e *ZiroomExecutor) executeForPage(ctx context.Context, pageUrl string, houseCallbacks []executeFunc) (nextPageUrl string, err error) {
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	err = c.Head("https://www.ziroom.com/")
	if err != nil {
		return "", err
	}
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.54"
	c.OnHTML("div.Z_pages", func(element *colly.HTMLElement) {
		nextPageUrl = element.ChildAttr("a.next", "href")
	})
	c.OnHTML("div.Z_list-box", func(element *colly.HTMLElement) {
		element.ForEach("div.item", func(i int, element *colly.HTMLElement) {
			rentDetail := &db.RentDetail{
				Platform: "ziroom",
			}
			log.WithContext(ctx).Debugf("获取到房源")

			// 获取房源url
			element.ForEach(".pic-box", func(i int, element *colly.HTMLElement) {
				element.ForEach("a", func(i int, element *colly.HTMLElement) {
					url := element.Attr("href")
					rentDetail.Url = url
					splits := strings.Split(strings.TrimSuffix(url, ".html"), "/x/")
					if len(splits) > 1 {
						rentDetail.HouseId = splits[1]
						log.WithContext(ctx).Debugf("HouseID为%s", rentDetail.HouseId)
					} else {
						// TODO 报错
					}
				})
			})
			// TODO 查看是否存在，存在则跳过
			currentUrl := ""
			currentDecode := []int{}
			// 获取房源信息
			element.ForEach("div.info-box", func(i int, element *colly.HTMLElement) {
				//title
				titleH5Lable := element.DOM.ChildrenFiltered("h5")
				rentDetail.Name = titleH5Lable.Text()
				element.ForEach("div.desc", func(i int, element *colly.HTMLElement) {
					for _, node := range element.DOM.ChildrenFiltered("div").Nodes {
						if len(node.Attr) == 0 {
							// 为面积
							splits := strings.Split(node.FirstChild.Data, "|")
							rentDetail.Area, err = strconv.ParseFloat(strings.Trim(splits[0], "㎡ "), 10)
							if err != nil {
								log.WithContext(ctx).WithError(err).Debugf("转换面积出错")
							}
							floors := strings.Split(splits[1], "/")
							if len(floors) < 2 {
								continue
							}
							rentDetail.Floor, err = strconv.Atoi(strings.TrimSpace(floors[0]))
							if err != nil {
								// TODO 处理错误
							}
							rentDetail.TotalFloor, err = strconv.Atoi(strings.Trim(floors[1], "层"))
							if err != nil {
								// TODO 处理错误
							}
						}
					}
				})
				log.WithContext(ctx).WithField("info", rentDetail).Debugf("基础信息扫描完毕")
				rentDetail.RentPriceHistory = &db.RentPriceHistory{}
				nodes := element.DOM.ChildrenFiltered("div.price").Nodes
				if len(nodes) == 0 {
					log.WithContext(ctx).WithField("info", rentDetail).WithField("url", pageUrl).Errorf("房源无价格信息")
					return
				}
				actualPriceNode := nodes[0]
				pricePicUrl := ""
				posArray := []int{}
				fc := actualPriceNode.FirstChild
				for {
					if fc.NextSibling == nil {
						break
					}
					if fc.Type != html.ElementNode {
						fc = fc.NextSibling
						continue
					}
					if len(fc.Attr) == 1 {
						rentDetail.RentPriceHistory.Unit = fc.Attr[0].Val
					}
					if len(fc.Attr) <= 1 {
						fc = fc.NextSibling
						continue
					}
					styles := strings.Split(fc.Attr[1].Val, ";")
					pricePicUrl = strings.Trim(strings.Trim(styles[0], "background-image: url("), ")")
					position := strings.Split(styles[1], " -")[1]
					positionFloat, err := strconv.ParseFloat(strings.Trim(position, "px"), 10)
					if err != nil {
						log.WithContext(ctx).WithError(err).WithField("info", rentDetail).WithField("podisiton", positionFloat).WithField("url", pageUrl).Errorf("房源价格offset转换失败")
						return
					}
					pos := positionFloat / 21.4
					posArray = append(posArray, int(pos))
					fc = fc.NextSibling
				}
				log.WithContext(ctx).WithField("pos_array", posArray).Debugf("价格offset扫描完毕。当前pic url为%s", pricePicUrl)
				if currentUrl != pricePicUrl {
					log.WithContext(ctx).Debugf("当前价格url与上次不同，更新点位数字信息")
					currentDecode, err = e.GetSerizalFromPicUrl(pricePicUrl)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("解析点位信息失败")
						return
					}
					currentUrl = pricePicUrl
				}
				priceStr := ""
				for _, posValue := range posArray {
					priceStr = fmt.Sprintf("%s%d", priceStr, currentDecode[posValue])
				}
				prictInt, err := strconv.Atoi(priceStr)
				if err != nil {
					log.WithContext(ctx).WithError(err).WithField("info", rentDetail).WithField("price_str", priceStr).WithField("url", pageUrl).Errorf("房源价格转换失败")
					return
				}
				rentDetail.RentPriceHistory.Price = float64(prictInt)
				rentDetail.RentPriceHistory.PriceStr = fmt.Sprintf("¥%d/月", prictInt)
				rentDetail.RentPriceHistory.ScrapedAt = time.Now()
				log.WithContext(ctx).WithField("info", rentDetail).Debugf("价格信息扫描完毕")
			})
			e.stat.HouseCount += 1
			for _, callback := range houseCallbacks {
				err := callback(ctx, rentDetail)
				if err != nil {
					log.WithContext(ctx).WithError(err).Errorf("执行回调失败.")
				}
			}
		})
	})
	metrics.VisitCounter.WithLabelValues("ziroom").Add(1)
	start := time.Now()
	err = c.Visit(pageUrl)
	if err != nil {
		return "", err
	}
	end := time.Now()
	timeUseMs := (end.UnixNano() - start.UnixNano()) / 1000000
	metrics.TimeUse.WithLabelValues("ziroom").Observe(float64(timeUseMs))
	e.stat.HouseListPageCount += 1
	return "", err
}
func (e *ZiroomExecutor) Status() statistic.ScrapeStatus {
	return *e.status
}
func (e *ZiroomExecutor) Report() statistic.ScrapeStat {
	return e.stat.Clone()
}

var picMap = map[string]string{}

func (e *ZiroomExecutor) GetSerizalFromPicUrl(urlScraped string) ([]int, error) {
	base64Str := ""
	// 验证图片是否存在
	if _, ok := picMap[urlScraped]; ok {
		// 图片已存在
		log.GlobalLogger().WithContext(context.Background()).Debugf("图片已存在，使用现有base64")
		base64Str = picMap[urlScraped]
	} else {
		log.GlobalLogger().WithContext(context.Background()).Debugf("图片不存在，开始下载图片")
		resp, err := resty.New().R().Get(fmt.Sprintf("https:%s", urlScraped))
		if err != nil {
			return nil, err
		}
		log.GlobalLogger().WithContext(context.Background()).Debugf("图片下载完成")
		pic := resp.Body()
		base64Str = base64.StdEncoding.EncodeToString(pic)
		picMap[urlScraped] = base64Str
	}
	log.GlobalLogger().WithContext(context.Background()).Debugf("开始进行OCR识别")
	resultStr, err := connector.OcrConnector.GetByBase64(context.Background(), base64Str)
	if err != nil {
		// TODO 处理ocr连接失败的问题
		return nil, err
	}
	log.GlobalLogger().WithContext(context.Background()).Debugf("OCR识别完成")
	result := []int{}
	for _, r := range []rune(resultStr) {
		num, err := strconv.Atoi(string(r))
		if err != nil {
			continue
		}
		result = append(result, num)
	}
	if len(([]rune(resultStr))) != 10 {
		return nil, fmt.Errorf("数字数量错误，应为10个，实际为%d", len(([]rune(resultStr))))
	}
	return result, nil
}
