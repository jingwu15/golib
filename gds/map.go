package gds

import(
    "fmt"
    "reflect"
)

//取交集
func Map_key_inner(d0 interface{}, d1 interface{}) (data interface{}, err error) {
    t0 := reflect.TypeOf(d0).String()
    t1 := reflect.TypeOf(d1).String()
    if t0 != t1 { return data, fmt.Errorf("数据类型不一致：%s, %s", t0, t1) }
    switch t0 {
    case "map[int]int":
        cm := []int{}
        dm1 := d1.(map[int]int)
        for k0, _ := range d0.(map[int]int) {
            if _, ok := dm1[k0]; ok {
                cm = append(cm, k0)
            }
        }
        return cm, nil
    case "map[string]int":
        cm := []string{}
        dm1 := d1.(map[string]int)
        for k0, _ := range d0.(map[string]int) {
            if _, ok := dm1[k0]; ok {
                cm = append(cm, k0)
            }
        }
        return cm, nil
    default:
    }
    return data, nil
}

//k/v 反转
func Map_vkMap(data interface{}) (ret interface{}, err error) {
    t0 := reflect.TypeOf(data).String()
    switch t0 {
    case "map[string]string":
        ret := map[string]string{}
        for k, v := range data.(map[string]string) {
            ret[v] = k
        }
        return ret, nil
    default:
    }
    return data, fmt.Errorf("数据类型支持")
}

//map 取 key 列表
func Map_kList(data interface{}) (ret interface{}, err error) {
    t0 := reflect.TypeOf(data).String()
    switch t0 {
    case "map[string]string", "map[string]int", "map[string]interface{}":
        ret := []string{}
        for k, _ := range data.(map[string]string) {
            ret = append(ret, k)
        }
        return ret, nil
    default:
    }
    return data, fmt.Errorf("数据类型支持")
}

//map 取 value 列表
func Map_vList(data interface{}) (ret interface{}, err error) {
    t0 := reflect.TypeOf(data).String()
    switch t0 {
    case "map[string]string", "map[int]string":
        ret := []string{}
        for _, v := range data.(map[string]string) {
            ret = append(ret, v)
        }
        return ret, nil
    default:
    }
    return data, fmt.Errorf("数据类型支持")
}

func Map_kStrs(data interface{}) (rs []string, e error) {
    ret, e := Map_kList(data)
    if e != nil { return rs, e }
    return ret.([]string), nil
}

func Map_vStrs(data interface{}) (rs []string, e error) {
    ret, e := Map_vList(data)
    if e != nil { return rs, e }
    return ret.([]string), nil
}

