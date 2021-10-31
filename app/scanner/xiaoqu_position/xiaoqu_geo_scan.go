package xiaoqu_position

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/abyss414/house/app/common/model/db"

	"github.com/abyss414/house/app/common/connector"

	"github.com/abyss414/house/app/common/dal"
	"github.com/abyss414/house/app/common/log"
)

func ScanForXiaoQuGeoInformation(ctx context.Context) {
	for {
		log.WithContext(ctx).Infof("开始新一轮扫描")
		xiaoQuList, err := dal.GlobalDBAccess.Xiaoqu.ListNoPosition(ctx)
		if err != nil {
			log.GlobalLogger().WithError(err).Errorf("Fail to list no position xiaoqu.Error %+v", err)
			time.Sleep(60 * time.Second)
			continue
		}
		log.WithContext(ctx).Infof("扫描到%d个小区", len(xiaoQuList))
		for i := 0; i < len(xiaoQuList); i++ {
			xiaoqu := xiaoQuList[i]
			result, err := connector.AmapClient.SearchPOI(ctx, []string{strings.Trim(xiaoqu.Name, " ")}, "", "北京", true, 1, 10, "")
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("Fail to search poi for xiaoqu %s", xiaoqu.Name)
				continue
			}
			for _, poi := range result.Pois {
				if poi.PName != "北京市" {
					continue
				}
				switch poi.TypeCode {
				case "120302", "120000", "120200", "120203", "120300", "120301", "120303":
					location := poi.Location
					locations := strings.Split(location, ",")
					if len(locations) < 2 {
						break
					}
					lonStr := locations[0]
					latStr := locations[1]
					xiaoqu.LonStr = lonStr
					xiaoqu.LatStr = latStr
					err = dal.GlobalDBAccess.Xiaoqu.Update(ctx, xiaoqu)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("Fail to update xiaoqu %s", xiaoqu.Name)
					}
				case "990000", "991000", "991001", "991400", "991401", "991500":
					posExist, err := dal.GlobalDBAccess.XiaoQuPosition.GetByXiaoQuIDAndName(ctx, xiaoqu.ID, poi.Name)
					if err != nil {
						log.WithContext(ctx).WithError(err).Errorf("Fail to get exist pos by xiaoqu id %d.name %s", xiaoqu.ID, poi.Name)
						break
					}
					if posExist == nil {
						obj := &db.XiaoquPosition{
							XiaoquID: int32(xiaoqu.ID),
							PoiType:  poi.TypeCode,
							Name:     poi.Name,
							Lon:      0,
							Lat:      0,
						}
						location := poi.Location
						locations := strings.Split(location, ",")
						if len(locations) < 2 {
							break
						}
						lonStr := locations[0]
						latStr := locations[1]
						lonDouble, err := strconv.ParseFloat(lonStr, 10)
						if err != nil {
							log.WithContext(ctx).WithError(err).Errorf("Fail to parse lontitude.Value %s", lonStr)
							break
						}
						latDouble, err := strconv.ParseFloat(latStr, 10)
						if err != nil {
							log.WithContext(ctx).WithError(err).Errorf("Fail to parse latitude.Value %s", lonStr)
							break
						}
						obj.Lon = lonDouble
						obj.Lat = latDouble
						err = dal.GlobalDBAccess.XiaoQuPosition.InsertXiaoQuPosition(ctx, obj)
						if err != nil {
							log.WithContext(ctx).WithError(err).Errorf("Fail to insert xiaoqu.Value %s", *obj)
							break
						}
					}
				}
			}
			if xiaoqu.LonStr == "" && xiaoqu.LatStr == "" {
				xiaoqu.LonStr = "-1"
				xiaoqu.LatStr = "-1"
				err = dal.GlobalDBAccess.Xiaoqu.Update(ctx, xiaoqu)
				if err != nil {
					log.WithContext(ctx).WithError(err).Errorf("Fail to update xiaoqu %s", xiaoqu.Name)
				}
			}
		}
		log.WithContext(ctx).Infof("当前扫描结束。5分钟后继续")
		time.Sleep(5 * time.Minute)
	}
}
