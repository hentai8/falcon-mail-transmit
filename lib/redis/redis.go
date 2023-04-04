package redis

import (
	"context"
	"errors"
	"falcon-mail-transmit/interval/config"
	"falcon-mail-transmit/lib/log"
	"github.com/go-redis/redis/v8"
	"os"
	"sync"
	"time"
)

var client *redis.Client
var once sync.Once

func GetInstance() *redis.Client {

	once.Do(func() {
		client = create()
	})
	return client
}

func create() *redis.Client {
	cfg, err := config.Load()

	if err != nil {
		log.Logger.Fatal("failed to load application configuration: ", err.Error())
		os.Exit(-1)
	}
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
		DB:   cfg.Redis.DB,
	})
	return client
}

func Ping() (err error) {
	client = GetInstance()
	if client.Ping(context.Background()).Val() != "PONG" {
		err = errors.New("failed to ping redis")
		return
	}
	return nil
}

func Set(key string, value interface{}, expire time.Duration) (err error) {
	client = GetInstance()
	_, err = client.Set(context.Background(), key, value, expire).Result()
	if err == nil {
		return
	}
	return err
}

func HSet(key string, field interface{}, value interface{}) (err error) {
	client = GetInstance()
	_, err = client.HSet(context.Background(), key, field, value).Result()
	if err == nil {
		return
	}
	return err
}

func HGetAll(key string) map[string]string {
	client = GetInstance()
	s, err := client.HGetAll(context.Background(), key).Result()
	if err != nil {
		return map[string]string{}
	}
	return s
}

func HDel(key string, field string) (err error) {
	client = GetInstance()
	_, err = client.HDel(context.Background(), key, field).Result()
	if err == nil {
		return
	}
	return err
}

func SetNX(key string, value interface{}, expire time.Duration) (result bool, err error) {
	client = GetInstance()
	resp := client.SetNX(context.Background(), key, value, expire)
	result, err = resp.Result()
	if result == true {
		return true, nil
	}
	return false, err
}

func DelUniqueTokenIfExist(keys []string, values []interface{}) (result int, err error) {
	ctx := context.Background()
	client = GetInstance()
	var script = redis.NewScript(`
local key = KEYS[1]
local unique_token = ARGV[1]
local value = redis.call("GET", key)
if not value then
  return 0
end

if value == unique_token then
	redis.call("DEL", key)
	return 1
end

return 0  


`)
	result, err = script.Run(ctx, client, keys, values...).Int()
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}

func Get(key string) string {
	client = GetInstance()
	s, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}
	return s
}

func Del(key string) error {
	client = GetInstance()
	_, err := client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func Keys(pattern string) []string {
	client = GetInstance()
	s, err := client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return []string{}
	}
	return s
}
