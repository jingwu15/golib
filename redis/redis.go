package redis

import (
	"fmt"
	"strconv"
	"strings"
	//    "reflect"
	"github.com/gomodule/redigo/redis"
	"github.com/jingwu15/golib/tconv"
	"github.com/jingwu15/golib/time"
	log "github.com/sirupsen/logrus"
)

//常量
const TIMEOUT_CONN int = 1
const TIMEOUT_READ int = 1
const TIMEOUT_WRITE int = 1
const POOL_MAX_IDLE int = 16
const POOL_MAX_ACTIVE int = 16

var configs = map[string]Config{}

type Config struct {
	Ukey               string
	Raw                map[string]interface{}
	ConnMethod         string // connect/sentinel/cluster
	Address            string
	SentinelAddress    []string
	SentinelMasterName string
	ClusterAddress     []string
	Password           string
	Timeout            int
	Db                 int
	Cmdlog             string
	Conn               redis.Conn
	//"address":             "172.16.100.112:6379",
	//"sentinel_address":    "172.16.100.112:26379, 172.16.100.112:26380, 172.16.100.112:26381",
	//"sentinel_mastername": "mymaster",
	//"cluster_address":     "172.16.100.112:26379, 172.16.100.112:26380, 172.16.100.112:26381",
	//"password":            "Gmck7X02",
	//"timeout":             "0",
	//"db":                  "1",
	//"cmdlog":              "off",
}

func ParseConfig(ukey string, cfg map[string]interface{}) (Config, error) {
	config := Config{}
	config.Ukey = ukey
	config.Raw = cfg
	if password, ok := cfg["password"].(string); ok {
		config.Password = password
	}
	if db, ok := cfg["db"].(string); ok {
		dbInt, _ := strconv.Atoi(db)
		config.Db = dbInt
	}
	if cmdlog, ok := cfg["cmdlog"].(string); ok {
		config.Cmdlog = string(cmdlog)
	}
	if address, ok := cfg["address"].(string); ok {
		config.Address = address
		config.ConnMethod = "connect"
	}
	if sentinelAddress, ok := cfg["sentinel_address"].(string); ok {
		config.SentinelAddress = strings.Split(sentinelAddress, ",")
		for i, r := range config.SentinelAddress {
			config.SentinelAddress[i] = strings.TrimSpace(r)
		}
		config.ConnMethod = "sentinel"
	}
	if sentinelMasterName, ok := cfg["sentinel_mastername"].(string); ok {
		config.SentinelMasterName = sentinelMasterName
	}
	if clusterAddress, ok := cfg["cluster_address"].(string); ok {
		config.ClusterAddress = strings.Split(clusterAddress, ",")
		for i, r := range config.ClusterAddress {
			config.ClusterAddress[i] = strings.TrimSpace(r)
		}
		config.ConnMethod = "cluster"
	}

	return config, nil
}

func SetCfg(ukey string, cfg map[string]interface{}) (err error) {
	configs[ukey], err = ParseConfig(ukey, cfg)
	return err
}

func GetClient(ukey string) (c Client, err error) {
	config, ok := configs[ukey]
	if !ok {
		return Client{}, fmt.Errorf("配置不存在")
	}
	if config.ConnMethod == "connect" {
		c, err := Connect(config.Address, config.Password, config.Timeout, config.Db)
		return Client{Conn: c, Address: config.Address, Auth: config.Password, Db: config.Db}, err
	}
	if config.ConnMethod == "sentinel" {
		c, err := ConnectSentinel(config.SentinelMasterName, config.SentinelAddress, config.Password, config.Timeout, config.Db)
		return Client{Conn: c, SentinelMasterName: config.SentinelMasterName, SentinelAddress: config.SentinelAddress, Auth: config.Password, Db: config.Db}, err
	}
	if config.ConnMethod == "cluster" {
		return Client{}, err
	}
	return Client{}, fmt.Errorf("配置错误")
}

func Connect(addr string, password string, db int, timeout int) (c redis.Conn, err error) {
	c, err = redis.DialTimeout("tcp", addr, time.Keep(timeout), time.Keep(timeout), time.Keep(timeout))
	if err != nil {
		return c, err
	}
	if password != "" {
		c.Do("AUTH", password)
	}
	c.Do("select", db)
	return c, nil
}

