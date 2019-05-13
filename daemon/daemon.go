package daemon

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"strconv"
	"strings"
	"syscall"
	"os/exec"
	"os/signal"
	"io/ioutil"
	"path/filepath"
	"github.com/spf13/viper"
	"github.com/erikdubbelboer/gspt"
	logchan "github.com/jingwu15/golib/logchan"
	log "github.com/sirupsen/logrus"
)

var runFlag int = 1

func findProcess(procTitle string) ([]int, error) {
	var err error
	matches, err := filepath.Glob("/proc/*/cmdline")
	if err != nil {
		return nil, err
	}
	var pid int
	var pids = []int{}
	var tmp []string
	var body []byte
	for _, filename := range matches {
		body, err = ioutil.ReadFile(filename)
		if err == nil {
			if bytes.HasPrefix(body, []byte(procTitle)) {
				tmp = strings.Split(filename, "/")
				pid, _ = strconv.Atoi(tmp[2])
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}

//初始化日志
func InitLog() {
	go logchan.LogWrite()
	log.SetFormatter(&log.JSONFormatter{})
	//log.SetOutput(ioutil.Discard)
    logLevelMap := map[string]log.Level{
        "trace": log.TraceLevel,
        "debug": log.DebugLevel,
        "info":  log.InfoLevel,
        "warn":  log.WarnLevel,
        "error": log.ErrorLevel,
        "fatal": log.FatalLevel,
        "panic": log.PanicLevel,
    }
    logLevel := viper.GetString("log.level")
    if level, ok := logLevelMap[logLevel]; ok {
	    log.SetLevel(level)
    } else {
	    log.SetLevel(log.DebugLevel)
    }
	config := map[string]string{
		"error":      viper.GetString("log.error"),
		"info":       viper.GetString("log.info"),
		"writeDelay": viper.GetString("log.delay"),
		"cutType":    "day",
	}
	logChanHook := logchan.NewLogChanHook(config)
	log.AddHook(&logChanHook)
}

//处理信号，停止及重启
func handleSignals(fun_Close func()) {
	var sig os.Signal
	var signalChan = make(chan os.Signal, 100)
	signal.Notify(
		signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)
	for {
		sig = <-signalChan
		switch sig {
		case syscall.SIGTERM:
			log.Info("stop")
			logchan.LogClose()
            fun_Close()
			runFlag = 0
		case syscall.SIGUSR2:
			log.Info("restart")
			logchan.LogClose()
            fun_Close()
		default:
		}
	}
}

func Run(procTitle string, fun_Run func(), fun_Close func()) {
	gspt.SetProcTitle(procTitle)

	go handleSignals(fun_Close)
	InitLog()

    fun_Run()

	log.Info(fmt.Sprintf("%s is running", procTitle))
	for {
		if runFlag == 1 {
			//未结束，一直等待
			time.Sleep(time.Duration(2) * time.Second)
			//log.Info(fmt.Sprintf("%s is running", procTitle))
		} else {
			log.Info(fmt.Sprintf("%s is shut down", procTitle))
			break
		}
	}
}

func Start(procTitle string) {
	var err error
    cmdArgs := os.Args
	cmdArgs[0], _ = filepath.Abs(cmdArgs[0])
    cmd := strings.Replace(strings.Join(cmdArgs, " "), "start", "run", 1)
    cmd = fmt.Sprintf("nohup %s &> /tmp/log_gyworker.log &", cmd)
	client := exec.Command("sh", "-c", cmd)
	err = client.Start()
	if err != nil {
		fmt.Println(procTitle, "start error:")
		fmt.Println(err)
		return
	}
	err = client.Wait()
	if err != nil {
		fmt.Println(procTitle, "start error:")
		fmt.Println(err)
		return
	}
	fmt.Println(procTitle, "is started")
	return
}

//func Restart(procTitle string) {
//	pids, err := findProcess(procTitle)
//	if err != nil {
//		log.Error(err)
//		return
//	}
//	for _, pid := range pids {
//		syscall.Kill(pid, syscall.SIGUSR2)
//	}
//	fmt.Println(procTitle, "is restarted")
//}

func Stop(procTitle string) {
	pids, err := findProcess(procTitle)
	if err != nil {
		log.Error(err)
		return
	}
	for _, pid := range pids {
		syscall.Kill(pid, syscall.SIGTERM)
	}
	fmt.Println(procTitle, "is stoped")
}
