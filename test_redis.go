package main

import (
	"fmt"
	"github.com/jingwu15/golib/redis"
)

//;;[demo]
//;;;; [可选]redis节点连接地址, 格式：<host>:<port>, 不支持多地址
//;;address = 172.16.100.26:6379
//;;;; [可选]redis sentinel节点连接地址, 格式：<host>:<port>, 多以逗号(,)分隔
//;;sentinel_mastername = mymaster
//;;sentinel_address = 172.16.100.26:26379, 172.16.100.26:26380
//;;;; [可选]redis cluster连接地址, 格式：<host>:<port>, 多以逗号(,)分隔
//;;cluster_address = 172.16.100.26:6379, 172.16.100.26:6380
//;;;; [必选] redis密码
//;;password = pi2paUAEDrTwfD9MzDnkTGDIm-QB0FLH
//;;;; [必选] redis连接超时设置
//;;timeout = 0
//;;;; [可选] redis连接的db, 最好设置，每个db一个连接, cluster不需要此项
//;;db = 0
//;;;; [可选], 值on/off，默认off，执行命令写入日志
//;;cmdlog = off

func main() {
	cfg := map[string]interface{}{
		//"address": "172.16.100.112:6379",
		"sentinel_address":    "172.16.100.112:26379,        172.16.100.112:26380, 172.16.100.112:26381",
		"sentinel_mastername": "mymaster",
		//"cluster_address":     "172.16.100.112:26379, 172.16.100.112:26380, 172.16.100.112:26381",
		"password": "Gmck7X02",
		"timeout":  "0",
		"db":       "1",
		"cmdlog":   "off",
	}
	fmt.Println(cfg)
	e := redis.SetCfg("default", cfg)
	fmt.Println(e)
	ct, e := redis.GetClient("default")
	fmt.Println(ct, e)
	ct.Set("a", "b")
	r, e := ct.Hset("jw", "tester", "tester")
	fmt.Println("hset jw tester tester", r, e)
	r, e = ct.Hset("jw", "tester0", "tester")
	fmt.Println("hset jw tester0 tester", r, e)
	s, e := ct.Hget("jw", "tester")
	fmt.Println("hget jw tester", s, e)
	r, e = ct.Hdel("jw", "tester")
	fmt.Println("hdel jw tester", r, e)

	//封装方式
	r, e = ct.Hset("jw", "tester", "tester")
	r, e = ct.Hset("jw", "tester0", "tester")
	r, e = ct.Hdel("jw", "tester", "tester0")
	fmt.Println("hdel jw tester tester0", r, e)

	//原生方式
	r, e = ct.Hset("jw", "tester", "tester")
	r, e = ct.Hset("jw", "tester0", "tester")
	r, e = redis.FInt(ct.Do("HDEL", redis.Strs2Ifs("jw", "tester0", "tester")...))
	fmt.Println("hdel jw tester tester0", r, e)
}
