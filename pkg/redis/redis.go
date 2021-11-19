package redis

import (
	"baseFrame/pkg/config"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

var ctx = context.Background()

type RedisClient struct {
	*redis.Client
}

func InitRedis(cfg *config.Config) (*RedisClient, error) {
	db, _ := strconv.Atoi(cfg.GetConfig("redis", "db"))
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.GetConfig("redis", "addr"),
		Password: cfg.GetConfig("redis", "auth"),
		DB:       db,
	})
	if err := c.Set(ctx, "test", 1, 1).Err(); err != nil {
		return nil, err
	}
	return &RedisClient{c}, nil
}

func (c *RedisClient) Set(key, value string, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Del(key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisClient) DelKeyByPrefix(key string) {
	keys, err := c.Client.Keys(ctx, key+"*").Result()
	if err == nil && len(keys) > 0 {
		c.Client.Del(ctx, keys...).Result()
	}
}

func (c *RedisClient) SetKeyWithExpireAt(key, value string, expireAt time.Time) error {
	expiration := time.Now().AddDate(0, 0, 1).Sub(time.Now())
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Get(key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *RedisClient) GetKeysByPrefix(key string) ([]string, error) {
	return c.Client.Keys(ctx, key+"*").Result()
}

func (c *RedisClient) Verify(key, value string) bool {
	v, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return v == value
}

func (c *RedisClient) IsExist(key string) bool {
	_, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return true
}

func (c *RedisClient) RPush(key, value string) error {
	return c.Client.RPush(ctx, key, value).Err()
}

func (c *RedisClient) GetList(key string) ([]string, error) {
	return c.Client.LRange(ctx, key, 0, -1).Result()
}

func (c *RedisClient) PopList(key string) []string {
	var err error
	var value string
	res := []string{}
	for {
		value, err = c.Client.LPop(ctx, key).Result()
		if err != nil {
			break
		}
		res = append(res, value)
	}
	return res
}

// 向集合添加元素
func (c *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	return c.Client.SAdd(ctx, key, members).Result()
}

// 指定元素是否在集合中
func (c *RedisClient) SIsMember(key string, member interface{}) bool {
	b, _ := c.Client.SIsMember(ctx, key, member).Result()
	return b
}

// 集合中元素列表
func (c *RedisClient) SMembers(key string) ([]string, error) {
	return c.Client.SMembers(ctx, key).Result()
}

// 从集合中删除
func (c *RedisClient) SRem(key string, members ...interface{}) (int64, error) {
	return c.Client.SRem(ctx, key, members).Result()
}
