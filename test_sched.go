package main

import (
    "fmt"
	"github.com/jingwu15/golib/time"
	"github.com/jingwu15/golib/sched"
)

func main() {
	////测试计划任务的时间计算
	var ttime, now time.Time
	var crontabSched = "1 */30 * * * *";
    var i int64
	now = time.Now()
    for i = 0; i < 4000; i++ {
        ttime = time.Unix(now.Unix() + i, 0)
	    fmt.Println(ttime.ToStr(), crontabSched, sched.Create(ttime, crontabSched).NextTime());
    }
}
