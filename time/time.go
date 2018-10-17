package time

import (
//	"fmt"
	"time"
	"strconv"
	"strings"
    "regexp"
    "errors"
)

type Duration = time.Duration
type Time struct {
	Current time.Time
}

//获取当前时间
func Now() Time {
	return Time{
		Current: time.Now(),
	}
}

//获取当前时间戳
func Unix(sec int64, nesc int64) Time {
	return Time{
		Current: time.Unix(sec, nesc),
	}
}

//获取时间对应的时间戳
func (t Time) Unix() int64 {
	return t.Current.Unix()
}

//获取时间对应的周
func (t Time) Weekday() string {
    weeks := map[time.Weekday]string{time.Monday: "1", time.Tuesday: "2", time.Wednesday: "3", time.Thursday: "4", time.Friday: "5", time.Saturday: "6", time.Sunday: "0"}
    return weeks[t.Current.Weekday()]
}

//获取时间对应的周
func (t Time) WeekInt(start int) int {
    var weeks map[time.Weekday]int
    if start == 0 {
        weeks = map[time.Weekday]int{time.Monday: 1, time.Tuesday: 2, time.Wednesday: 3, time.Thursday: 4, time.Friday: 5, time.Saturday: 6, time.Sunday: 0}
    } else {
        weeks = map[time.Weekday]int{time.Monday: 1, time.Tuesday: 2, time.Wednesday: 3, time.Thursday: 4, time.Friday: 5, time.Saturday: 6, time.Sunday: 7}
    }
    return weeks[t.Current.Weekday()]
}

//获取时间对应的周
func (t Time) WeekStr(start int) string {
    var weeks map[time.Weekday]string
    if start == 0 {
       weeks = map[time.Weekday]string{time.Monday: "1", time.Tuesday: "2", time.Wednesday: "3", time.Thursday: "4", time.Friday: "5", time.Saturday: "6", time.Sunday: "0"}
    } else {
       weeks = map[time.Weekday]string{time.Monday: "1", time.Tuesday: "2", time.Wednesday: "3", time.Thursday: "4", time.Friday: "5", time.Saturday: "6", time.Sunday: "7"}
    }
    return weeks[t.Current.Weekday()]
}

//格式化时间为字符串
//支持各种linux通用的时间描述符，同date
func (t Time) Format(format string) string {
	//2006-01-02T15:04:05.999999999Z07:00
	format = strings.Replace(format, "a", "Mon", -1)
	format = strings.Replace(format, "A", "Monday", -1)

	format = strings.Replace(format, "b", "Jan", -1)
	format = strings.Replace(format, "B", "January", -1)

	format = strings.Replace(format, "c", "Mon Jan 2 2006 15:04:05 PM MST", -1) //RFC1123 如 Mon 09 Apr 2018 11:41:34 AM CST
	//format = strings.Replace(format, "C", "20", -1) 	//世纪 @不可用

	format = strings.Replace(format, "d", "02", -1)			//一月中的第几天
	format = strings.Replace(format, "D", "01/02/2006", -1)	//日期，如 %m/%d/%y

	format = strings.Replace(format, "e", " 2", -1)			//一月中的第几天，空格补齐,如 02->_2

	format = strings.Replace(format, "F", "2006-01-02 15:04:05", -1)	//日期，如 %y-%m-%d

	//format = strings.Replace(format, "g", "06", -1)		//年，两位数
	//format = strings.Replace(format, "G", "2006", -1)		//年，四位数

	format = strings.Replace(format, "H", "15", -1)		//小时数，24时制
	format = strings.Replace(format, "I", "03", -1)		//小时数，12时制

	format = strings.Replace(format, "j", strconv.Itoa(t.Current.YearDay()), -1)		//一年中的第几天

	format = strings.Replace(format, "m", "01", -1)			//月
	format = strings.Replace(format, "M", "04", -1)			//分钟

	format = strings.Replace(format, "N", strconv.Itoa(t.Current.Nanosecond()), -1)	//纳秒 000000000..999999999

	format = strings.Replace(format, "Y", "2006", -1)		//年，四位数
	format = strings.Replace(format, "y", "06", -1)			//年，两位数

	format = strings.Replace(format, "s", "05", -1)		//秒，时间戳
	format = strings.Replace(format, "S", "05", -1)		//秒，60制

	format = strings.Replace(format, "z", "0700", -1)		//时区，0700
	format = strings.Replace(format, "Z", "MST", -1)		//时区，MST

	//fmt.Println(format, t.Current.Nanosecond())
	return t.Current.Format(format)
}

