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
    LOG_PRE = ""
    logkeyDefault = "default"
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

func SetLogDir(logdir string) error {
    fp, e := os.OpenFile(logdir, os.O_RDWR, 0644)
    if os.IsNotExist(e) {           //目录不存在，则新建
        e = os.MkdirAll(logdir, os.ModePerm)
        if e != nil {       //新建目录失败, 使用默认目录
            return fmt.Errorf("新建日志目录失败!")
        }
    }
    fp.Close()
    LOG_DIR = logdir
    return nil
}

func SetLogPre(logpre string) {
    LOG_PRE = logpre
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

        fname := fmt.Sprintf("%s/%s%s.log", LOG_DIR, LOG_PRE, key)
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

func SetLevel(key, level string) {
    logger, okLog := loggers[key]
    _, okLevel := levels[level]
    if okLog && okLevel { logger.SetLevel(levels[level]) }
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

//func Debug(args ...interface{}) {
//    loggers[logkeyDefault].Debug(args...)
//}
//func Debugf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Debugf(format, args...)
//}
//func Debugln(args ...interface{}) {
//    loggers[logkeyDefault].Debugln(args...)
//}
//func Error(args ...interface{}) {
//    loggers[logkeyDefault].Error(args...)
//}
//func Errorf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Errorf(format, args...)
//}
//func Errorln(args ...interface{}) {
//    loggers[logkeyDefault].Errorln(args...)
//}
//func Fatal(args ...interface{}) {
//    loggers[logkeyDefault].Fatal(args...)
//}
//func Fatalf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Fatalf(format, args...)
//}
//func Fatalln(args ...interface{}) {
//    loggers[logkeyDefault].Fatalln(args...)
//}
//func Info(args ...interface{}) {
//    loggers[logkeyDefault].Info(args...)
//}
//func Infof(format string, args ...interface{}) {
//    loggers[logkeyDefault].Infof(format, args...)
//}
//func Infoln(args ...interface{}) {
//    loggers[logkeyDefault].Infoln(args...)
//}
//func Panic(args ...interface{}) {
//    loggers[logkeyDefault].Panic(args...)
//}
//func Panicf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Panicf(format, args...)
//}
//func Panicln(args ...interface{}) {
//    loggers[logkeyDefault].Panicln(args...)
//}
//func Print(args ...interface{}) {
//    loggers[logkeyDefault].Print(args...)
//}
//func Printf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Printf(format, args...)
//}
//func Println(args ...interface{}) {
//    loggers[logkeyDefault].Println(args...)
//}
//func Warn(args ...interface{}) {
//    loggers[logkeyDefault].Warn(args...)
//}
//func Warnf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Warnf(format, args...)
//}
//func Warning(args ...interface{}) {
//    loggers[logkeyDefault].Warning(args...)
//}
//func Warningf(format string, args ...interface{}) {
//    loggers[logkeyDefault].Warningf(format, args...)
//}
//func Warningln(args ...interface{}) {
//    loggers[logkeyDefault].Warningln(args...)
//}
//func Warnln(args ...interface{}) {
//    loggers[logkeyDefault].Warnln(args...)
//}
