package data

import (
	"GraduateDesign/conf"
	"github.com/go-redis/redis/v7"
	"log"
)

var Client *redis.Client

// 开启redis连接池
func initRedisConnection(config conf.AppConfig) {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Address,
		Password: config.App.Redis.Password, // It's ok if password is "".
		DB:       0,                         // 默认0号数据库
	})

	if config.App.FlushAllForTest {
		_, err := FlushAll()
		if err != nil {
			log.Println("Redis flushAll时出错：", err.Error())
		}
	}

}

// 用于测试
func FlushAll() (string, error) {
	return Client.FlushAll().Result()
}

// TODO: 确保redis加载lua脚本，若未加载则加载
