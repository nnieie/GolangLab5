package db

import (
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

var DB *gorm.DB

func initMysqlDSN() string {
	return strings.Join([]string{config.Mysql.Username, ":", config.Mysql.Password,
		"@tcp(", config.Mysql.Addr, ")/", config.Mysql.Database,
		"?charset=", config.Mysql.Charset, "&parseTime=True"}, "")
}

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(initMysqlDSN()))
	if err != nil {
		logger.Fatalf("mysql connect error: %v", err)
	}
	logger.Infof("mysql connected")
	DB = DB.Table(constants.FollowsTableName)
}
