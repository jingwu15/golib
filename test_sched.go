package main

import (
	"github.com/jingwu15/golib/sched"
)

func main() {
    //sched.Tester("1 * * * * *",  "2019-09-25 23:45:12", 1)
    //sched.Tester("1 * * * * *",  "2019-09-25 13:45:12", 1)
    //sched.Tester("1 1 10 * * *",  "2019-09-25 03:45:12", 10)
    //sched.Tester("1 1 10 * * */2",  "2019-09-25 13:45:12", 10)
    //sched.Tester("1 1 10 * * 1",  "2019-09-25 13:45:12", 10)
    sched.Tester("1 */30 * * * *",  "2020-02-29 01:05:12", 1)        //验证润年
    return
}
