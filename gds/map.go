package gds

import(
    "fmt"
    "reflect"
    "strings"
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
    return data, fmt.Errorf("数据类型不支持")
}

//map 取 key 列表
func Map_kList(data interface{}) (ret interface{}, err error) {
    t0 := reflect.TypeOf(data).String()
    fmt.Println("t0----", t0)
    switch t0 {
    case "map[string]string":
        ret := []string{}
        for k, _ := range data.(map[string]string) { ret = append(ret, k) }
        return ret, nil
    case "map[string]int":
        ret := []string{}
        for k, _ := range data.(map[string]int) { ret = append(ret, k) }
        return ret, nil
    case "map[string]interface{}", "map[string]interface {}":
        ret := []string{}
        for k, _ := range data.(map[string]interface{}) { ret = append(ret, k) }
        return ret, nil
    default:
    }
    return data, fmt.Errorf("数据类型不支持")
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
    return data, fmt.Errorf("数据类型不支持")
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

//合并map, 不支持自定义的结构体
//a := map[string]interface{}{
//    "db": map[string]interface{}{
//        "dns": map[string]interface{}{
//            "host": "test.yundun.com",
//            "port": "3306",
//        },
//    },
//    "tables": []string{"table1", "table2", "table3"},
//    "order_status": map[int]string{1: "成功", 2: "失败"},
//    "invoice_status": map[string]string{"1": "成功", "2": "失败"},
//    "nil_t": nil,
//}
//b := map[string]interface{}{
//    "db": map[string]interface{}{
//        "dns": map[string]interface{}{
//            "host": "127.0.0.1",
//            "dbname": "dns",
//        },
//        "cp": map[string]interface{}{
//            "host": "test.yundun.com",
//            "port": "3306",
//        },
//    },
//    "tables": []string{"table1", "table2", "table3", "table4", "table5", "table5"},
//    "order_status": map[int]string{1: "成功", 2: "失败", 3:"支付失败"},
//    "invoice_status": map[string]string{"1": "成功", "2": "失败", "3":"支付失败"},
//    "nil_t": nil,
//}
//result := gds.Map_merge(&a, &b)
//fmt.Println(result)
func Map_merge(raw interface{}, ext interface{}) map[string]interface{} {
    rawType := reflect.TypeOf(raw).String()
    extType := reflect.TypeOf(ext).String()
    if rawType == extType {
        switch(rawType) {
        case "map[string]interface {}", "*map[string]interface {}":
            rawData := map[string]interface{}{}
            extData := map[string]interface{}{}
            if(rawType == "map[string]interface {}") {
                rawData = (raw).(map[string]interface{})
                extData = (ext).(map[string]interface{})
            } else {
                rawDataP := (raw).(*map[string]interface{})
                extDataP := (ext).(*map[string]interface{})
                rawData = *rawDataP
                extData = *extDataP
            }
            keysRaw, _ := Map_kStrs(rawData)
            keysExt, _ := Map_kStrs(extData)
            keys := keysRaw
            for _,key := range keysExt {
                keys = append(keys, key)
            }
            keysMap := List_strMap(keys)
            keys, _ = Map_kStrs(keysMap)
            for _, k := range keys {
                if _, ok := extData[k]; !ok { continue }                            //ext 中不存在key, 跳过
                if _, ok := rawData[k]; !ok { rawData[k] = extData[k]; continue; }  //raw 中不存在key, 直接赋值
                if rawData[k] == nil || extData[k] == nil { continue; }

                rawVType := reflect.TypeOf(rawData[k]).String()
                extVType := reflect.TypeOf(extData[k]).String()
                if(rawVType != extVType) { rawData[k] = extData[k]; continue;    }  //都存在key, 但类型不一样，赋新值

                //都存在key, 类型一样
                vType := ""
                if strings.HasPrefix(rawVType, "[]") {
                    vType = strings.TrimLeft(rawVType, "[]")         //数组化简
                    fmt.Println("rawvtype----", rawVType, "vtype----", vType)
                } else if strings.HasPrefix(rawVType, "map[") {     //map, 取值的数据类型
                    tmp := strings.SplitN(rawVType, "]", 2)
                    vType = fmt.Sprintf("m_%s", tmp[1])
                    //fmt.Println("mmm---", vType)
                } else {
                    vType = rawVType
                }
                switch(vType) {
                case "m_interface {}":         //复合数据类型
                    rawData[k] = Map_merge(rawData[k], extData[k])
                case "byte","string", "bool", "m_byte","m_string", "m_bool":                                       //基本数据类型
                    rawData[k] = extData[k]
                case "int", "int8", "int16", "int32", "int64", "m_int", "m_int8", "m_int16", "m_int32", "m_int64":                      //基本数据类型
                    rawData[k] = extData[k]
                case "uint", "uint8", "uint16", "uint32", "uint64", "m_uint", "m_uint8", "m_uint16", "m_uint32", "m_uint64":                 //基本数据类型
                    rawData[k] = extData[k]
                case "float32", "float64", "complex64", "complex128", "m_float32", "m_float64", "m_complex64", "m_complex128":               //基本数据类型
                    rawData[k] = extData[k]
                default:                                                            //未检查出来的数据类型
                    fmt.Println("不支持的数据类型----------key: ", k, "type: ", rawVType)
                }
            }
            return rawData
        default:
            fmt.Println("不支持的数据类型----------", rawType)
        }
        return raw.(map[string]interface{})
    } else {
        fmt.Println("else")
        //直接赋值
        return ext.(map[string]interface{})
    }
}

