package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/ulysseskk/house/app/common/config"
	"github.com/ulysseskk/house/app/common/connector"
	customLog "github.com/ulysseskk/house/app/common/log"
	"github.com/ulysseskk/house/app/common/tracing"
	"github.com/ulysseskk/house/app/scanner/price_day_timestamp"
	"github.com/ulysseskk/house/app/scanner/unit_price"
	"github.com/ulysseskk/house/app/scanner/xiaoqu_position"
)

func main() {
	flag.Parse()
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Fail to load config.Error %+v", err)
	}
	err = customLog.InitLogger()
	if err != nil {
		panic(err)
	}
	connector.Init()
	tracing.InitTracer("house_scanner")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		parentCtx := context.WithValue(context.Background(), "method", "ScanForFillUnitPrice")
		ctx := tracing.StartTracingFromCtx(parentCtx, "ScanForFillUnitPrice")
		unit_price.ScanForFillUnitPrice(ctx)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		parentCtx := context.WithValue(context.Background(), "method", "ScanForPriceDayTimestamp")
		ctx := tracing.StartTracingFromCtx(parentCtx, "ScanForPriceDayTimestamp")
		price_day_timestamp.ScanForPriceDayTimestamp(ctx)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		parentCtx := context.WithValue(context.Background(), "method", "ScanForXiaoQuGeoInformation")
		ctx := tracing.StartTracingFromCtx(parentCtx, "ScanForXiaoQuGeoInformation")
		xiaoqu_position.ScanForXiaoQuGeoInformation(ctx)
		wg.Done()
	}()
	wg.Wait()
}
