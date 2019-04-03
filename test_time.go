package main

import (
	"fmt"
	"github.com/jingwu15/golib/time"
)

func main() {
    var err error

	fmt.Println("=========================time.week============================")
	fmt.Println(time.Now().Weekday())
	fmt.Println(time.Now().WeekInt(0))
	fmt.Println(time.Now().WeekInt(1))
	fmt.Println(time.Now().WeekStr(0))
	fmt.Println(time.Now().WeekStr(1))
	fmt.Println("=========================time.week============================")

	fmt.Println("=========================time.StrToTime=======================")
	var t time.Time
	t, err = time.StrToTime("2018-01-01 12:34:56")
	fmt.Println(t, err)
	t, err = time.StrToTime("2018-11-21")
	fmt.Println(t, err)
	t, err = time.StrToTime("20181121")
	fmt.Println(t, err)
	t, err = time.StrToTime("2018-1-2")
	fmt.Println(t, err)
	t, err = time.StrToTime("12:34:56")
	fmt.Println(t, err)
	t, err = time.StrToTime("123456")
	fmt.Println(t, err)
	t, err = time.StrToTime("Jan 2 15:04:05 2016")
	fmt.Println(t, err) //ANSIC
	t, err = time.StrToTime("Jan 2 15:04:05 AWST 2016")
	fmt.Println(t, err) //UnixDate
	t, err = time.StrToTime("Jan 02 15:04:05 +0800 2016")
	fmt.Println(t, err) //RubyDate
	t, err = time.StrToTime("02 Jan 26 15:04 MST")
	fmt.Println(t, err) //RFC822  "02 Jan 06 15:04 MST"
	t, err = time.StrToTime("02 Jan 2126 15:04:05 MST")
	fmt.Println(t, err) //RFC1123 "02 Jan 2006 15:04:05 MST"
	t, err = time.StrToTime("02 Jan 06 15:04 +0800")
	fmt.Println(t, err) //RFC822Z
	t, err = time.StrToTime("02 Jan 2106 15:04:05 +0800")
	fmt.Println(t, err) //RFC1123Z
	fmt.Println("=========================time.StrToTime=======================")

}