func ConnectSentinel(masterName string, addrs []string, password string, db int, timeout int) (c redis.Conn, err error) {
	//timeout, timeourRead, timeoutWrite
	cinfo := []string{}
	for _, addr := range addrs {
		conn, err := redis.DialTimeout("tcp", addr, time.Keep(timeout), time.Keep(timeout), time.Keep(timeout))
		if err == nil {
			cinfo, err = redis.Strings(conn.Do("SENTINEL", "get-master-addr-by-name", masterName))
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return c, fmt.Errorf("redigo: no sentinels available, error: %s", err.Error())
	}
	addr := fmt.Sprintf("%s:%s", cinfo[0], cinfo[1])
	c, err = redis.DialTimeout("tcp", addr, time.Keep(timeout), time.Keep(timeout), time.Keep(timeout))
	if err != nil {
		return c, err
	}
	if password != "" {
		c.Do("AUTH", password)
	}
	c.Do("select", db)
	return c, nil
}

//单个客户端连接
type Client struct {
	Conn               redis.Conn
	Address            string
	SentinelAddress    []string
	SentinelMasterName string
	ClusterAddress     []string
	Auth               string
	Db                 int
}

//单个连接池
type Pool struct {
	Pool      redis.Pool
	Addr      string
	Auth      string
	MaxIdle   int
	MaxActive int
}

//-------------------------------------连接相关的操作-------------------------------------------------

func GetTimeout(cfg map[string]string) (timeoutC, timeoutR, timeoutW int) {
	toc, tor, tow := TIMEOUT_CONN, TIMEOUT_READ, TIMEOUT_WRITE
	if v, ok := cfg["timeout"]; ok {
		to, _ := tconv.StrToInt(v)
		toc, tor, tow = to, to, to
	}
	if v, ok := cfg["timeout_conn"]; ok {
		toc, _ = tconv.StrToInt(v)
	}
	if v, ok := cfg["timeout_read"]; ok {
		tor, _ = tconv.StrToInt(v)
	}
	if v, ok := cfg["timeout_write"]; ok {
		tow, _ = tconv.StrToInt(v)
	}
	return toc, tor, tow
}

func GetPoolCfg(cfg map[string]int) map[string]int {
	c := map[string]int{}
	if v, ok := cfg["max_idle"]; ok {
		c["max_idle"] = v
	} else {
		c["max_idle"] = POOL_MAX_IDLE
	}
	if v, ok := cfg["max_active"]; ok {
		c["max_active"] = v
	} else {
		c["max_active"] = POOL_MAX_ACTIVE
	}
	return c
}

func New(cfg map[string]string) (client Client, err error) {
	addr := fmt.Sprintf("%s:%s", cfg["host"], cfg["port"])
	toc, tor, tow := GetTimeout(cfg)

	c, err := redis.DialTimeout("tcp", addr, time.Keep(toc), time.Keep(tor), time.Keep(tow))
	if err != nil {
		return Client{Address: addr}, err
	} else {
		c.Do("AUTH", cfg["auth"])
		return Client{Conn: c, Address: addr, Auth: cfg["auth"]}, nil
	}
}

func (c Client) Close() {
	c.Conn.Close()
}

func NewPool(cfg map[string]string, cfgs ...map[string]int) Pool {
	addr := fmt.Sprintf("%s:%s", cfg["host"], cfg["port"])
	toc, tor, tow := GetTimeout(cfg)

	opt := GetPoolCfg(cfgs[0])
	redisPool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialTimeout("tcp", addr, time.Keep(toc), time.Keep(tor), time.Keep(tow))
			if err != nil {
				log.Error(fmt.Sprintf("redis connect error[%s]\n", err))
			}
			return conn, err
		},
		MaxIdle:   opt["max_idle"],
		MaxActive: opt["max_active"],
	}
	return Pool{Pool: redisPool, Addr: addr, Auth: cfg["auth"]}
}

//从连接池中取一个连接
func (p Pool) Get() Client {
	c := p.Pool.Get()
	c.Do("AUTH", p.Auth)
	return Client{Conn: c, Address: p.Addr, Auth: p.Auth}
}

//关闭连接池中
func (p Pool) Close() {
	p.Pool.Close()
}

//-------------------------------------普通K-V操作-------------------------------------------------
func (c Client) Set(key string, value string) (string, error) {
	return FStr(c.Do("SET", key, value))
}

func (c Client) Get(key string) (string, error) {
	return FStr(c.Do("GET", key))
}

