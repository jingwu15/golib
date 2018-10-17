package nexttime

import (
	"strings"
	"yundun/crontab_builder/lib/time"
	"strconv"
	"bytes"
)

func GetNextTime(crontabTime string) string{
	//这里还要加个验证
	Sched := strings.Split(crontabTime, " ")
	if len(Sched) >= 5 {
		SchedArr := FormatCrotab(Sched)
		Current := time.Now().Unix()
		//计算下一次执行时间
		nextTime := CalcNextTime(SchedArr, Current)
		return nextTime
	}else{
		return "err"
	}

}

//格式化时间成键对值时间格式
func FormatCrotab(Sched []string) map[string]string {
	response := make(map[string]string)
	response["minute"] 	= Sched[0]
	response["hour"] 	= Sched[1]
	response["day"] 	= Sched[2]
	response["month"] 	= Sched[3]
	response["week"] 	= Sched[4]
	return response
}

/**
SchedArr 规则时间
Current 当前时间
 */
func CalcNextTime(SchedArr map[string]string, Current int64) string{
	currentYear 	:= time.Unix(Current, 0).Format("2006") 		//年
	currentMonth 	:= time.Unix(Current, 0).Format("01")  		//月
	currentDay 		:= time.Unix(Current, 0).Format("02")    		//日
	currentWeek 	:= time.Now().WeekStr(0)								//周
	currentHour 	:= time.Unix(Current, 0).Format("15")   		//时
	currentMinute 	:= time.Unix(Current, 0).Format("04") 			//分
	//currentS := time.Unix(Current,0).Format("05")								//预留
	nextYear 	:= time.Unix(Current, 0).Format("2006") 		//年

	//月
	nextMonth := CalcNexMonth(SchedArr["month"],Current)
	nextMonthInt,err:=strconv.Atoi(nextMonth)
	currentMonthInt,err:=strconv.Atoi(currentMonth)
	currentYearInt,err:=strconv.Atoi(currentYear)
	if nextMonthInt < currentMonthInt && err == nil{
		nextYear = strconv.Itoa(currentYearInt+1)
	}

	//周
	nextWeek := CalcNexWeek(SchedArr["week"],Current)
	nextWeekInt,_ := strconv.Atoi(nextWeek)
	currentWeekInt,_ := strconv.Atoi(currentWeek)

	//日
	nextDay  := CalcNexDay(SchedArr["day"],Current)
	nextDayInt,err:=strconv.Atoi(nextDay)
	currentDayInt,_:=strconv.Atoi(currentDay)
	if nextDayInt < currentDayInt && err == nil{
		nextMonth = strconv.Itoa(nextMonthInt+1)
	}
	//时
	nextHour := CalcNexHour(SchedArr["hour"],Current)
	nextHourInt,err:=strconv.Atoi(nextHour)
	currentHourInt,_:=strconv.Atoi(currentHour)
	if nextHourInt < currentHourInt && err == nil{
		nextDay = strconv.Itoa(nextDayInt+1)
	}
	//分
	nextMinute := CalcNexMinute(SchedArr["minute"],Current)
	nextMinuteInt,err:=strconv.Atoi(nextMinute)
	currentMinuteInt,_:=strconv.Atoi(currentMinute)
	if nextMinuteInt < currentMinuteInt && err == nil{
		nextHour = strconv.Itoa(nextHourInt+1)
	}

	//差额天
	diffDay := 0
	if nextWeekInt != currentWeekInt{
		if nextWeekInt > currentWeekInt{
			diffDay = nextWeekInt - currentWeekInt
		}else if nextWeekInt < currentWeekInt{
			diffDay = 7 - currentWeekInt + nextWeekInt
		}
	}
	//处理时间差
	nextDayInt,_ =strconv.Atoi(nextDay)
	nextDay 	 = strconv.Itoa(nextDayInt+diffDay)

	var mytime bytes.Buffer
	mytime.Write([]byte(nextYear))
	mytime.Write([]byte("-"))
	if len(nextMonth) <= 1 {
		mytime.WriteString("0")
		mytime.WriteString(nextMonth)
	}else {
		mytime.WriteString(nextMonth)
	}
	mytime.Write([]byte("-"))
	mytime.WriteString(nextDay)
	mytime.Write([]byte(" "))
	if len(nextHour) <= 1 {
		mytime.WriteString("0")
		mytime.WriteString(nextHour)
	}else {
		mytime.WriteString(nextHour)
	}
	mytime.Write([]byte(":"))

	if len(nextMinute) <= 1 {
		mytime.WriteString("0")
		mytime.WriteString(nextMinute)
	}else {
		mytime.WriteString(nextMinute)
	}
	mytime.Write([]byte(":"))
	mytime.Write([]byte("01"))
	nextTimeStr := string(mytime.String())

	nextTimeT,_   := time.StrToTime(nextTimeStr)
	nextTimeUnix  := nextTimeT.Unix()//需要时区
	time.Unix(nextTimeUnix,0).Format("F")
	return time.Unix(nextTimeUnix,0).Format("F")
}

