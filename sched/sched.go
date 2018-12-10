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
    Crontab string
}

func Create(current time.Time, crontab string) Sched {
	return Sched {
		Current: current,
        Crontab: crontab,
	}
}

//计算下一次执行时间
func (c Sched) NextTime() string {
	sched := strings.Split(c.Crontab, " ")
	schedArr := FormatCrontab(sched)
	nextTime := c.CalcNextTime(schedArr)
	return nextTime
}

//格式化时间成键对值时间格式
func FormatCrontab(sched []string) map[string]string {
	response := make(map[string]string)
	response["second"]	= sched[0]
	response["minute"]	= sched[1]
	response["hour"]	= sched[2]
	response["day"]	    = sched[3]
	response["month"]	= sched[4]
	response["week"]	= sched[5]
	return response
}

//SchedArr 规则时间
func (c Sched) CalcNextTime(SchedArr map[string]string) string{
    current := c.Current
	currentYear	    := current.Year()		    //年
	currentMonth	:= current.Month()		    //月
	currentDay		:= current.Day()		    //日
	currentWeek	    := current.WeekInt(0)	    //周
	currentHour	    := current.Hour()		    //时
	currentMinute	:= current.Minute()		    //分
	currentSecond   := current.Second()		    //预留
	nextYear	    := current.Year()		    //年

	//月
	nextMonth := c.CalcNextMonth(SchedArr["month"])
	if nextMonth < currentMonth && nextYear == currentYear {
		nextYear = currentYear + 1
	}
	//周
	nextWeek := c.CalcNextWeek(SchedArr["week"])
	//日
	nextDay  := c.CalcNextDay(SchedArr["day"])
	if nextDay < currentDay && nextMonth == currentMonth {
		nextMonth = nextMonth + 1
	}
	//时
	nextHour := c.CalcNextHour(SchedArr["hour"])
	if nextHour < currentHour && nextDay == currentDay {
		nextDay = nextDay + 1
	}
	//分
	nextMinute := c.CalcNextMinute(SchedArr["minute"])
	if nextMinute < currentMinute && nextHour == currentHour {
		nextHour = nextHour + 1
	}
	//秒
	nextSecond := c.CalcNextSecond(SchedArr["second"])
	if nextSecond < currentSecond && nextMinute == currentMinute {
		nextMinute = nextMinute + 1
	}

	//差额天
	diffDay := 0
	if nextWeek != currentWeek {
		if nextWeek > currentWeek {
			diffDay = nextWeek - currentWeek
		} else if nextWeek < currentWeek {
			diffDay = 7 - currentWeek + nextWeek
		}
	}
	//处理时间差
	nextDay	= nextDay + diffDay

    nextTimeStr  := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", nextYear, nextMonth, nextDay, nextHour, nextMinute, nextSecond)
	nextTimeT,_  := time.StrToTime(nextTimeStr)
	nextTimeUnix := nextTimeT.Unix()       //需要时区
	time.Unix(nextTimeUnix, 0).Format("F")
	return time.Unix(nextTimeUnix,0).Format("F")
}

//月处理
func (c Sched) CalcNextMonth(month string) int {
	rows := CalcUse(month, 1, 12)
	currentMonth := c.Current.Month()		//月
	if len(rows) == 0 {
		return currentMonth
	}
	used := FilterMin(rows,currentMonth)
	res := rows[0]
	if len(used) > 0 {
		res = used[0]
	}
	return res
}

//周处理
func (c Sched) CalcNextWeek(schedWeek string) int {
	rows := CalcUse(schedWeek, 0, 6)
	currentMonth := c.Current.WeekInt(0)    //周
	if len(rows) == 0 {
		return currentMonth
	}
	used := FilterMin(rows,currentMonth)
	res := rows[0]
	if len(used) > 0 {
		res = used[0]
	}
	return res
}


//天处理
func (c Sched) CalcNextDay(scheday string) int {
	currentDay	:= c.Current.Day()		//日
    monthDayMap := map[int]int{1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31}
	rows        := CalcUse(scheday,1, monthDayMap[c.Current.Month()])
	if len(rows) == 0 {
		return currentDay
	}
	used := FilterMin(rows,currentDay)
	res := rows[0]
	if len(used) > 0 {
		res = used[0]
	}
	return res
}

//时处理
func (c Sched) CalcNextHour(scheHour string) int {
	currentHour	:= c.Current.Hour()		//时
	rows := CalcUse(scheHour, 0, 23)
	if len(rows) == 0 {
		return currentHour
	}
	used := FilterMin(rows, currentHour)
	res := rows[0]
	if len(used) > 0{
		res = used[0]
	}
	return res
}

//分处理
func (c Sched) CalcNextMinute(scheMinute string) int {
	currentMinute	:= c.Current.Minute()			//分
	rows := CalcUse(scheMinute, 0, 59)
	if len(rows) == 0 {
		return currentMinute
	}
	used := FilterMin(rows,currentMinute)
	res := rows[0]
	if len(used) > 0{
		res = used[0]
	}
	return res
}

//秒处理
func (c Sched) CalcNextSecond(scheSecond string) int {
	currentSecond	:= c.Current.Second()			//分
	rows := CalcUse(scheSecond, 0, 59)
	if len(rows) == 0 {
		return currentSecond
	}
	used := FilterMin(rows,currentSecond)
	res := rows[0]
	if len(used) > 0 {
		res = used[0]
	}
	return res
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

	return uses
}

func FilterMin(rows []int,currentMonth int) []int{
	arr := []int{}
	for i:=0; i<len(rows)-1;i++ {
		if rows[i] > currentMonth {
			arr = append(arr,rows[i])
		}
	}
	return arr
}

