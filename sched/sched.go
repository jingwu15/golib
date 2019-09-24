package sched

////测试计划任务的时间计算
//import(
//    "yundun/crontab_sched/lib/time"
//    "yundun/crontab_sched/lib/sched"
//}
//now := time.Now()
//var crontab string
//crontab = "1 1 * * * *";        fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 4 * * * *";        fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 */1 * * * *";      fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 */4 * * * *";      fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 */7 * * * *";      fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 22,25 * * * *";    fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());
//crontab = "1 20-30 * * * *";    fmt.Println(now.ToStr(), crontab, sched.Create(now, crontab).NextTime());

import (
    "fmt"
    "regexp"
	"strings"
	"strconv"
    "github.com/jingwu15/golib/time"
)

type Sched struct {
	Current time.Time
    Sched string
    SchedFmt string
    Scheds map[string]string
}

func Create(current time.Time, crontab string) Sched {
    schedFmt, scheds := FormatCrontab(crontab)
    //fmt.Println("Sched: ", crontab, schedFmt)
	return Sched {
		Current: current,
        Sched: crontab,
        SchedFmt: schedFmt,
        Scheds: scheds,
	}
}

//计算下一次执行时间
func (c Sched) NextTime() string {
	nextTime := c.CalcNextTime(c.Scheds)
    if c.Current.ToStr() == nextTime {
        //更新当前时间，
        c.Current = time.Unix(c.Current.Unix() + 1, 0)
	    nextTime = c.CalcNextTime(c.Scheds)
    }
    return nextTime
}

//格式化时间成键对值时间格式, 并校正不规范的格式
func FormatCrontab(schedRaw string) (string, map[string]string) {
    scheds := strings.Split(schedRaw, " ")
    keys := []string{"week", "month", "day", "hour", "minute", "second"}
    //标识父级是否己启用，即：值 不为 *
    pflag := false
	schedFmt := make(map[string]string)
    for i, key := range keys {
        v := scheds[5 - i]
        if v == "*/1" { v = "*" }
        //父级已启用， 子级置为初始值，
        if pflag && v == "*" && key == "day" { v = "1" }
        if pflag && v == "*" && (key == "hour" || key == "minute" || key == "second") { v = "0" }

        if key != "week" && v != "*" { pflag = true }
	    schedFmt[key]	= v
    }
    return fmt.Sprintf("%s %s %s %s %s %s", schedFmt["second"], schedFmt["minute"], schedFmt["hour"], schedFmt["day"], schedFmt["month"], schedFmt["week"]), schedFmt
}