//以秒为单位进行休眠
func Sleep(sec Duration) {
	time.Sleep(sec * time.Second)
}

//以毫秒为单位进行休眠
func SleepMsec(msec Duration) {
	time.Sleep(msec * time.Millisecond)
}

func StrToTime(timeStr string) (Time, error) {
    regYear := "[1-9][0-9]{3}"
    regMonth := "(0[0-9])|(1[0-2])"
    regDay := "([0-2][0-9])|(3[0-1])"
    regHour24 := "([0-1][0-9])|(2[0-3])"
    //regHour12 := "(0[0-9])|(1[0-1])"
    regMinute := "[0-5][0-9]"
    regSecond := "[0-5][0-9]"

    regs := []map[string]string{
        //YYYY-MM-DD
        map[string]string{
            //"reg": "^" + regYear + "-" + regMonth + "-" + regDay + "$",
            "reg": "^" + regYear + "-" + regMonth + "-" + regDay + "$",
            "format_str": "YYYY-MM-DD",
            "format": "2006-01-02 15:04:05 -0700",
            "pre": "",
            "suffix": " 00:00:00 +0800",
        },
        //YYYYMMDD
        map[string]string{
            //"reg": "^" + regYear + regMonth + regDay + "$",
            "reg": "^" + regYear + regMonth + regDay + "$",
            "format_str": "YYYYMMDD",
            "format": "20060102 15:04:05 -0700",
            "pre": "",
            "suffix": " 00:00:00 +0800",
        },
        //YYYY-MM-DD HH:II:SS
        map[string]string{
            //"reg": "^" + regYear + "-" + regMonth + "-" + regDay + " "+ regHour24 + ":" + regMinute + ":" + regSecond + "$",
            "reg": "^" + regYear + "-" + regMonth + "-" + regDay + " "+ regHour24 + ":" + regMinute + ":" + regSecond + "$",
            "format_str": "YYYY-MM-DD HH:II:SS",
            "format": "2006-01-02 15:04:05 -0700",
            "pre": "",
            "suffix": " +0800",
        },
        //YYYYMMDDHHIISS
        map[string]string{
            "reg": "^" + regYear + regMonth + regDay + regHour24 + regMinute + regSecond + "$",
            "format_str": "YYYYMMDDHHIISS",
            "format": "20060102150405  -0700",
            "pre": "",
            "suffix": " +0800",
        },
        ////HH:II:SS
        //map[string]string{
        //    "reg": "^" + regHour24 + ":" + regMinute + ":" + regSecond + "$",
        //    "format_str": "HH:II:SS",
        //    "format": "20060102 15:04:05",
        //    "pre": "00000000 ",
        //    "suffix": "",
        //},
        ////HHIISS
        //map[string]string{
        //    "reg": "^" + regHour24 + regMinute + regSecond + "$",
        //    "format_str": ":HIISS",
        //    "format": "20060102 150405",
        //    "pre": "00000000 ",
        //    "suffix": "",
        //},
    }
    var err error
    var newTime Time
    var tmpTime time.Time
    flag := 0
    for _,row := range regs {
        match, _ := regexp.MatchString(row["reg"], timeStr)
        formatLen := strings.Count(row["format_str"], "")
        timeLen := strings.Count(timeStr, "")
        if match && formatLen == timeLen {
            tmp := row["pre"] + timeStr + row["suffix"]
            tmpTime, err = time.Parse(row["format"], tmp)
            newTime.Current = tmpTime.Local()
            flag = 1
            break
        }
    }
    if flag == 0 {
        err = errors.New("do'not matched")
    }

    return newTime, err
}

