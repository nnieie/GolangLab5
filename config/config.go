package config

import (
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

var (
	Mysql        *mySQL
	Redis        *redis
	Etcd         *etcd
	Kafka        *kafkaConfig
	OTel         *otelConfig
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

	// 调试信息：打印正在读取的服务名
	log.Printf("[Config] Loading config for service: %s", serviceName)

	addr := runtimeViper.GetString("services." + serviceName + ".addr")
	name := runtimeViper.GetString("services." + serviceName + ".name")

	// 调试信息：打印读取到的配置
	log.Printf("[Config] Service name: %s, addr: %s", name, addr)

	// 如果地址为空，说明配置读取失败
	if addr == "" {
		log.Printf("[Config] WARNING: Failed to read addr for service '%s'", serviceName)
		log.Printf("[Config] Available services in config: %v", runtimeViper.Get("services"))
	}

	Service = &service{
		Name: name,
		Addr: addr,
	}

	Mysql = &c.MySQL
	Redis = &c.Redis
	Etcd = &c.Etcd
	Kafka = &c.Kafka
	OTel = &c.OTel
	CFR2Config = &c.R2
}

func TelemetryEndpoint() string {
	if OTel != nil && OTel.Endpoint != "" {
		return OTel.Endpoint
	}
	return constants.OpenTelemetryCollectorEndpoint
}

func LoadR2ConfigFromEnv() {
	if CFR2Config == nil {
		CFR2Config = &cfR2Config{}
	}

	CFR2Config = &cfR2Config{
		Endpoint:        getEnvOrFallback("R2_Endpoint", CFR2Config.Endpoint),
		AccessKeyID:     getEnvOrFallback("R2_ACCESS_KEY_ID", CFR2Config.AccessKeyID),
		SecretAccessKey: getEnvOrFallback("R2_SECRET_ACCESS_KEY", CFR2Config.SecretAccessKey),
	}
}

func getEnvOrFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if fallback == "" {
		logger.Warnf("missing optional env %s, keeping config file value", key)
	}
	return fallback
}
