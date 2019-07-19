package gds
//基本数据类型操作：数组或列表

import (
    "fmt"
)

//取数组的交集
func List_intersect(items0 interface{}, items1 interface{}) (data interface{}, err error) {
    t0 := GetType(items0)
    t1 := GetType(items1)
    if t0 != t1 {
        return data, fmt.Errorf("数据类型不一致：%s, %s", t0, t1)
    }
    switch t0 {
    case "[]int":
        return Map_key_intersect(List_intToMap(items0.([]int)), List_intToMap(items1.([]int)))
    case "[]string":
        return Map_key_intersect(List_strToMap(items0.([]string)), List_strToMap(items1.([]string)))
    default:
    }
    return data, nil
}

//int列表转换成字典
func List_intToMap(keys []int) map[int]int {
    var m = map[int]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

//int列表转换成字典
func List_strToMap(keys []string) map[string]int {
    var m = map[string]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

