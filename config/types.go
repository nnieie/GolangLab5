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

type cfR2Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type config struct {
	MySQL   mySQL
	Redis   redis
	Etcd    etcd
	Service service
}
