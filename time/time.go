package time

import (
	"fmt"
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


func Wrap(t time.Time) Time {
	return Time{
		Current: t,
	}
}

func (t Time) UnWrap() time.Time {
	return t.Current;
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

func (t Time) Date() (year, month, day int) {
    monthMap := map[time.Month]int{
        time.January: 1, time.February: 2, time.March: 3,     time.April: 4,    time.May: 5,       time.June: 6,
        time.July: 7,    time.August: 8,   time.September: 9, time.October: 10, time.November: 11, time.December: 12,
    }
    return t.Current.Year(), monthMap[t.Current.Month()], t.Current.Day()
}

func (t Time) Clock() (hour , min, sec int) {
    return t.Current.Clock()
}

func (t Time) Year() int {
    return t.Current.Year()
}

func (t Time) Month() int {
    monthMap := map[time.Month]int{
        time.January: 1, time.February: 2, time.March: 3,     time.April: 4,    time.May: 5,       time.June: 6,
        time.July: 7,    time.August: 8,   time.September: 9, time.October: 10, time.November: 11, time.December: 12,
    }
    return monthMap[t.Current.Month()]
}

func (t Time) YearDay() int {
    return t.Current.YearDay()
}

func (t Time) Day() int {
    return t.Current.Day()
}

func (t Time) Hour() int {
    return t.Current.Hour()
}

func (t Time) Minute() int {
    return t.Current.Minute()
}

func (t Time) Second() int {
    return t.Current.Second()
}

func (t Time) Nanosecond() int {
    return t.Current.Nanosecond()
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

func (t Time) ToStr() string {
    return t.Format(`Y-m-d H:M:s`)
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


//var t time.Time
//t, err = time.StrToTime("2018-01-01 12:34:56"); fmt.Println(t, err)
//t, err = time.StrToTime("2018-11-21"); fmt.Println(t, err)
//t, err = time.StrToTime("20181121"); fmt.Println(t, err)
//t, err = time.StrToTime("2018-1-2"); fmt.Println(t, err)
//t, err = time.StrToTime("12:34:56"); fmt.Println(t, err)
//t, err = time.StrToTime("123456"); fmt.Println(t, err)
//t, err = time.StrToTime("Jan 2 15:04:05 2016"); fmt.Println(t, err)           //ANSIC
//t, err = time.StrToTime("Jan 2 15:04:05 AWST 2016"); fmt.Println(t, err)      //UnixDate
//t, err = time.StrToTime("Jan 02 15:04:05 +0800 2016"); fmt.Println(t, err)    //RubyDate
//t, err = time.StrToTime("02 Jan 26 15:04 MST"); fmt.Println(t, err)           //RFC822  "02 Jan 06 15:04 MST"
//t, err = time.StrToTime("02 Jan 2126 15:04:05 MST"); fmt.Println(t, err)      //RFC1123 "02 Jan 2006 15:04:05 MST"
//t, err = time.StrToTime("02 Jan 06 15:04 +0800"); fmt.Println(t, err)         //RFC822Z
//t, err = time.StrToTime("02 Jan 2106 15:04:05 +0800"); fmt.Println(t, err)    //RFC1123Z
func StrToTime(timeStr string) (Time, error) {
    var err     error
    var now = time.Now()

    monthMap := map[time.Month]string{
        time.January: "1", time.February: "2", time.March: "3",     time.April: "4",    time.May: "5",       time.June: "6",
        time.July: "7",    time.August: "8",   time.September: "9", time.October: "10", time.November: "11", time.December: "12",
    }
    monthSMap := map[string]string{
        "Jan": "1", "Feb": "2", "Mar": "3", "Apr": "4", "May": "5", "Jun": "6",
        "Jul": "7", "Aug": "8", "Sep": "9", "Oct": "10", "Nov": "11", "Dec": "12",
        "January": "1", "February": "2", "March": "3",     "April": "4",                      "June": "6",
        "July": "7",    "August": "8",   "September": "9", "October": "10", "November": "11", "December": "12",
    }

    fields := map[string]string{
        "year":     strconv.Itoa(now.Year()),
        "month":    monthMap[now.Month()],
        "monthE":   "",
        "day":      strconv.Itoa(now.Day()),
        "hour":     "00",   //strconv.Itoa(now.Hour()),
        "minute":   "00",   //strconv.Itoa(now.Minute()),
        "second":   "00",   //strconv.Itoa(now.Second()),
        "zoneUtc":  "+0800",
        "zoneGmt":  "",
        //"zoneTime": "",
        "apm":      "",
    }

    regYear    := "(?P<year>(?:[0-9]{2})|(?:[1-9][0-9]{3}))"
    regYearA   := "(?P<year>[1-9][0-9]{3})"
    regMonth   := "(?P<month>(?:[1-9])|(?:0[0-9])|(?:1[0-2]))"
    regMonthA  := "(?P<month>(?:0[0-9])|(?:1[0-2]))"
    regMonthE  := "(?P<monthE>Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec|January|February|March|April|June|July|August|September|October|November|December)"
    regDay     := "(?P<day>(?:[1-9])|(?:[0-2][0-9])|(?:3[0-1]))"
    regDayA    := "(?P<day>(?:[0-2][0-9])|(?:3[0-1]))"
    regHour    := "(?P<hour>(?:[1-9])|(?:[0-1][0-9])|(?:2[0-3]))"
    regHourA   := "(?P<hour>(?:[0-1][0-9])|(?:2[0-3]))"
    regMinute  := "(?P<minute>[0-5]?[0-9])"
    regMinuteA := "(?P<minute>[0-5][0-9])"
    regSecond  := "(?P<second>[0-5]?[0-9])"
    regSecondA := "(?P<second>[0-5][0-9])"
    regZoneUtc := "(?P<zoneUtc>[+-](?:[0-1][0-9]{3}))"
    regZoneGmt := "(?P<zoneGmt>(?:[A-Z]{2,4}))"
    //regZoneTime:= "(?P<zoneTime>[0-9]{2}:[0-9]{2})"

    regs := []string{
        "^#year#-#month#-#day#$",                                      //YYYY-MM-DD
        "^#yearA##monthA##dayA#$",                                     //YYYYMMDD
        "^#year#-#month#-#day# #hour#:#minute#:#second#$",             //YYYY-MM-DD HH:II:SS
        "^#yearA##monthA##dayA##hourA##minuteA##secondA#$",            //YYYYMMDDHHIISS
        "^#hour#:#minute#:#second#$",                                  //HH:II:SS
        "^#hourA##minuteA##secondA#$",                                 //HHIISS
        "^#monthE# #day# #hour#:#minute#:#second# #year#$",            //ANSIC "Jan _2 15:04:05 2006"
        "^#monthE# #day# #hour#:#minute#:#second# #zoneGmt# #year#$",  //UnixDate "Jan _2 15:04:05 MST 2006"
        "^#monthE# #day# #hour#:#minute#:#second# #zoneUtc# #year#$",  //RubyDate "Jan 02 15:04:05 -0700 2006"
        "^#day# #monthE# #year# #hour#:#minute# #zoneGmt#$",           //RFC822  "02 Jan 06 15:04 MST"
        "^#day# #monthE# #year# #hour#:#minute# #zoneUtc#$",           //RFC822Z "02 Jan 06 15:04 -0700" // 使用数字表示时区的RFC822
        "^#day#-#monthE#-#year# #hour#:#minute#:#second# #zoneGmt#$",  //RFC850  "02-Jan-06 15:04:05 MST"
        "^#day# #monthE# #year# #hour#:#minute#:#second# #zoneGmt#$",  //RFC1123 "02 Jan 2006 15:04:05 MST"
        "^#day# #monthE# #year# #hour#:#minute#:#second# #zoneUtc#$",  //RFC1123Z "02 Jan 2006 15:04:05 -0700" // 使用数字表示时区的RFC1123
        //"^#yearA#-#monthA#-#dayA#T#hourA#:#minuteA#:#secondA#Z#zoneTime#$", //RFC3339 "2006-01-02T15:04:05Z07:00"
        //RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
    }

    flag := 0
    for _,regRule := range regs {
        //regRule := row["reg"]
        regRule = strings.Replace(regRule, "#year#",    regYear,    -1)
        regRule = strings.Replace(regRule, "#yearA#",   regYearA,   -1)
        regRule = strings.Replace(regRule, "#month#",   regMonth,   -1)
        regRule = strings.Replace(regRule, "#monthA#",  regMonthA,  -1)
        regRule = strings.Replace(regRule, "#day#",     regDay,     -1)
        regRule = strings.Replace(regRule, "#dayA#",    regDayA,    -1)
        regRule = strings.Replace(regRule, "#hour#",    regHour,    -1)
        regRule = strings.Replace(regRule, "#hourA#",   regHourA,   -1)
        regRule = strings.Replace(regRule, "#minute#",  regMinute,  -1)
        regRule = strings.Replace(regRule, "#minuteA#", regMinuteA, -1)
        regRule = strings.Replace(regRule, "#second#",  regSecond,  -1)
        regRule = strings.Replace(regRule, "#secondA#", regSecondA, -1)
        regRule = strings.Replace(regRule, "#monthE#",  regMonthE,  -1)
        regRule = strings.Replace(regRule, "#zoneUtc#", regZoneUtc, -1)
        regRule = strings.Replace(regRule, "#zoneGmt#", regZoneGmt, -1)
        //regRule = strings.Replace(regRule, "#zoneTime#",regZoneTime,-1)

        reg := regexp.MustCompile(regRule)
        if reg.MatchString(timeStr) {
            matchs := reg.FindStringSubmatch(timeStr)
            groups := reg.SubexpNames()
            for i, name := range groups {
                if i != 0 && name != "" {
                    fields[name] = matchs[i]
                }
            }
            flag = 1
            break
        }
    }
    if flag == 0 {
        return Time{}, errors.New("do'not matched")
    }

    //修正部分数据
    if fields["monthE"] != "" {
        fields["month"] = monthSMap[fields["monthE"]]
    }
    if (strings.Count(fields["year"], "")-1) == 2 {
        fields["year"]   = "20" + fields["year"]
    }
    if (strings.Count(fields["month"], "")-1) == 1 {
        fields["month"]  = "0" + fields["month"]
    }
    if (strings.Count(fields["day"], "")-1) == 1 {
        fields["day"]    = "0" + fields["day"]
    }
    if (strings.Count(fields["hour"], "")-1) == 1 {
        fields["hour"]   = "0" + fields["hour"]
    }
    if (strings.Count(fields["minute"], "")-1) == 1 {
        fields["minute"] = "0" + fields["minute"]
    }
    if (strings.Count(fields["second"], "")-1) == 1 {
        fields["second"] = "0" + fields["second"]
    }

    var data string
    var format string
    if fields["zoneGmt"] != "" {
        data   = fmt.Sprintf("%s-%s-%s %s:%s:%s %s", fields["year"], fields["month"], fields["day"], fields["hour"], fields["minute"], fields["second"], fields["zoneGmt"])
        format = fmt.Sprintf("%s-%s-%s %s:%s:%s %s", "2006", "01", "02", "15", "04", "05", "MST")
    //} else if fields["zoneTime"] != "" {
    //    data   = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ%s", fields["year"], fields["month"], fields["day"], fields["hour"], fields["minute"], fields["second"], fields["zoneTime"])
    //    format = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ%s", "2006", "01", "02", "15", "04", "05", "07:00")
    //    //RFC3339     = "2006-01-02T15:04:05Z07:00"
    } else {
        data   = fmt.Sprintf("%s-%s-%s %s:%s:%s %s", fields["year"], fields["month"], fields["day"], fields["hour"], fields["minute"], fields["second"], fields["zoneUtc"])
        format = fmt.Sprintf("%s-%s-%s %s:%s:%s %s", "2006", "01", "02", "15", "04", "05", "-0700")
    }
    //fmt.Println(data)
    //fmt.Println(format)
    //fmt.Println(fields)

    tmp, err := time.Parse(format, data)
    if err == nil {
        return Time{Current: tmp}, nil
    } else {
        return Time{}, err
    }
}

