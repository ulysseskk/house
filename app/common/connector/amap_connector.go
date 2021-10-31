package connector

import (
	"github.com/abyss414/house/app/common/config"
	"github.com/abyss414/house/app/common/proxy/amap"
)

var AmapClient *amap.AmapClient

func InitAmapClient() {
	AmapClient = amap.NewAmapClient(config.GlobalConfig().Amap.Host, config.GlobalConfig().Amap.Key)
}
