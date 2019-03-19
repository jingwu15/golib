package main

import (
	"fmt"
	libBS "github.com/jingwu15/golib/beanstalk"
)

func main() {
	tube := "tester"
	addr := "127.0.0.1:11300"
	bs, err := libBS.New(addr)
	if err != nil {
		fmt.Println("new", addr, err)
	}
	err = bs.Use(tube)
	fmt.Println("cmd-use", err)

	tubeUsed, err := bs.ListTubeUsed()
	fmt.Println("cmd-list-tube-used", tubeUsed, err)

	//for i := 0; i < 1000; i++ {
	//	//jobid, err := bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
	//	//fmt.Println("cmd-put", jobid, err)
	//	_, _ = bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
	//}
	//bs.ReConn()
	//fmt.Println("reConn", bs.Conn)
	//bs.Use(tube)
	//fmt.Println("cmd-use", tube)
	//for i := 0; i < 1000; i++ {
	//	jobid, err := bs.Put([]byte(fmt.Sprintf("tester----%04d", i)), 1, 0, 30)
	//	fmt.Println("cmd-put", jobid, err)
	//}

	count, err := bs.Watch(tube)
	fmt.Println("cmd-watch", tube, count, err)

	count, err = bs.Ignore("default")
	fmt.Println("cmd-ignore", "default", count, err)

	id, body, err := bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	err = bs.Delete(id)
	fmt.Println("cmd-delete", err)

	id, body, err = bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	err = bs.Release(id, 1, 0)
	fmt.Println("cmd-release", err)

	id, body, err = bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	err = bs.Bury(id, 1)
	fmt.Println("cmd-bury", id, err)

	id, body, err = bs.PeekBuried()
	fmt.Println("cmd-peek-buried", id, string(body), err)

	id, body, err = bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	err = bs.Touch(id)
	fmt.Println("cmd-touch", id, err)

	id, body, err = bs.Peek(id)
	fmt.Println("cmd-peek", id, string(body), err)

	id, body, err = bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	id, body, err = bs.PeekReady()
	fmt.Println("cmd-peek-ready", id, string(body), err)

	id, body, err = bs.PeekDelayed()
	fmt.Println("cmd-peek-delayed", id, string(body), err)

	count, err = bs.Kick(1000)
	fmt.Println("cmd-kick", count, err)

	id, body, err = bs.Reserve(1)
	fmt.Println("cmd-reserve", id, string(body), err)

	stats, err := bs.Stats()
	fmt.Println("cmd-stats", stats, err)

	stats, err = bs.StatsJob(id)
	fmt.Println("cmd-stats-job", id, stats, err)

	stats, err = bs.StatsTube(tube)
	fmt.Println("cmd-stats-tube", tube, stats, err)

	tubes, err := bs.ListTubes()
	fmt.Println("cmd-list-tubes", tubes, err)

	tubeUsed, err = bs.ListTubeUsed()
	fmt.Println("cmd-list-tube-used", tubeUsed, err)

	tubes, err = bs.ListTubeWatched()
	fmt.Println("cmd-list-tube-watched", tubes, err)

	err = bs.PauseTube(tube, 1)
	fmt.Println("cmd-pause-tube", err)

	err = bs.Quit()
	fmt.Println("cmd-quit", err)
}
