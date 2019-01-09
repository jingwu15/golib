package logchan

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

var writeDelay string = "0"

var cutTypes = map[string]string{"day": "1", "month": "1", "hour": "1"}
var logCutType = "" //     day/month/hour
var logFilePs = make(map[string]*os.File)
var logFiles = make(map[string][]string)
var levelFileMap = make(map[string][]string)

type LogChanHook struct {
}

var logChan = make(chan string, 1000000)

func LogFormatLevel(level interface{}) string {
	var format string
	switch level {
	case logrus.DebugLevel:
		format = "debug"
	case logrus.InfoLevel:
		format = "info"
	case logrus.WarnLevel:
		format = "warn"
	case logrus.ErrorLevel:
		format = "error"
	case logrus.FatalLevel:
		format = "fatal"
	case logrus.PanicLevel:
		format = "panic"
	default:
		format = "debug"
	}
	return format
}

func NewLogChanHook(config map[string]string) LogChanHook {
	var ok bool
	if _, ok = config["debug"]; ok {
		logFiles[config["debug"]] = []string{"debug", "info", "warn", "error", "fatal", "panic"}
	}
	if _, ok = config["info"]; ok {
		logFiles[config["info"]] = []string{"info", "warn", "error", "fatal", "panic"}
	}
	if _, ok = config["warn"]; ok {
		logFiles[config["warn"]] = []string{"warn", "error", "fatal", "panic"}
	}
	if _, ok = config["error"]; ok {
		logFiles[config["error"]] = []string{"error", "fatal", "panic"}
	}
	if _, ok = config["fatal"]; ok {
		logFiles[config["fatal"]] = []string{"fatal", "panic"}
	}
	if _, ok = config["panic"]; ok {
		logFiles[config["panic"]] = []string{"panic"}
	}
	if _, ok = config["writeDelay"]; ok {
		writeDelay = config["writeDelay"]
	}
	if _, ok = config["cutType"]; ok {
		if _, ok = cutTypes[config["cutType"]]; ok {
			logCutType = config["cutType"]
		}
	}
	for logfile, levels := range logFiles {
		for _, level := range levels {
			if _, ok = levelFileMap[level]; !ok {
				levelFileMap[level] = []string{}
			}
			levelFileMap[level] = append(levelFileMap[level], logfile)
		}
	}

	logChanHook := LogChanHook{}
	return logChanHook
}

func (hook *LogChanHook) Fire(entry *logrus.Entry) error {
	tmp, err := entry.String()
	timeRawStr := entry.Time.String()
	timeRawStr = strings.Replace(timeRawStr, " ", "", -1)
	timeRawStr = strings.Replace(timeRawStr, "-", "", -1)
	timeRawStr = strings.Replace(timeRawStr, ":", "", -1)
	timeStr := timeRawStr[0:14]
	var line bytes.Buffer
	line.WriteString(LogFormatLevel(entry.Level))
	line.WriteString(",")
	line.WriteString(timeStr)
	line.WriteString(",")
	line.WriteString(tmp)
	if err == nil {
		logChan <- line.String()
		return nil
	} else {
		return err
	}
}

func GetLogFile(logfile string, timeStr string) string {
	var logfilekey string = ""
	if logCutType == "day" {
		logfilekey = logfile + "." + timeStr[0:8]
	} else if logCutType == "month" {
		logfilekey = logfile + "." + timeStr[0:6]
	} else if logCutType == "hour" {
		logfilekey = logfile + "." + timeStr[0:10]
	} else {
		logfilekey = logfile
	}
	if _, ok := logFilePs[logfilekey]; !ok {
		fh, err := os.OpenFile(logfilekey, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		logFilePs[logfilekey] = fh
	}
	return logfilekey
}

func LogWrite() {
	delay, _ := strconv.Atoi(writeDelay)
	var limit, i int
	var line string
	var ok bool
    limit := 0
	for {
        if limit > 100 {
            LogClose()
            limit = 0
        }
		limit = len(logChan)
		var bodys = make(map[string]*bytes.Buffer)
		for i = 0; i < limit; i++ {
			line = <-logChan
			tmp := strings.SplitN(line, ",", 3)
			if _, ok = levelFileMap[tmp[0]]; ok {
				for _, tmpfile := range levelFileMap[tmp[0]] {
					logfile := GetLogFile(tmpfile, tmp[1])
					if _, ok = bodys[logfile]; !ok {
						bodys[logfile] = bytes.NewBufferString("")
					}
					bodys[logfile].WriteString(tmp[2])
				}
			}
		}
		for logfile, body := range bodys {
			logFilePs[logfile].WriteString(body.String())
		}
		time.Sleep(time.Second * time.Duration(delay))
        limit++
	}
}

func LogClose() {
	var ok bool
	var bodys = make(map[string]*bytes.Buffer)
	limit := len(logChan)
	for i := 0; i < limit; i++ {
		line := <-logChan
		tmp := strings.SplitN(line, ",", 3)
		if _, ok = levelFileMap[tmp[0]]; ok {
			for _, tmpfile := range levelFileMap[tmp[0]] {
				logfile := GetLogFile(tmpfile, tmp[1])
				if _, ok = bodys[logfile]; !ok {
					bodys[logfile] = bytes.NewBufferString("")
				}
				bodys[logfile].WriteString(tmp[2])
			}
		}
	}
	for logfile, body := range bodys {
		logFilePs[logfile].WriteString(body.String())
	}
	for fkey, fp := range logFilePs {
		fp.Close()
        delete(logFilePs, fkey)
	}
}

func (hook *LogChanHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
