package gds
//基本数据类型操作：数组或列表

import (
    //"fmt"
    "strings"
    //"reflect"
)

//字符串，拆分为字典，如：a,b,c 拆分为字典为： {"a":1, "b":1, "c": 1}
func Str_strMap(raw, sep string, trims ...string) map[string]int {
    list := strings.Split(raw, sep)
    m := map[string]int{}
	for _, i := range list {
		m[i] = 1
	}
	return m
}

//字符串，拆分为数组，如：a,b,c 拆分为数组为： {"a", "b", "c"}
func Str_strs(raw, sep string, trims ...string) []string {
    return strings.Split(raw, sep)
}

