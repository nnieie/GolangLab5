package config

type service struct {
	Name string
	Addr string
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type redis struct {
	Addr     string
	Password string
}

type etcd struct {
	Addr string
}

type kafkaConfig struct {
	Brokers []string
	Topic   string
}

type otelConfig struct {
	Endpoint          string
	TraceEnabled      *bool    `mapstructure:"traceEnabled"`
	TraceSampleRatio  *float64 `mapstructure:"traceSampleRatio"`
	RedisTraceEnabled *bool    `mapstructure:"redisTraceEnabled"`
}

type cfR2Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type config struct {
	MySQL   mySQL
	Redis   redis
	Etcd    etcd
	Kafka   kafkaConfig
	OTel    otelConfig `mapstructure:"otel"`
	R2      cfR2Config `mapstructure:"r2"`
	Service service
}
