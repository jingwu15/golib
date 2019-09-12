package daemon

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"syscall"
	"os/exec"
	"os/signal"
	"io/ioutil"
	"path/filepath"
	log "github.com/sirupsen/logrus"
	logchan "github.com/jingwu15/golib/logchan"
)

var runFlag int = 1

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

func Run(pidfile string, fun_Run func(), fun_Close func()) error {
    fpid, err := os.OpenFile(pidfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
    if !os.IsNotExist(err) {    //进程文件存在, 检查进程是否真的存在
        body, err := ioutil.ReadAll(fpid)
        if err != nil { return fmt.Errorf("读取pidfile文件失败") }
        pid, _ := strconv.Atoi(string(body))
        if pid > 0 {            //进程ID存在
            pidfileOld := fmt.Sprintf("/proc/%s", string(body))
            _, err := os.Open(pidfileOld)
            if !os.IsNotExist(err) { return fmt.Errorf("程序运行中") }
        }
    }

    //写入新的进程ID
    pid := os.Getpid()
    pidStr := strconv.Itoa(pid)
    _, err = fpid.WriteString(pidStr)
    if err != nil { return fmt.Errorf("写入pidfile失败： %s", err.Error()) }

	go handleSignals(fun_Close)

    fun_Run()

	for {
		if runFlag == 1 {
			//未结束，一直等待
			time.Sleep(time.Duration(2) * time.Second)
			//log.Info(fmt.Sprintf("%s is running", procTitle))
		} else {
			//log.Info(fmt.Sprintf("%s is shut down", procTitle))
			break
		}
	}
    return nil
}

func Start(pidfile string) (err error) {
    fpid, err := os.OpenFile(pidfile, os.O_CREATE|os.O_RDWR, 0777)
    //进程文件已存在
    if !os.IsNotExist(err) {
        body, err := ioutil.ReadAll(fpid)
        fpid.Close()
        if err != nil { return fmt.Errorf("读取 %s 文件失败", pidfile) }
        pidStr := string(body)
        if pidStr == "" {
            os.Remove(pidfile)
        } else {
            pidfileOld := fmt.Sprintf("/proc/%s", pidStr)
            fpidProc, err := os.Open(pidfileOld)
            fpidProc.Close()
            if os.IsNotExist(err) {             //只存在一个空的 pidfile, 删除
                os.Remove(pidfile)
            } else {
                return fmt.Errorf("程序已运行")
            }
        }
    }

    cmdArgs := os.Args
	cmdArgs[0], _ = filepath.Abs(cmdArgs[0])
    cmd := strings.Replace(strings.Join(cmdArgs, " "), "start", "run", 1)
    //cmd = fmt.Sprintf("nohup %s &> /tmp/log_gyworker.log &", cmd)
    cmd = fmt.Sprintf("nohup %s &", cmd)
	client := exec.Command("sh", "-c", cmd)
	err = client.Start()
	if err != nil { return err }
	err = client.Wait()
	if err != nil { return err }
	return nil
}

//重启，先停后启
func Restart(pidfile string) (err error) {
    err = Stop(pidfile)
    if err != nil { return err }

    err = Start(pidfile)
    if err != nil { return err }

    return nil
}

//重载，要用户程序支持
func Reload(pidfile string) (err error) {
    fpid, err := os.OpenFile(pidfile, os.O_CREATE|os.O_RDWR, 0777)
    //进程文件不存在, 退出
    if os.IsNotExist(err) { fpid.Close(); return fmt.Errorf("进程已终止或 %s 文件不存在", pidfile) }

    //进程文件存在, 检查进程是否真的存在
    body, err := ioutil.ReadAll(fpid)
    if err != nil { fpid.Close(); return fmt.Errorf("读取 %s 文件失败", pidfile) }

    pid, _ := strconv.Atoi(string(body))
    //进程ID存在, 但错误
    if pid <= 0 { fpid.Close(); return fmt.Errorf("%s 文件内容错误，非数字", pidfile) }
    pidfileOld := fmt.Sprintf("/proc/%s", string(body))
    fpidProc, err := os.Open(pidfileOld)
    if os.IsNotExist(err) {
        os.Remove(pidfile)
        fpidProc.Close()
        return fmt.Errorf("进程已终止")
    }

	syscall.Kill(pid, syscall.SIGUSR2)
    return nil
}

func Stop(pidfile string) (err error) {
    fpid, err := os.OpenFile(pidfile, os.O_CREATE|os.O_RDWR, 0777)
    //进程文件不存在, 退出
    if os.IsNotExist(err) { fpid.Close(); return fmt.Errorf("进程已终止或 %s 文件不存在", pidfile) }

    //进程文件存在, 检查进程是否真的存在
    body, err := ioutil.ReadAll(fpid)
    if err != nil { fpid.Close(); return fmt.Errorf("读取 %s 文件失败", pidfile) }

    pid, _ := strconv.Atoi(string(body))
    //进程ID存在, 但错误
    if pid <= 0 { fpid.Close(); return fmt.Errorf("%s 文件内容错误，非数字", pidfile) }
    pidfileOld := fmt.Sprintf("/proc/%s", string(body))
    fpidProc, err := os.Open(pidfileOld)
    if os.IsNotExist(err) {
        os.Remove(pidfile)
        fpidProc.Close()
        return fmt.Errorf("进程已终止")
    }

    //强杀
	syscall.Kill(pid, syscall.SIGKILL)
    os.Remove(pidfile)

    //交由用户程序处理完任务，再退出
	//syscall.Kill(pid, syscall.SIGTERM)
    return nil
}
