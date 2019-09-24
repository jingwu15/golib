package main

import (
	"github.com/jingwu15/golib/sched"
)


func main() {
    sched.Tester("1 30 4 * * *",  "2019-09-24 12:12:12", 20)
    return
}
