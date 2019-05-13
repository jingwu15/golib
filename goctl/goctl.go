package goctl

import(
    "fmt"
    "reflect"
    "runtime"
    "strings"
    "github.com/jingwu15/golib/time"
    log "github.com/sirupsen/logrus"
)

var logger = log.WithFields(log.Fields{"btype": "goctl"})

func GetFunctionName(i interface{}) string {
    // 获取函数名称
    fname := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
    return strings.Replace(fname, "/", ".", -1)
}

//只限制启动多少个go程, go程退出则不管
func GoLimit(fun func()(error), max int) {
    for i:=0; i < max; i++ {
        go fun()
    }
}

func GoToChan(fun func()(error), keepQ chan int) {
    fname := GetFunctionName(fun)
    keepQ <- 1         //启动
    logger.Info(fmt.Sprintf("func[GoToChan] fun[%s] run start", fname))
    fun()              //退出
    logger.Info(fmt.Sprintf("func[GoToChan] fun[%s] run end", fname))
    keepQ <- -1
}

//保持多少个go程在运行状态
func GoKeep(fun func()(error), max int) {
    keepQ := make(chan int, max)
    for i:=0; i < max; i++ {
        go GoToChan(fun, keepQ)
    }

    for{
        select{
        case v, _  := <- keepQ:
            if v == -1 {
                go GoToChan(fun, keepQ)
            }
		case <-time.After(10):
        }
    }
//GoRun:

}

