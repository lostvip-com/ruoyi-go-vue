package lv_redis

import (
	"context"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_conf"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	redisClient *RedisClient
	onceRedis   = sync.Once{}
)

type RedisClient struct {
	client *redis.Client
}

func GetInstance(indexDb int) *RedisClient {
	if redisClient == nil {
		onceRedis.Do(func() {
			redisClient = NewRedisClient(indexDb)
		})
	}
	return redisClient
}

func NewRedisClient(indexDb int) *RedisClient {
	conf := lv_conf.CfgDefault{}
	addr := conf.GetValueStr("application.redis.host")
	port := conf.GetValueStr("application.redis.port")
	password := conf.GetValueStr("application.redis.password")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: password, // 没有密码，默认值
		DB:       indexDb,  // 默认DB 0
	})
	redisClient = new(RedisClient)
	redisClient.client = rdb
	if redisClient.client.Ping(context.Background()).Val() == "" {
		msg := ` 
			  ------------>连接 reids 错误：
			  无法链接到redis!!!! 检查相关配置:
			  host: %v
			  port: %v
			  password: %v
             `
		host := conf.GetValueStr("application.redis.host")
		lv_log.Error(fmt.Sprintf(msg, host, conf.GetValueStr("application.redis.port"), conf.GetValueStr("application.redis.password")))
		panic("redis 错误:" + host + " port:" + port)
	}
	return redisClient
}

func (rcc *RedisClient) HMSet(key string, mp map[string]any, expiration time.Duration) error {
	err := rcc.client.HSet(context.Background(), key, mp).Err()
	err = rcc.client.Expire(context.Background(), key, expiration).Err()
	return err
}

func (rcc *RedisClient) Expire(key string, duration time.Duration) error {
	return rcc.client.Expire(context.Background(), key, duration).Err()
}

func (rcc *RedisClient) Exists(key string) (int64, error) {
	return rcc.client.Exists(context.Background(), key).Result()
}

func (rcc *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	rcc.client.Set(context.Background(), key, value, expiration)
	return nil
}

func (rcc *RedisClient) Get(key string) (data string, err error) {
	data, err = rcc.client.Get(context.Background(), key).Result()
	return data, err
}

func (rcc *RedisClient) Del(keys ...string) error {
	var err error = nil
	for _, key := range keys {
		err = rcc.client.Del(context.Background(), key).Err()
	}
	return err
}

func (rcc *RedisClient) HSet(key string, values ...interface{}) error {
	err := rcc.client.HSet(context.Background(), key, values...).Err()
	return err
}

func (rcc *RedisClient) HGet(key, field string) (string, error) {
	data, err := rcc.client.HGet(context.Background(), key, field).Result()
	return data, err
}

func (rcc *RedisClient) HDel(key string, fields ...string) error {
	return rcc.client.HDel(context.Background(), key, fields...).Err()
}

func (rcc *RedisClient) HGetAll(key string) (map[string]string, error) {
	return rcc.client.HGetAll(context.Background(), key).Result()
}

func (rcc *RedisClient) Execute(script string, keys []string, args ...interface{}) (interface{}, error) {
	return rcc.client.Eval(context.TODO(), script, keys, args).Result()
}

func (rcc *RedisClient) Close() {
	err := rcc.client.Close()
	if err != nil {
		return
	}
}
func (r RedisClient) Scan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	return r.client.Scan(context.TODO(), cursor, match, count).Result()
}
