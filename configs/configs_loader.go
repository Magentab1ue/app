package configs

import (
	"approval-service/logs"
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func LoadConfigs(cfg *Config) {

	config := vault.DefaultConfig()

	config.Address = os.Getenv("VAULT_ADDR")
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	ctx := context.Background()

	secret, err := client.KVv2((os.Getenv("VAULT_TYPE"))).Get(ctx, os.Getenv("VAULT_PATH"))
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	//App env
	cfg.App.Port = secret.Data["APP_PORT"].(string)
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.App.Port))
	//postgres env
	cfg.Postgres.DatabaseName = secret.Data["DB_DATABASE_APPROVAL"].(string)
	logs.Info(fmt.Sprintf("Env DB_DATABASE_APPROVAL : %s", cfg.Postgres.DatabaseName))

	cfg.Postgres.Host = secret.Data["DB_HOST"].(string)
	logs.Info(fmt.Sprintf("Env DB_HOST : %s", cfg.Postgres.Host))

	cfg.Postgres.Password = secret.Data["DB_PASSWORD_APPROVAL"].(string)
	logs.Info(fmt.Sprintf("Env DB_PASSWORD_APPROVAL : %s", cfg.Postgres.Password))

	cfg.Postgres.Port = secret.Data["DB_PORT"].(string)
	logs.Info(fmt.Sprintf("Env DB_PORT : %s", cfg.Postgres.Port))

	cfg.Postgres.Schema = secret.Data["DB_SCHEMA_APPROVAL"].(string)
	logs.Info(fmt.Sprintf("Env DB_SCHEMA_APPROVAL : %s", cfg.Postgres.Schema))

	cfg.Postgres.SslMode = secret.Data["DB_SSLMODE"].(string)
	logs.Info(fmt.Sprintf("Env DB_SSLMODE : %s", cfg.Postgres.SslMode))

	cfg.Postgres.Username = secret.Data["DB_USERNAME_APPROVAL"].(string)
	logs.Info(fmt.Sprintf("Env DB_USERNAME_APPROVAL : %s", cfg.Postgres.Username))

	// Kafka
	cfg.Kafkas.Hosts = []string{secret.Data["KAFKA_SERVERS"].(string)}
	logs.Info(fmt.Sprintf("Env KAFKA_SERVERS : %s", cfg.Kafkas.Hosts))

	cfg.Kafkas.Group = secret.Data["KAFKA_GROUP_ID"].(string)
	logs.Info(fmt.Sprintf("Env KAFKA_GROUP_ID : %s", cfg.Kafkas.Group))

	cfg.Kafkas.Group = secret.Data["KAFKA_CLIENT_ID"].(string)
	logs.Info(fmt.Sprintf("Env KAFKA_CLIENT_ID : %s", cfg.Kafkas.ClientID))


	//redis
	cfg.Redis.Host = secret.Data["REDIS_HOST"].(string)
	logs.Info(fmt.Sprintf("Env REDIS_HOST : %s", cfg.Redis.Host))

	cfg.Redis.Port = secret.Data["REDIS_PORT"].(string)
	logs.Info(fmt.Sprintf("Env REDIS_PORT : %s", cfg.Redis.Port))

	cfg.Redis.Password = secret.Data["REDIS_PASSWORD"].(string)
	logs.Info(fmt.Sprintf("Env REDIS_PASSWORD : %s", cfg.Redis.Password))

}
