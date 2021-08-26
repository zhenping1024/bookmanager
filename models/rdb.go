package models

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Rdb *redis.Client

// 初始化连接
func InitClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB

	})

	_, err = Rdb.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
