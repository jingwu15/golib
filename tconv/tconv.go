package tconv

import (
	"bytes"
	"strconv"
)

func ParseBool(str string) (value bool, err error) {
	return strconv.ParseBool(str)
}

func ParseInt(s string, base int, bitSize int) (i int64, err error) {
	return strconv.ParseInt(s, base, bitSize)
}

func ParseUint(s string, base int, bitSize int) (n uint64, err error) {
	return strconv.ParseUint(s, base, bitSize)
}

func ParseFloat(s string, bitSize int) (f float64, err error) {
	return strconv.ParseFloat(s, bitSize)
}

func FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

func FormatInt(i int64, base int) string {
	return strconv.FormatInt(i, base)
}

func FormatUint(i uint64, base int) string {
	return strconv.FormatUint(i, base)
}

func FormatFloat(f float64, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(f, fmt, prec, bitSize)
}

func Atoi(s string) (i int, err error) {
	return strconv.Atoi(s)
}

func StrToInt(s string) (i int, err error) {
	return strconv.Atoi(s)
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func IntsToStrs(rows []int) []string {
    strs := []string{}
    for _, row := range rows {
        strs = append(strs, strconv.Itoa(row))
    }
	return strs
}

//将[][]byte 转为json列表,如：[][]byte{[]byte(`a`), []byte(`a`)} => ["a","b"]
func ByteList_to_json(rawList [][]byte) []byte {
	var items = make([][]byte, 3)
	items[0] = []byte(`[`)
	items[1] = bytes.Join(rawList, []byte(`,`))
	items[2] = []byte(`]`)
	return bytes.Join(items, []byte(""))
}