/**
	月处理
 */
func CalcNexMonth(month string,current int64) string{
	rows := CalCuse(month,1,12)
	currentMonth := time.Unix(current, 0).Format("01")  		//月
	if len(rows) == 0{
		return currentMonth
	}else{
		used := FilterMin(rows,currentMonth)
		res := rows[0]
		if len(used) > 0{
			res = used[0]
		}
		return res
	}
}

/**
	周处理
 */
func CalcNexWeek(schedWeek string,current int64) string{
	rows := CalCuse(schedWeek,0,6)
	currentMonth := time.Unix(current,0).WeekStr(0)//周
	if len(rows) == 0{
		return currentMonth
	}else{
		used := FilterMin(rows,currentMonth)
		res := rows[0]
		if len(used) > 0{
			res = used[0]
		}
		return res
	}
}


/**
	天处理
 */
func CalcNexDay(scheday string,current int64) string{
	currentDay 		:= time.Unix(current, 0).Format("02")    		//日
	currentDayInt,_   := strconv.ParseInt(currentDay, 10, 64)
	rows := CalCuse(scheday,1,currentDayInt)
	if len(rows) == 0{
		return currentDay
	}else{
		used := FilterMin(rows,currentDay)
		res := rows[0]
		if len(used) > 0{
			res = used[0]
		}
		return res
	}
}

/**
	时处理
 */
func CalcNexHour(scheHour string,current int64) string{
	currentHour 	:= time.Unix(current, 0).Format("15")   		//时
	rows := CalCuse(scheHour,0,23)
	if len(rows) == 0{
		return currentHour
	}else{
		used := FilterMin(rows,currentHour)
		res := rows[0]
		if len(used) > 0{
			res = used[0]
		}
		return res
	}
}

/**
	分处理
 */
func CalcNexMinute(scheMinute string,current int64) string{
	currentMinute 	:= time.Unix(current, 0).Format("04") 			//分
	rows := CalCuse(scheMinute,1,59)
	if len(rows) == 0{
		return currentMinute
	}else{
		used := FilterMin(rows,currentMinute)
		res := rows[0]
		if len(used) > 0{
			res = used[0]
		}
		return res
	}
}

func CalCuse(schedValue string,star int64 ,end int64) []string{
	response := []string{}
	//原是先判断是否是数字/数字字符串，现放在else里面
	if strings.Count(schedValue,",") > 0{
		response = strings.Split(schedValue,",")
	}else if strings.Count(schedValue,"-") > 0{
		response = strings.Split(schedValue,"-")
	}else if strings.Count(schedValue,"/") > 0{
		multipleS := schedValue[2:len(schedValue)]
		multiple,_ := strconv.ParseInt(multipleS, 10, 64)
		for i:= star;i<=end ;i++  {
			if i % multiple == 0{
			response = append(response,strconv.FormatInt(i,10))
			}
		}
	}else if _, err := strconv.Atoi(schedValue); err == nil{
		response = append(response,schedValue)
	}else{

	}
	return response
}

func FilterMin(rows []string,currentMonth string) []string{
	currentMonthInt,_ := strconv.Atoi(currentMonth)
	arr := []string{}
	for i:=0; i<len(rows)-1;i++  {
		rowsInt,_ := strconv.Atoi(rows[i])
		if rowsInt > currentMonthInt{
			arr = append(arr,rows[i])
		}
	}
	return arr
}
