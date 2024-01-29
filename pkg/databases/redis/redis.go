package databases

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"approval-service/configs"
	"approval-service/logs"
	"approval-service/pkg/utils"
)

func NewRedisClient(cfg *configs.Config) *redis.Client {
	url, err := utils.UrlBuilder("redis", cfg)
	logs.Info(fmt.Sprintf("Connent redis with server %s", url))
	if err != nil {
		logs.Error(zap.Error(err))
	}
	return redis.NewClient(&redis.Options{
		Addr:     url,
		Password: cfg.Redis.Password,
	})
}
