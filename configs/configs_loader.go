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
	logs.Info(fmt.Sprintf("Env VAULT_ADDR : %s", config.Address))

	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))
	logs.Info(fmt.Sprintf("Env VAULT_TOKEN : %s", os.Getenv("VAULT_TOKEN")))

	ctx := context.Background()

	secret, err := client.KVv2((os.Getenv("VAULT_TYPE"))).Get(ctx, os.Getenv("VAULT_PATH"))
	logs.Info(fmt.Sprintf("Env VAULT_TYPE : %s", os.Getenv("VAULT_TYPE")))
	logs.Info(fmt.Sprintf("Env VAULT_PATH : %s", os.Getenv("VAULT_PATH")))
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	setData := func(key string) string {
		if val, ok := secret.Data[key]; ok {
			if strVal, ok := val.(string); ok {
				logs.Info(fmt.Sprintf("Env %s : %s", key, strVal))
				return strVal
			}
		}
		logs.Error("Key " + key + " not found or not a string in Vault secret")
		return ""
	}

	//App env
	cfg.App.Host = setData("APP_HOST")
	if cfg.App.Host == "" {
		cfg.App.Host = "0.0.0.0"
	}
	cfg.App.Port = setData("APP_PORT")

	//postgres env
	cfg.Postgres.DatabaseName = setData("DB_DATABASE_APPROVAL")
	cfg.Postgres.Host = setData("DB_HOST")
	cfg.Postgres.Password = setData("DB_PASSWORD_APPROVAL")
	cfg.Postgres.Port = setData("DB_PORT")
	cfg.Postgres.Schema = setData("DB_SCHEMA_APPROVAL")
	cfg.Postgres.SslMode = setData("DB_SSLMOD")
	cfg.Postgres.Username = setData("DB_USERNAME_APPROVAL")

	// Kafka
	cfg.Kafkas.Servers = []string{setData("KAFKA_SERVERS")}
	cfg.Kafkas.Port = setData("KAFKA_PORT")
	cfg.Kafkas.Group = setData("KAFKA_GROUP_ID")
	cfg.Kafkas.ClientID = setData("KAFKA_CLIENT_ID")

	//redis
	cfg.Redis.Host = setData("REDIS_HOST")
	cfg.Redis.Port = setData("REDIS_PORT")
	cfg.Redis.Password = setData("REDIS_PASSWORD")

	//Cors allow
	cfg.Cors.AllowOrigins = setData("AUTH_CORS_ALLOW_ORIGIN")
}
