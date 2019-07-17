package tconv

import (
    "fmt"
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

func AtoiForce(s string) int {
    i := 0
    i, _ = strconv.Atoi(s)
    return i
}

func StrToIntForce(s string) int {
    i := 0
    i, _ = strconv.Atoi(s)
    return i
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

func KeysMapIF(data map[string]interface{}) []string {
	keys := []string{}
	for key, _ := range data {
        keys = append(keys, key)
	}
	return keys
}

func KeysMapStr(data map[string]string) []string {
	keys := []string{}
	for key, _ := range data {
        keys = append(keys, key)
	}
	return keys
}

//string列表转换成字典
func KeyStrToMap(keys []string) map[string]int {
    var m = map[string]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

//int列表转换成字典
func KIntToMap(keys []int) map[int]int {
    var m = map[int]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

//int64列表转换成字典
func KInt64ToMap(keys []int64) map[int64]int {
    var m = map[int64]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}


//float64列表转换成字典
func KFloat64ToMap(keys []float64) map[float64]int {
    var m = map[float64]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

func KeysMapInt64(data map[string]int64) []string {
	keys := []string{}
	for key, _ := range data {
        keys = append(keys, key)
	}
	return keys
}

//K str V str转换成字典
func KStrVStrToMap(keys, values []string) map[string]string {
	var m = map[string]string{}
	for i, k := range keys {
		m[k] = values[i]
	}
	return m
}

func FormatValues(items interface{}, format string) (keys []string) {
    keys = []string{}
    switch items.(type) {
    case []string:
        for _, v := range items.([]string) {
            keys = append(keys, fmt.Sprintf(format, v))
        }
    case []int:
        for _, v := range items.([]int) {
            keys = append(keys, fmt.Sprintf(format, v))
        }
    case []float64:
        for _, v := range items.([]float64) {
            keys = append(keys, fmt.Sprintf(format, v))
        }
    default:
    }
    return keys
}

