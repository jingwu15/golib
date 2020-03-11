package log

import (
    "os"
    "fmt"
    "io/ioutil"
    "github.com/sirupsen/logrus"
    "github.com/rifflock/lfshook"
)

type Cfg struct {
    Dir string
    Pre string
    Level string
}
type Logger struct {
    logger *logrus.Logger
    entry  *logrus.Entry
}

var (
    keyDefault = "default"
    loggers = map[string]*Logger{}
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

var cfg = Cfg{Dir: "/tmp", Pre: "log_default", Level: "debug"}
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

func Get(keys ...string) *Logger {
    key := "default"
    if len(keys) > 0 { key = keys[0] }
    if _, ok := loggers[key]; !ok {
        logger := logrus.New()
        logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
        logger.SetOutput(os.Stdout)
        logger.SetLevel(levels[cfg.Level])

        fname := ""
        if key == "default" {
            fname = fmt.Sprintf("%s/%s.log", cfg.Dir, cfg.Pre)
        } else {
            fname = fmt.Sprintf("%s/%s_%s.log", cfg.Dir, cfg.Pre, key)
        }
        logger.Hooks.Add(lfshook.NewHook(
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
        if key == "default" {
            loggers[key] = &Logger{logger: logger, entry: logger.WithFields(map[string]interface{}{})}
        } else {
            loggers[key] = &Logger{logger: logger, entry: logger.WithFields(map[string]interface{}{"_logkey": key})}
        }
    }
    return loggers[key]
}

func Close_stdout(key string) {
    if _, ok := loggers[key]; ok {
        loggers[key].logger.SetOutput(ioutil.Discard)
    }
}

func Output_file(key, fname string) {
    fp, e := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
    if _, ok := loggers[key]; ok && e == nil {
        loggers[key].logger.SetOutput(fp)
    }
}

func Format(key, format string) {
    _, okLog := loggers[key]
    _, okFormat := formats[format]
    if okLog && okFormat {
        loggers[key].logger.SetFormatter(formats[format])
    }
}

func (l *Logger)WithFields(params map[string]interface{}) {
    for k, v := range l.entry.Data { params[k] = v }
    l.entry = l.logger.WithFields(params)
}

func (l *Logger)Debug(args ...interface{})  { l.entry.Debug(args...)  }
func (l *Logger)Info(args ...interface{})   { l.entry.Info(args...)   }
func (l *Logger)Warn(args ...interface{})   { l.entry.Warn(args...)   }
func (l *Logger)Error(args ...interface{})  { l.entry.Error(args...)  }
func (l *Logger)Fatal(args ...interface{})  { l.entry.Fatal(args...)  }
func (l *Logger)Panic(args ...interface{})  { l.entry.Panic(args...)  }
func (l *Logger)Print(args ...interface{})  { l.entry.Print(args...)  }

func (l *Logger)Debugln(args ...interface{}) { l.entry.Debugln(args...) }
func (l *Logger)Infoln(args ...interface{})  { l.entry.Infoln(args...)  }
func (l *Logger)Warnln(args ...interface{})  { l.entry.Warnln(args...)  }
func (l *Logger)Errorln(args ...interface{}) { l.entry.Errorln(args...) }
func (l *Logger)Fatalln(args ...interface{}) { l.entry.Fatalln(args...) }
func (l *Logger)Panicln(args ...interface{}) { l.entry.Panicln(args...) }
func (l *Logger)Println(args ...interface{}) { l.entry.Println(args...) }

func (l *Logger)Debugf(format string, args ...interface{})  { l.entry.Debugf(format, args...) }
func (l *Logger)Infof(format string, args ...interface{})   { l.entry.Infof(format, args...)  }
func (l *Logger)Warnf(format string, args ...interface{})   { l.entry.Warnf(format, args...)  }
func (l *Logger)Errorf(format string, args ...interface{})  { l.entry.Errorf(format, args...) }
func (l *Logger)Fatalf(format string, args ...interface{})  { l.entry.Fatalf(format, args...) }
func (l *Logger)Panicf(format string, args ...interface{})  { l.entry.Panicf(format, args...) }
func (l *Logger)Printf(format string, args ...interface{})  { l.entry.Printf(format, args...) }

func WithFields(params map[string]interface{}) {    Get(keyDefault).WithFields(params)          }

func Debug(args ...interface{}) { Get(keyDefault).Debug(args...)  }
func Info(args ...interface{})  { Get(keyDefault).Info(args...)   }
func Warn(args ...interface{})  { Get(keyDefault).Warn(args...)   }
func Error(args ...interface{}) { Get(keyDefault).Error(args...)  }
func Fatal(args ...interface{}) { Get(keyDefault).Fatal(args...)  }
func Panic(args ...interface{}) { Get(keyDefault).Panic(args...)  }
func Print(args ...interface{}) { Get(keyDefault).Print(args...)  }

func Debugln(args ...interface{}) { Get(keyDefault).Debugln(args...) }
func Infoln(args ...interface{})  { Get(keyDefault).Infoln(args...)  }
func Warnln(args ...interface{})  { Get(keyDefault).Warnln(args...)  }
func Errorln(args ...interface{}) { Get(keyDefault).Errorln(args...) }
func Fatalln(args ...interface{}) { Get(keyDefault).Fatalln(args...) }
func Panicln(args ...interface{}) { Get(keyDefault).Panicln(args...) }
func Println(args ...interface{}) { Get(keyDefault).Println(args...) }

func Debugf(format string, args ...interface{}) { Get(keyDefault).Debugf(format, args...)  }
func Infof(format string, args ...interface{})  { Get(keyDefault).Infof(format, args...)   }
func Warnf(format string, args ...interface{})  { Get(keyDefault).Warnf(format, args...)   }
func Errorf(format string, args ...interface{}) { Get(keyDefault).Errorf(format, args...)  }
func Fatalf(format string, args ...interface{}) { Get(keyDefault).Fatalf(format, args...)  }
func Panicf(format string, args ...interface{}) { Get(keyDefault).Panicf(format, args...)  }
func Printf(format string, args ...interface{}) { Get(keyDefault).Printf(format, args...)  }
