package connector

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ulysseskk/house/app/common/config"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var globalDB *gorm.DB

func InitMysql() {
	var err error
	mysqlConf := config.GlobalConfig().Mysql
	if mysqlConf == nil {
		log.Fatalf("Mysql 配置为空！")
	}
	globalDB, err = InitMysqlConnector(config.GlobalConfig().Mysql.Host, config.GlobalConfig().Mysql.Port, config.GlobalConfig().Mysql.User, config.GlobalConfig().Mysql.Password, config.GlobalConfig().Mysql.DbName)
	if err != nil {
		log.Fatalf("初始化mysql失败。错误%+v", err)
	}
}

func InitMysqlConnector(ip string, port int, username, password, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?", username, password, ip, port, dbName)
	db, err := gorm.Open(mysql.Dialector{
		Config: &mysql.Config{
			DSN: dsn,
		},
	}, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		FullSaveAssociations:                     false,
		Logger:                                   nil,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		Plugins:                                  nil,
	})
	if err != nil {
		return nil, err
	}
	db = db.Set("gorm:save_associations", false).Set("gorm:association_save_reference", false)
	mysqlConn, err := db.DB()
	if err != nil {
		return nil, err
	}
	mysqlConn.SetMaxOpenConns(4)
	mysqlConn.SetMaxIdleConns(3)
	mysqlConn.SetConnMaxLifetime(1 * time.Hour)
	db.Logger.LogMode(logger.Info)
	return db, nil
}
func GetMysqlConnector(ctx context.Context) *gorm.DB {
	return globalDB.WithContext(ctx)
}
