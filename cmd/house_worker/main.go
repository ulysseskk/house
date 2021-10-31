package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/abyss414/house/app/common/constant"
	"github.com/abyss414/house/app/common/tracing"
	"github.com/abyss414/house/app/scrape/framework"

	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/connector"
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
