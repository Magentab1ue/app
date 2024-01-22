package configs

import (
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

	//postgres env
	cfg.Postgres.Host = secret.Data["DB_HOST"].(string)
	cfg.Postgres.Port = secret.Data["DB_PORT"].(string)
	cfg.Postgres.Username = secret.Data["DB_USER"].(string)
	cfg.Postgres.Password = secret.Data["DB_PASSWORD"].(string)
	cfg.Postgres.DatabaseName = secret.Data["DB_NAME"].(string)
	cfg.Postgres.SslMode = secret.Data["DB_SSLMODE"].(string)
	cfg.Postgres.Schema = secret.Data["DB_SCHEMA"].(string)

	// Kafka
	cfg.Kafkas.Hosts = []string{secret.Data["KAFKA_SERVER"].(string)}
	cfg.Kafkas.Group = secret.Data["KAFKA_GROUP"].(string)

	cfg.Redis.Host = secret.Data["REDIS_HOST"].(string)
	cfg.Redis.Port = secret.Data["REDIS_PORT"].(string)
	cfg.Redis.Password = secret.Data["REDIS_PASSWORD"].(string)

	fmt.Printf("%s\n\n", cfg.Redis.Password)
}
