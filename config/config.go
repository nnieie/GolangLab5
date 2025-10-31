package config

import (
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/nnieie/golanglab5/pkg/logger"
)

var (
	Mysql        *mySQL
	Redis        *redis
	Etcd         *etcd
	Service      *service
	CFR2Config   *cfR2Config
	runtimeViper = viper.New()
)

func Init(serviceName string) {
	ReadConfigFile(serviceName)
	LoadR2ConfigFromEnv()
}

func ReadConfigFile(serviceName string) {
	c := new(config)
	runtimeViper.SetConfigName("config")   // 不用写扩展名
	runtimeViper.SetConfigType("yaml")     // 文件类型
	runtimeViper.AddConfigPath("./config") // 在项目的config目录查找
	err := runtimeViper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = runtimeViper.Unmarshal(c)
	if err != nil {
		log.Fatalln(err)
	}

	addr := runtimeViper.GetString("services." + serviceName + ".addr")
	Service = &service{
		Name: runtimeViper.GetString("services." + serviceName + ".name"),
		Addr: addr,
	}

	Mysql = &c.MySQL
	Redis = &c.Redis
	Etcd = &c.Etcd
}

func LoadR2ConfigFromEnv() {
	CFR2Config = &cfR2Config{
		Endpoint:        getEnv("R2_Endpoint", ""),
		AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	logger.Warnf("get key from env err")
	return defaultValue
}
