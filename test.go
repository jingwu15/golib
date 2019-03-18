package main

import (
	"fmt"
	libBS "github.com/jingwu15/golib/beanstalk"
)

func main() {
	tube := "tester"
	addr := "127.0.0.1:11300"
	bs, err := libBS.New(addr)
	fmt.Println("new", bs.Conn)
	if err != nil {
		fmt.Println(addr, "连接失败")
	}
	bs.Use(tube)
	for i := 0; i < 1000; i++ {
		//jobid, err := bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
		//fmt.Println(jobid, err)
		_, _ = bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
	}
	bs.ReConn()
	fmt.Println("reconn", bs.Conn)
	bs.Use(tube)
	for i := 0; i < 1000; i++ {
		jobid, err := bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
		fmt.Println(jobid, err)
	}
}
