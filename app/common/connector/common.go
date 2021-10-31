package connector

func Init() {
	InitMysql()
	InitAmapClient()
	InitOcrConnector()
	InitGlobalRedisClient()
}
