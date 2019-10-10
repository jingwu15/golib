package log

import (
    "os"
    "fmt"
    "io/ioutil"
    "github.com/sirupsen/logrus"
    "github.com/rifflock/lfshook"
)

var (
    LOG_DIR = "/tmp/"
    LOG_PRE = "log_"
    keyDefault = "default"
    loggers = map[string]*logrus.Logger{}
    levels = map[string]logrus.Level{
        "trace": logrus.TraceLevel,
        "debug": logrus.DebugLevel,
        "info":  logrus.InfoLevel,
        "warn":  logrus.WarnLevel,
        "error": logrus.ErrorLevel,
        "fatal": logrus.FatalLevel,
        "panic": logrus.PanicLevel,
    }
)

type Cfg struct {
    Dir string
    Pre string
    Level string
}

var cfg = Cfg{Dir: "/tmp", Pre: "log_", Level: "debug"}
func InitCfg(cfgUser Cfg) error {
    if cfgUser.Dir == ""   { cfgUser.Dir   = cfg.Dir   }
    if cfgUser.Pre == ""   { cfgUser.Pre   = cfg.Pre   }
    if cfgUser.Level == "" { cfgUser.Level = cfg.Level }
    fp, e := os.OpenFile(cfgUser.Dir, os.O_RDWR, 0644)
    if os.IsNotExist(e) {           //目录不存在，则新建
        e = os.MkdirAll(cfgUser.Dir, os.ModePerm)
        if e != nil {       //新建目录失败, 使用默认目录
            return fmt.Errorf("新建日志目录失败!")
        }
    }
    fp.Close()
    if _, ok := levels[cfgUser.Level]; !ok { return fmt.Errorf("日志等级设置错误!") }

    cfg.Dir   = cfgUser.Dir
    cfg.Level = cfgUser.Level
    cfg.Pre   = cfgUser.Pre
    return nil
}

var formats = map[string]logrus.Formatter{
    "json": &logrus.JSONFormatter{},
    "text": &logrus.TextFormatter{FullTimestamp: true, DisableColors: true, TimestampFormat: "2006-01-02 15:04:05"},
}

func Get(keys ...string) *logrus.Logger {
    key := "default"
    if len(keys) > 0 { key = keys[0] }
    if _, ok := loggers[key]; !ok {
        loggers[key] = logrus.New()
        loggers[key].SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
        loggers[key].SetOutput(os.Stdout)
        loggers[key].SetLevel(levels[cfg.Level])

        fname := ""
        if key == "default" {
            fname = fmt.Sprintf("%s/%s.log", cfg.Dir, cfg.Pre)
        } else {
            fname = fmt.Sprintf("%s/%s_%s.log", cfg.Dir, cfg.Pre, key)
        }
        loggers[key].Hooks.Add(lfshook.NewHook(
            lfshook.PathMap{
                logrus.TraceLevel: fname,
                logrus.DebugLevel: fname,
                logrus.InfoLevel: fname,
                logrus.WarnLevel: fname,
                logrus.ErrorLevel: fname,
                logrus.FatalLevel: fname,
                logrus.PanicLevel: fname,
	        },
		    &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"},
	    ))
        loggers[key].WithFields(logrus.Fields{"btype": key})
    }
    return loggers[key]
}

func Close_stdout(key string) {
    if _, ok := loggers[key]; ok {
        loggers[key].SetOutput(ioutil.Discard)
    }
}

func Output_file(key, fname string) {
    fp, e := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if _, ok := loggers[key]; ok && e == nil {
        loggers[key].SetOutput(fp)
    }
}

func Format(key, format string) {
    logger, okLog := loggers[key]
    _, okFormat := formats[format]
    if okLog && okFormat {
        logger.SetFormatter(formats[format])
    }
}

func Debug(args ...interface{}) {
    Get(keyDefault).Debug(args...)
}
func Debugf(format string, args ...interface{}) {
    Get(keyDefault).Debugf(format, args...)
}
func Debugln(args ...interface{}) {
    Get(keyDefault).Debugln(args...)
}
func Error(args ...interface{}) {
    Get(keyDefault).Error(args...)
}
func Errorf(format string, args ...interface{}) {
    Get(keyDefault).Errorf(format, args...)
}
func Errorln(args ...interface{}) {
    Get(keyDefault).Errorln(args...)
}
func Fatal(args ...interface{}) {
    Get(keyDefault).Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
    Get(keyDefault).Fatalf(format, args...)
}
func Fatalln(args ...interface{}) {
    Get(keyDefault).Fatalln(args...)
}
func Info(args ...interface{}) {
    Get(keyDefault).Info(args...)
}
func Infof(format string, args ...interface{}) {
    Get(keyDefault).Infof(format, args...)
}
func Infoln(args ...interface{}) {
    Get(keyDefault).Infoln(args...)
}
func Panic(args ...interface{}) {
    Get(keyDefault).Panic(args...)
}
func Panicf(format string, args ...interface{}) {
    Get(keyDefault).Panicf(format, args...)
}
func Panicln(args ...interface{}) {
    Get(keyDefault).Panicln(args...)
}
func Print(args ...interface{}) {
    Get(keyDefault).Print(args...)
}
func Printf(format string, args ...interface{}) {
    Get(keyDefault).Printf(format, args...)
}
func Println(args ...interface{}) {
    Get(keyDefault).Println(args...)
}
func Warn(args ...interface{}) {
    Get(keyDefault).Warn(args...)
}
func Warnf(format string, args ...interface{}) {
    Get(keyDefault).Warnf(format, args...)
}
func Warning(args ...interface{}) {
    Get(keyDefault).Warning(args...)
}
func Warningf(format string, args ...interface{}) {
    Get(keyDefault).Warningf(format, args...)
}
func Warningln(args ...interface{}) {
    Get(keyDefault).Warningln(args...)
}
func Warnln(args ...interface{}) {
    Get(keyDefault).Warnln(args...)
}