func (c Client) Mget(keys []string) ([]string, error) {
	ks := make([]interface{}, len(keys))
	for i, v := range keys {
		ks[i] = v
	}
	return FStrs(c.Do("MGET", ks...))
}

//-------------------------------------Hash(哈希表)操作-------------------------------------------------
func (c Client) Hset(key string, hkey string, value string) (int, error) {
	return FInt(c.Do("HSET", key, hkey, value))
}

func (c Client) Hgetall(key string) (map[string]string, error) {
	return FMapStrStr(c.Do("HGETALL", key))
}

func (c Client) Hget(key string, hkey string) (string, error) {
	return FStr(c.Do("HGET", key, hkey))
}

func (c Client) Hincrby(key string, hkey string, incr int) (int, error) {
	return c.Hincr(key, hkey, incr)
}

func (c Client) Hincr(key string, hkey string, incr int) (int, error) {
	return FInt(c.Do("HINCRBY", key, hkey, incr))
}

func (c Client) Hincrbyfloat(key string, hkey string, incr float64) (float64, error) {
	return c.HincrFloat(key, hkey, incr)
}

func (c Client) HincrFloat(key string, hkey string, incr float64) (float64, error) {
	return FFloat(c.Do("HINCRBYFLOAT", key, hkey, incr))
}

func (c Client) Hmget(key string, hkeys ...string) (resp map[string]string, err error) {
	ks := []interface{}{key}
	for _, v := range hkeys {
		ks = append(ks, v)
	}
	rs, err := FBytess(c.Do("HMGET", ks...))
	if err != nil {
		return resp, err
	}

	resp = map[string]string{}
	for i, hkey := range hkeys {
		if len(rs[i]) > 0 {
			resp[hkey] = string(rs[i])
		}
	}
	return resp, err
}

func (c Client) Hdel(key string, hkeys ...string) (int, error) {
	ks := []interface{}{key}
	for _, v := range hkeys {
		ks = append(ks, v)
	}
	return FInt(c.Do("HDEL", ks...))
}

//------------------------------------------------List(列表)操作-----------------------------------
func (c Client) Lindex(key string, index string) (string, error) {
	return FStr(c.Do("LINDEX", key, index))
}

func (c Client) Llen(key string) (int, error) {
	return FInt(c.Do("LLEN", key))
}

func (c Client) Lpop(key string) (string, error) {
	return FStr(c.Do("LPOP", key))
}

func (c Client) Rpush(key string, values ...string) (string, error) {
	result, err := c.Do("RPUSH", values)
	return FStr(result, err)
}

//-----------------------------------------redis连接封-----------------------------------------------

func (c Client) Do(cmd string, args ...interface{}) (result interface{}, err error) {
	result, err = c.Conn.Do(cmd, args...)
	if err != nil {
		log.Error(fmt.Sprintf("redis error[%s]\n", err))
	}
	return result, err
}

func Strs2Ifs(keys ...string) []interface{} {
	ks := []interface{}{}
	for _, v := range keys {
		ks = append(ks, v)
	}
	return ks
}

//------------------------------------------数据类型转换---------------------------------------------
func FBytess(result interface{}, err error) ([][]byte, error) {
	if err == nil {
		return redis.ByteSlices(result, err)
	} else {
		return result.([][]byte), err
	}
}

func FBytes(result interface{}, err error) ([]byte, error) {
	if err == nil {
		return redis.Bytes(result, err)
	} else {
		return result.([]byte), err
	}
}

func FInt(result interface{}, err error) (int, error) {
	if err == nil {
		return redis.Int(result, err)
	} else {
		return result.(int), err
	}
}

func FFloat(result interface{}, err error) (float64, error) {
	if err == nil {
		return redis.Float64(result, err)
	} else {
		return result.(float64), err
	}
}

func FStr(result interface{}, err error) (string, error) {
	if err == nil {
		return redis.String(result, err)
	} else {
		return result.(string), err
	}
}
func FStrs(result interface{}, err error) ([]string, error) {
	rows, err := FBytess(result, err)
	rs := []string{}
	for _, row := range rows {
		rs = append(rs, string(row))
	}
	return rs, err
}
func FMapStrStr(result interface{}, err error) (map[string]string, error) {
	if err == nil {
		return redis.StringMap(result, err)
	} else {
		return result.(map[string]string), err
	}
}
