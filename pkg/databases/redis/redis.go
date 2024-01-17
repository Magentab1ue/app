package databases

import (
	"approval-service/configs"
	"approval-service/logs"
	"approval-service/pkg/utils"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func NewRedisClient(cfg *configs.Config) *redis.Client {
	url, err := utils.UrlBuilder("redis", cfg)
	if err != nil {
		logs.Error(zap.Error(err))
	}
	return redis.NewClient(&redis.Options{
		Addr: url,
	})
}
