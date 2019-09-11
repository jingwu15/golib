package gds
//基本数据类型操作：数组或列表

import (
    "fmt"
    //"reflect"
)

//取数组的交集
func List_inner(itemss ...interface{}) (data interface{}, err error) {
    if len(itemss) < 2 { return data, fmt.Errorf("参数必须多于2个") }
    dtype := GetType(itemss[0])
    switch dtype {
    case "[]int":
        dcom, itemss := itemss[0], itemss[1:]
        for _, items := range itemss {
            dcom, err = Map_key_inner(List_intMap(dcom.([]int)), List_intMap(items.([]int)))
            if err != nil { return dcom, err }
        }
        return dcom, nil
    case "[]string":
        dcom, itemss := itemss[0], itemss[1:]
        for _, items := range itemss {
            dcom, err = Map_key_inner(List_strMap(dcom.([]string)), List_strMap(items.([]string)))
            if err != nil { return dcom, err }
        }
        return dcom, nil
    default:
    }
    return data, nil
}

func List_inner_strs(itemss ...interface{}) (ret []string, err error) {
    data, err := List_inner(itemss...)
    if err != nil { return ret, err }
    ret, ok := data.([]string)
    if !ok { return ret, fmt.Errorf("不能转换为[]string类型") }
    return ret, nil
}

//int列表转换成字典
func List_intMap(keys []int) map[int]int {
    var m = map[int]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

//int列表转换成字典
func List_strMap(keys []string) map[string]int {
    var m = map[string]int{}
	for _, i := range keys {
		m[i] = 1
	}
	return m
}

