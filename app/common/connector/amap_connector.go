package connector

import (
	"github.com/ulysseskk/house/app/common/config"
	"github.com/ulysseskk/house/app/common/proxy/amap"
)

var AmapClient *amap.AmapClient

func InitAmapClient() {
	AmapClient = amap.NewAmapClient(config.GlobalConfig().Amap.Host, config.GlobalConfig().Amap.Key)
}
