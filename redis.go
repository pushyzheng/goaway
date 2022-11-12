package main

import (
	"fmt"
	"github.com/go-redis/redis"
	logger "github.com/sirupsen/logrus"
)

var redisConn *redis.Client

func initRedis() {
	fmt.Println(Conf.Server.Statistics)
	if Conf.Server.Statistics {
		redisConn = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
		})
		_, err := redisConn.Ping().Result()
		if err != nil {
			panic("connect redis error: " + err.Error())
		}
		logger.Info("init redis succeed")
	}
}
