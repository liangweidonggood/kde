package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

//redis 配置信息
type Redis struct {
	DB       int    `yaml:"db"`       //数据库
	Addr     string `yaml:"addr"`     //地址
	Password string `yaml:"password"` //密码
}

var Ctx = context.Background()

/**
生成redis实例
*/
func RedisClient() *redis.Client {
	redisCfg := GVA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	pong, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Println("redis connect ping failed, err:", err)
	}
	log.Println("redis connect ping response:", pong)
	return client
}
