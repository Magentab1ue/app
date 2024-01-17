package configs

type Config struct {
	App      Fiber
	Postgres PostgresSql
	PostgresX PostgresSqlX
	Kafkas   Kafka
	Minio    Minio
	Redis    Redis
}

type Redis struct {
	Host string
	Port string
}

type Minio struct {
	Host      string
	Port      string
	AccessKey string
	SecretKey string
	Name      string
	Secure    string
}

type Fiber struct {
	Port string
}

type PostgresSql struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SslMode      string
}

type PostgresSqlX struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	SslMode      string
}

type Kafka struct {
	Hosts []string
	Group string
}
