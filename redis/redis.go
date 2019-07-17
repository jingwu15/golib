package redis

import (
	"fmt"
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

//单个客户端连接
type Client struct {
	Conn redis.Conn
	Addr string
	Auth string
}

//单个连接池
type Pool struct {
	Pool      redis.Pool
	Addr      string
	Auth      string
	MaxIdle   int
	MaxActive int
}

//哨兵, 主从
type Sentinel struct {
	Master       string         //master name
	Addrs        []string       //哨兵地址
	CAddrs       map[string]int //当前正在连接的客户端          {"127.0.0.1:6379": 0}
	Auth         string         //密码
	TimeoutConn  int            //连接超时
	TimeoutRead  int            //读取超提
	TimeoutWrite int            //写入超时
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
		return Client{Addr: addr}, err
	} else {
		c.Do("AUTH", cfg["auth"])
		return Client{Conn: c, Addr: addr, Auth: cfg["auth"]}, nil
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
	return Client{Conn: c, Addr: p.Addr, Auth: p.Auth}
}

//关闭连接池中
func (p Pool) Close() {
	p.Pool.Close()
}

//新建哨钱
func NewSentinel(master string, items []map[string]string, cfgs ...map[string]interface{}) Sentinel {
	addrs := []string{}
	for _, row := range items {
		addrs = append(addrs, fmt.Sprintf("%s:%s", row["host"], row["port"]))
	}
	s := Sentinel{Master: master, Addrs: addrs, Auth: "", TimeoutConn: TIMEOUT_CONN, TimeoutRead: TIMEOUT_READ, TimeoutWrite: TIMEOUT_WRITE}
	if v, ok := cfgs[0]["auth"]; ok {
		s.Auth = v.(string)
	}
	if v, ok := cfgs[0]["timeout_conn"]; ok {
		s.TimeoutConn = v.(int)
	}
	if v, ok := cfgs[0]["timeout_read"]; ok {
		s.TimeoutRead = v.(int)
	}
	if v, ok := cfgs[0]["timeout_write"]; ok {
		s.TimeoutWrite = v.(int)
	}
	return s
}

//哨兵取一个连接，未来实现轮询，权重，Hash方式
func (s Sentinel) Get() (Client, error) {
	var err error
	cinfo := []string{}
	for _, addr := range s.Addrs {
		conn, err := redis.DialTimeout("tcp", addr, time.Keep(s.TimeoutConn), time.Keep(s.TimeoutRead), time.Keep(s.TimeoutWrite))
		if err == nil {
			cinfo, err = redis.Strings(conn.Do("SENTINEL", "get-master-addr-by-name", s.Master))
			if err == nil {
				break
			} else {
				fmt.Println("cinfo------fail", err, cinfo)
			}
		}
	}
	if err != nil {
		return Client{}, fmt.Errorf("redigo: no sentinels available, error: %s", err.Error())
	}
	addr := fmt.Sprintf("%s:%s", cinfo[0], cinfo[1])
	c, err := redis.DialTimeout("tcp", addr, time.Keep(s.TimeoutConn), time.Keep(s.TimeoutRead), time.Keep(s.TimeoutWrite))
	if err != nil {
		return Client{}, err
	}
	if s.Auth != "" {
		c.Do("AUTH", s.Auth)
	}
	return Client{Conn: c, Addr: addr}, nil
}

////哨兵取一个连接池, 未来实现轮询，权重，Hash方式，暂不实现
//func (s Sentinel) GetPool() (Pool, error) {
//}

//-------------------------------------普通K-V操作-------------------------------------------------
func (c Client) Set(key string, value string) (string, error) {
	return FStr(c.do("SET", key, value))
}

func (c Client) Get(key string) (string, error) {
	return FStr(c.do("GET", key))
}

func (c Client) Mget(keys []string) ([]string, error) {
    ks := make([]interface{}, len(keys))
    for i, v := range keys { ks[i] = v }
    return FStrs(c.do("MGET", ks...))
}

//-------------------------------------Hash(哈希表)操作-------------------------------------------------
func (c Client) Hset(key string, hkey string, value string) (int, error) {
	return FInt(c.do("HSET", key, hkey, value))
}

func (c Client) Hgetall(key string) (map[string]string, error) {
	return FMapStrStr(c.do("HGETALL", key))
}

func (c Client) Hget(key string, hkey string) (string, error) {
	return FStr(c.do("HGET", key, hkey))
}

func (c Client) Hmget(key string, hkeys []string) (resp map[string]string, err error) {
    ks := []interface{}{key}
    for _, v := range hkeys { ks = append(ks, v) }
    rs, err := FBytess(c.do("HMGET", ks...))
    if err != nil { return resp, err }

    resp = map[string]string{}
    for i, hkey := range hkeys {
        if len(rs[i]) > 0 {
            resp[hkey] = string(rs[i])
        }
    }
    return resp, err
}

//------------------------------------------------List(列表)操作-----------------------------------
func (c Client) Lindex(key string, index string) (string, error) {
	return FStr(c.do("LINDEX", key, index))
}

func (c Client) Llen(key string) (int, error) {
	return FInt(c.do("LLEN", key))
}

func (c Client) Lpop(key string) (string, error) {
	return FStr(c.do("LPOP", key))
}

func (c Client) Rpush(key string, values ...string) (string, error) {
	//var resp []byte
	//args := []interface{}{key}
	//for _,v := range values {
	//	args = append(args, v)
	//}
	result, err := c.do("RPUSH", values)
	return FStr(result, err)
}

//-----------------------------------------redis连接封-----------------------------------------------

func (c Client) do(cmd string, args ...interface{}) (result interface{}, err error) {
	result, err = c.Conn.Do(cmd, args...)
	if err != nil {
		log.Error(fmt.Sprintf("redis error[%s]\n", err))
	}
	return result, err
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
    for _, row := range rows { rs = append(rs, string(row)) }
    return rs, err
}
func FMapStrStr(result interface{}, err error) (map[string]string, error) {
	if err == nil {
		return redis.StringMap(result, err)
	} else {
		return result.(map[string]string), err
	}
}
