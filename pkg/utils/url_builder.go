package utils

import (
	"errors"
	"fmt"

	"approval-service/configs"
)

func UrlBuilder(urlType string, cfg *configs.Config) (string, error) {

	var url string

	switch urlType {
	case "fiber":
		url = fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	case "postgres":
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s ",
			cfg.Postgres.Host,
			cfg.Postgres.Port,
			cfg.Postgres.Username,
			cfg.Postgres.Password,
			cfg.Postgres.DatabaseName,
			cfg.Postgres.SslMode,
			cfg.Postgres.Schema,
		)
	case "redis":
		url = fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	default:
		err := fmt.Sprintf("error,url builder Unknown url type: %s", urlType)
		return "", errors.New(err)
	}
	return url, nil
}
