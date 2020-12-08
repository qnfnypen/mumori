package opredis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: viper.GetString("Driver.Redis.Addr"),
		Password: viper.GetString("Driver.Redis.Password"),
		DB: viper.GetInt("Driver.Redis.DB"),
	})

	_,err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatal().Str("error",err.Error()).Msg("连接Redis数据库失败")
	}
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	rdb.Close()
}