//SchedArr 规则时间
func (c Sched) CalcNextTime(SchedArr map[string]string) string{
    current := c.Current
    timeMap := map[string][]int{
        "year":   []int{current.Year(), current.Year() + 1},
        "month":  []int{},
        "week":   []int{},
        "day":    []int{},
        "hour":   []int{},
        "minute": []int{},
        "second": []int{},
    }
    currentMap := map[string]int{
        "year":   current.Year(),
        "month":  current.Month(),
        "week":   current.WeekInt(0),
        "day":    current.Day(),
        "hour":   current.Hour(),
        "minute": current.Minute(),
        "second": current.Second(),
    }
    nextMap := map[string]int{
        "year":   0,
        "month":  0,
        "week":   0,
        "day":    0,
        "hour":   0,
        "minute": 0,
        "second": 0,
    }

	timeMap["month"]  = c.CalcNextMonth(SchedArr["month"])       //月
	timeMap["week"]   = c.CalcNextWeek(SchedArr["week"])         //周
	timeMap["day"]    = c.CalcNextDay(SchedArr["day"])           //日
	timeMap["hour"]   = c.CalcNextHour(SchedArr["hour"])         //时
	timeMap["minute"] = c.CalcNextMinute(SchedArr["minute"])     //分
	timeMap["second"] = c.CalcNextSecond(SchedArr["second"])     //秒
    //nextSecond = timeMap["second"][0]

    incrPkey := ""
    keys := []string{"second", "minute", "hour", "day", "month", "year"}
    for i, key := range keys {
        if timeMap[key][1] < timeMap[key][0] {
            nextMap[key] = timeMap[key][1]
            incrPkey = keys[i+1]
        } else if timeMap[key][1] == timeMap[key][0] {      //周期内只有一个值
            nextMap[key] = timeMap[key][0]
            if currentMap[key] >= timeMap[key][0] {         //当前的时间点已过，父级增长
                incrPkey = keys[i+1]
            }
        } else {
            if incrPkey == key {
                nextMap[key] = timeMap[key][1]
                incrPkey = ""
            } else {
                nextMap[key] = timeMap[key][0]
            }
        }
    }

    nextTimeStr  := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", nextMap["year"], nextMap["month"], nextMap["day"], nextMap["hour"], nextMap["minute"], nextMap["second"])
	nextTimeT, _ := time.StrToTime(nextTimeStr)
	nextTimeUnix := nextTimeT.Unix()       //需要时区

    //星期数据要两次修正            1 30 22 * * 3       周三
    //1. 取到的星期小于正确的星期，应补足                                       在 2019-09-23 12:12:12 周一，取到的下次时间为  2019-09-23 22:30:01 周一,    差两天，补2天
    //2. 取到的星期等于正确的星期，但当前时间早于取到的时间，正确，不处理       在 2019-09-25 12:12:12 周三，取到的下次时间为  2019-09-25 22:30:01 周三,    正确不处理
    //2. 取到的星期晚于正确的星期，应将星期置为下星期                           在 2019-09-25 22:52:12 周三，取到的下次时间为  2019-09-26 22:30:01 周四,    应将星期置为下星期三
	//差额天
    if SchedArr["week"] != "*" {            //只有明确启用星期才有用
        nextTimeWeek := nextTimeT.WeekInt(0)

        if timeMap["week"][0] < timeMap["week"][1] {
            if timeMap["week"][0] < nextTimeWeek {
                nextMap["week"] = timeMap["week"][1]
            } else {
                nextMap["week"] = timeMap["week"][0]
            }
        } else if timeMap["week"][0] == timeMap["week"][1] {
            nextMap["week"] = timeMap["week"][0]
        } else {
            if timeMap["week"][0] > nextTimeWeek {
                nextMap["week"] = timeMap["week"][1]
            } else {
                nextMap["week"] = timeMap["week"][0]
            }
        }

        if nextTimeWeek < nextMap["week"] {            //星期过了
            diffUnix := int64((nextTimeWeek - nextMap["week"]) * 86400)
            nextTimeUnix = nextTimeUnix - diffUnix
        } else if nextMap["week"] == nextTimeWeek {    //星期正确
        } else {                                //星期未到
            diffUnix := int64((7 - nextTimeWeek + nextMap["week"]) * 86400)
            nextTimeUnix = nextTimeUnix + diffUnix
        }
    }

	time.Unix(nextTimeUnix, 0).Format("F")
	return time.Unix(nextTimeUnix,0).Format("F")
}

//月处理
func (c Sched) CalcNextMonth(month string) []int {
	rows := CalcUse(month, 1, 12)
	currentMonth := c.Current.Month()		//月
	if len(rows) == 0 {
		return []int{currentMonth, currentMonth}
	}
	used := FilterMin(rows,currentMonth)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
	//res := rows[0]
	//if len(used) > 0 {
	//	res = used[0]
	//}
	//return res
}

//周处理
func (c Sched) CalcNextWeek(schedWeek string) []int {
	rows := CalcUse(schedWeek, 0, 6)
	currentWeek := c.Current.WeekInt(0)    //周
	if len(rows) == 0 {
		return []int{currentWeek, currentWeek}
	}
	used := FilterMin(rows, currentWeek)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
}


