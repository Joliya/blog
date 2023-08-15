/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 14:45
 * @Description:
 */

package redis

import (
	"blog/pkg/utils"
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

var RedisClient *redis.Client

const ErrRedisNotFound = redis.Nil

type Config struct {
	Addr         string
	Password     string
	DB           int
	MinIdleConn  int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
	EnableTrace  bool
}

func Init(c *Config) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		MinIdleConns: c.MinIdleConn,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		PoolSize:     c.PoolSize,
		PoolTimeout:  c.PoolTimeout,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if utils.IsNotNil(err) {
		panic(err)
	}
	if c.EnableTrace {
		RedisClient.AddHook(redisotel.NewTracingHook())
	}
	return RedisClient
}

func InitTestRedis() {
	mr, err := miniredis.Run()
	if utils.IsNotNil(err) {
		panic(err)
	}
	// 打开下面命令可以测试链接关闭的情况
	// defer mr.Close()
	RedisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	fmt.Println("mini redis addr:", mr.Addr())
}
