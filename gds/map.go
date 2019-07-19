package gds

import(
    "fmt"
    "reflect"
)

//取交集
func Map_key_intersect(d0 interface{}, d1 interface{}) (data interface{}, err error) {
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
