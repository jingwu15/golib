package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var poolInit = 0
var redisPool redis.Pool

//-------------------------------------普通K-V操作-------------------------------------------------

func Set(key string, value string) (int64, error) {
	args := []interface{}{key, value}
	result, err := redis_do("SET", args...)
	if err == nil {
		return redis.Int64(result, err)
	} else {
		return 0, err
	}
}

func Get(key string) ([]byte, error) {
	var resp []byte
	args := []interface{}{key}
	result, err := redis_do("GET", args...)
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return resp, err
	}
}

//-------------------------------------Hash(哈希表)操作-------------------------------------------------

func Hset(key string, hkey string, value string) (int64, error) {
	args := []interface{}{key, hkey, value}
	result, err := redis_do("HSET", args...)
	if err == nil {
		return redis.Int64(result, err)
	} else {
		return 0, err
	}
}

func Hgetall(key string) (map[string]string, error) {
	args := []interface{}{key}
	result, err := redis_do("HGETALL", args...)
	if err == nil {
		return redis.StringMap(result, err)
	} else {
		return map[string]string{}, err
	}
}

func Hget(key string, hkey string) ([]byte, error) {
	var resp []byte
	args := []interface{}{key, hkey}
	result, err := redis_do("HGET", args...)
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return resp, err
	}
}

//------------------------------------------------List(列表)操作-----------------------------------

func Lindex(key string, index string) ([]byte, error) {
	var resp []byte
	args := []interface{}{key, index}
	result, err := redis_do("LINDEX", args...) 
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return resp, err
	}
}

func Llen(key string) (int64, error) {
	args := []interface{}{key}
	result, err := redis_do("LLEN", args...)
	if err == nil {
		return redis.Int64(result, err)
	} else {
		return 0, err
	}
}

func Lpop(key string) ([]byte, error) {
	var resp []byte
	args := []interface{}{key}
	result, err := redis_do("LPOP", args...)
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return resp, err
	}
}

func Rpush(key string, values ...string) ([]byte, error) {
	var resp []byte
	args := []interface{}{key}
	for _,v := range values {
		args = append(args, v)
	}
	result, err := redis_do("RPUSH", args...)
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return resp, err
	}
}

//-----------------------------------------redis连接封-----------------------------------------------

func redis_do(cmd string, args ...interface{}) (result interface{}, err error) {
	client := GetClient().Get()
	defer client.Close()
	result, err = client.Do(cmd, args...)
	if err != nil {
		log.Error(fmt.Sprintf("redis error[%s]\n", err))
	}
	return result, err
}

func GetClient() *redis.Pool {
	//初始化
	if poolInit == 0 {
		newPool()
		poolInit = 1
	}
	return &redisPool
}

func newPool() {
	redisPool = redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(viper.GetString("redis_api_url"))
			if err != nil {
				log.Error(fmt.Sprintf("redis connect error[%s]\n", err))
			}
			return conn, err
		},
		MaxIdle:   16,
		MaxActive: 16,
	}

	return
}