//天处理
func (c Sched) CalcNextDay(scheday string) []int {
	currentDay	:= c.Current.Day()		//日
    monthDayMap := map[int]int{1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
	rows        := CalcUse(scheday, 1, monthDayMap[c.Current.Month()])
	if len(rows) == 0 {
		return []int{currentDay, currentDay}
	}
	used := FilterMin(rows,currentDay)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
	//res := rows[0]
	//if len(used) > 0 {
	//	res = used[0]
	//}
	//return res
}

//时处理
func (c Sched) CalcNextHour(scheHour string) []int {
	currentHour	:= c.Current.Hour()		//时
	rows := CalcUse(scheHour, 0, 23)
	if len(rows) == 0 {
		return []int{currentHour, currentHour}
	}
	used := FilterMin(rows, currentHour)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
	//res := rows[0]
	//if len(used) > 0{
	//	res = used[0]
	//}
	//return res
}

//分处理
func (c Sched) CalcNextMinute(scheMinute string) []int {
	currentMinute	:= c.Current.Minute()			//分
	rows := CalcUse(scheMinute, 0, 59)
	if len(rows) == 0 {
		return []int{currentMinute, currentMinute}
	}
	used := FilterMin(rows,currentMinute)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
	//res := rows[0]
	//if len(used) > 0{
	//	res = used[0]
	//}
	//return res
}

//秒处理
func (c Sched) CalcNextSecond(scheSecond string) []int {
	currentSecond	:= c.Current.Second()			//分
	rows := CalcUse(scheSecond, 0, 59)
	if len(rows) == 0 {
		return []int{currentSecond, currentSecond}
	}
	used := FilterMin(rows,currentSecond)
    //保证至少有两个值, 如果未取到值，则取第一个值
	if len(used) == 0 {
	    used = append(used,rows[0])
    }
	if len(used) == 1 {
	    used = append(used,rows[0])
    }
    return used
	//res := rows[0]
	//if len(used) > 0 {
	//	res = used[0]
	//}
	//return res
}

func CalcUse(schedValue string, start int, end int) []int{
	uses := []int{}
    // 1
    result, _ := regexp.MatchString(`^\\d{1,4}$`, schedValue)
    if result {
        val, err := strconv.Atoi(schedValue)
        if(err == nil) {
            uses = append(uses, val)
        }
    }

    // 1,2,3
    result, _ = regexp.MatchString(`^\d{1,4}(,\d{1,4})?$`, schedValue)
    if result {
        rows := strings.Split(schedValue,",")
        for _,row := range rows {
            val, err := strconv.Atoi(row)
            if(err == nil) {
                uses = append(uses, val)
            }
        }
    }

    //1-6
    result, _ = regexp.MatchString(`^\d{1,4}-\d{1,4}$`, schedValue)
    if result {
        rows := strings.Split(schedValue,"-")
        vStart,errS  := strconv.Atoi(rows[0])
        vEnd, errE   := strconv.Atoi(rows[1])
        if errS == nil && errE == nil {
            for i := vStart; i <= vEnd; i++ {
                uses = append(uses, i)
            }
        }
    }

    // */30
    result, _ = regexp.MatchString(`^*/[1-9]\d{0,4}$`, schedValue)
    if result {
		tmp := schedValue[2:len(schedValue)]
        mul, err := strconv.Atoi(tmp)
        if err == nil {
		    for i:= start; i<=end; i++ {
	            if i % mul == 0{
	                uses = append(uses, i)
	            }
            }
        }
    }
    if schedValue == "*" {
	    for i:= start; i<=end; i++ {
	        uses = append(uses, i)
        }
    }

	return uses
}

func FilterMin(rows []int, current int) []int{
	arr := []int{}
	for i:=0; i<len(rows);i++ {
		if rows[i] >= current {
			arr = append(arr,rows[i])
		}
	}
	return arr
}

func ResetTimeIntArr(items1, items2 []int) []int {
    if items1[1] < items1[0] {
        items2 = append(items2, items2[0])
        items2 = items2[1:]
    }
    return items2
}

//测试计划任务的时间计算
//Tester("1 * 10 * * */1",  "2019-09-24 12:12:12", 10)
func Tester(schedStr, start string, day int) {
	var ttime, now time.Time
    now, e := time.StrToTime(start)
    if e != nil {
        fmt.Println("开始时间错误: ", start)
        return
    }
    lastTime, nextTime := "", ""
    for i := 0; i < (day * 86400); i++ {
        ttime = time.Unix(now.Unix() + int64(i), 0)
        if nextTime == "" {
            nextTime = Create(ttime, schedStr).NextTime()
        }
        if ttime.ToStr() == nextTime {              //执行时间到了
            lastTime = nextTime
            //重新计算下次的执行时间
            schedObj := Create(ttime, schedStr)
            nextTime = schedObj.NextTime()
            fmt.Printf("sched: %s\tschedFmt: %s\tnow: %s\tlast: %s\tnext: %s\texec_done\n", schedStr, schedObj.SchedFmt, ttime.ToStr(), lastTime, nextTime);
        } else {                                    //执行时间未到
            //fmt.Printf("crontab: %s\tnow: %s\tlast: %s\tnext: %s\n", schedStr, ttime.ToStr(), lastTime, nextTime);
        }
    }
}
