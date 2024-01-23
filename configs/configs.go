package configs

type Config struct {
	App      Fiber
	Postgres PostgresSql
	Kafkas   Kafka
	Redis    Redis
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type Fiber struct {
	Host string
	Port string
}

type PostgresSql struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SslMode      string
	Schema       string
}

type Kafka struct {
	Servers  []string
	Port     string
	Group    string
	ClientID string
}
