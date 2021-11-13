package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ulysseskk/house/app/common/constant"
	"github.com/ulysseskk/house/app/common/tracing"
	"github.com/ulysseskk/house/app/scrape/framework"

	"github.com/ulysseskk/house/app/common/config"
	"github.com/ulysseskk/house/app/common/connector"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Fail to load config.Error %+v", err)
	}
	connector.Init()
	tracing.InitTracer("house_scraper")
	wg := &sync.WaitGroup{}
	if config.GlobalConfig().Scrapper.EnableZiroom {
		wg.Add(1)
		go func() {
			s := framework.NewScraper("", nil, constant.PlatformZiroom)
			s.Start()
			for !s.Finished() {
				time.Sleep(2 * time.Second)
			}
			s.Report().Insert(context.Background())
			wg.Done()
		}()
	}
	if config.GlobalConfig().Scrapper.EnableLianjia {
		wg.Add(1)
		go func() {
			s := framework.NewScraper("", nil, constant.PlatformLianjia)
			s.Start()
			for !s.Finished() {
				time.Sleep(2 * time.Second)
			}
			s.Report().Insert(context.Background())
			wg.Done()
		}()
	}
	if config.GlobalConfig().Scrapper.Enable5I5j {
		wg.Add(1)
		go func() {
			s := framework.NewScraper("", nil, constant.Platform5i5j)
			s.Start()
			for !s.Finished() {
				time.Sleep(2 * time.Second)
			}
			s.Report().Insert(context.Background())
			wg.Done()
		}()
	}
	wg.Wait()
}
