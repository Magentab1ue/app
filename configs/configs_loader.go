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
	cfg.Postgres.Host = secret.Data["DB_HOST"].(string)
	logs.Info(fmt.Sprintf("Env DB_HOST : %s", cfg.Postgres.Host))

	cfg.Postgres.Port = secret.Data["DB_PORT"].(string)
	logs.Info(fmt.Sprintf("Env DB_PORT : %s", cfg.Postgres.Port))

	cfg.Postgres.Username = secret.Data["DB_USER"].(string)
	logs.Info(fmt.Sprintf("Env DB_USER : %s", cfg.Postgres.Username))

	cfg.Postgres.Password = secret.Data["DB_PASSWORD"].(string)
	logs.Info(fmt.Sprintf("Env DB_PASSWORD : %s", cfg.Postgres.Port))

	cfg.Postgres.DatabaseName = secret.Data["DB_NAME"].(string)
	logs.Info(fmt.Sprintf("Env DB_NAME : %s", cfg.Postgres.Password))

	cfg.Postgres.SslMode = secret.Data["DB_SSLMODE"].(string)
	logs.Info(fmt.Sprintf("Env DB_SSLMODE : %s", cfg.Postgres.SslMode))

	cfg.Postgres.Schema = secret.Data["DB_SCHEMA"].(string)
	logs.Info(fmt.Sprintf("Env DB_SCHEMA : %s", cfg.Postgres.Schema))

	// Kafka
	cfg.Kafkas.Hosts = []string{secret.Data["KAFKA_SERVER"].(string)}
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.Kafkas.Hosts))

	cfg.Kafkas.Group = secret.Data["KAFKA_GROUP"].(string)
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.Kafkas.Group))

	cfg.Redis.Host = secret.Data["REDIS_HOST"].(string)
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.Redis.Host))

	cfg.Redis.Port = secret.Data["REDIS_PORT"].(string)
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.Redis.Port))

	cfg.Redis.Password = secret.Data["REDIS_PASSWORD"].(string)
	logs.Info(fmt.Sprintf("Env APP_PORT : %s", cfg.Redis.Password))

}
