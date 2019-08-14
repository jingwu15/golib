package gds

import (
    "fmt"
    "net"
    "regexp"
    "strings"
)

//通用验证库，基于map结构(不能有struct)
//type  支持验证的数据类型，此处的数据类型不是golang的数据类型， 仅支持 int/str/list/map
//vstr  验证的字符串，与验证函数是一一对应的，如：intval/range/notEmpty
//使用方式：
//    d := map[string]interface{}{
//        "id": 26,
//        "ids": []int{6, 7},
//        "rules": []map[string]interface{}{
//            map[string]interface{}{"action": "anticc", "sort": 10},
//            map[string]interface{}{"action": "xxxxx", "sort": 1000},
//        },
//    }
//    cfgs := map[string]interface{}{
//        "id": map[string]interface{}{"name": "id", "vstr": "intval/range", "type": "int", "range_data": []int{35, 60}},
//        "ids": map[string]interface{}{"name": "ids", "vstr": "listInt", "type": "list"},
//        "rules": map[string]interface{}{"name": "rules", "vstr": "notEmpty/unique", "type": "list", "child": map[string]interface{}{
//                                                                                                        "action": map[string]interface{}{"type": "str"},
//                                                                                                        "sort": map[string]interface{}{"type": "int", "vstr": "intval"},
//                                                                                                    },
//        },
//    }
//    dv, es := Valid(d, cfgs)
//    fmt.Println("dv--------", dv, es)

//验证函数映射
func fMap(fkey string) (fun interface{}, e error) {
    funs := map[string]interface{}{
        "in": V_in,
        "gt0": V_gt0,
        "isIp": V_isIp,
        "range": V_range,
        "intval": V_intval,
        "noEmpty": V_notEmpty,
        "notEmpty": V_notEmpty,
        "isDomain": V_isDomain,
        "isNumber": V_isNumber,
    }

    if fun, ok := funs[fkey]; ok {
        return fun, nil
    } else {
        return nil, fmt.Errorf("函数不存在")
    }
}
//大于0
func V_gt0(value interface{}) (e error) {
    dtype := GetType(value)
    if dtype != "int" && dtype != "float64"      { return fmt.Errorf("数据类型应为整型或浮点型") }
    if dtype == "int" && value.(int) <= 0        { return fmt.Errorf("数据应大于0")              }
    if dtype == "float64" && value.(float64) > 0 { return fmt.Errorf("数据应大于0")              }
    return nil
}

//是否为数字
func V_isNumber(value interface{}) (e error) {
    if GetType(value) != "string" { return fmt.Errorf("域名应为字符串") }
    NotLine := "^\\d{1,}$"
    match, e := regexp.MatchString(NotLine, value.(string))
    if !match { return fmt.Errorf("必须为数字") }
    return nil
}

//是否为域名
func V_isDomain(value interface{}) (e error) {
    if GetType(value) != "string" { return fmt.Errorf("域名应为字符串") }
   NotLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}"
   match, _ := regexp.MatchString(NotLine, value.(string))
   if !match { return fmt.Errorf("域名格式错误") }
   return nil
}

func V_isIp(value interface{}) (e error) {
    if GetType(value) != "string" { return fmt.Errorf("IP应为字符串") }
    ip := net.ParseIP(value.(string))
    if ip == nil { return fmt.Errorf("IP格式错误") }
    return nil
}

func V_intval(value interface{}) (d int, e error) {
    switch GetType(value) {
    case "int":
        return value.(int), nil
    case "string":
        return I_int(value), nil
    case "float64":
        return I_int(value), nil
    default:
        return d, fmt.Errorf("数据不能转为整型!")
    }
}

func V_in(value interface{}, vd interface{}) (e error) {
    vdType := GetType(vd)
    flag := false
    switch GetType(value) {
    case "int":
        if vdType == "[]int" {
            for _, v := range vd.([]int) {
                if v == value { flag = true }
            }
        }
    case "string":
        if vdType == "[]string" {
            for _, v := range vd.([]string) {
                if v == value { flag = true }
            }
        }
    case "float64":
        if vdType == "[]float64" {
            for _, v := range vd.([]float64) {
                if v == value { flag = true }
            }
        }
    default:
        return fmt.Errorf("不支持的数据类型")
    }
    if !flag { return fmt.Errorf("数据超出地范围：[%s]", strings.Join(I_strs(vd), ",")) }
    return nil
}

func V_range(value interface{}, vd interface{}) (e error) {
    v := vd.([]int)
    if value.(int) < v[0] || value.(int) > v[1] {
        return fmt.Errorf("数值䞖界，范围：[%d-%d]", v[0], v[1])
    } else {
        return nil
    }
}

func V_notEmpty(value interface{}) (e error) {
    switch GetType(value) {
    case "string":
        if value.(string) == "" { return fmt.Errorf("数值不能为空!") }
    case "[]string":
        if len(value.([]string)) == 0 { return fmt.Errorf("数值不能为空!") }
    case "[]int":
        if len(value.([]int)) == 0 { return fmt.Errorf("数值不能为空!") }
    case "[]float64":
        if len(value.([]float64)) == 0 { return fmt.Errorf("数值不能为空!") }
    case "[]interface{}":
        if len(value.([]interface{})) == 0 { return fmt.Errorf("数值不能为空!") }
    case "[]map[string]interface{}":
        if len(value.([]map[string]interface{})) == 0 { return fmt.Errorf("数值不能为空!") }
    default:
    }
    return nil
}

//验证数据 数值(int/str) 列表(int/str/dict) 字典(int_i/str_i)
//设置默认值 default
//必须 must/option
/*
{
    "type": "int",
    "vstr": "list",
    "data_default": sss
}
*/
func Valid(params map[string]interface{}, cfgs map[string]interface{}) (data map[string]interface{}, es []string) {
    //转换配置参数
    data, es = map[string]interface{}{}, []string{}
    for field, row := range cfgs {
        cfg := row.(map[string]interface{})
        if v, ok := params[field]; ok {
            dOne, e := ValidOne(v, cfg)
            if e != nil { es = append(es, fmt.Sprintf("%s %s", field, e.Error())); continue }
            data[field] = dOne
        } else {
            vs := strings.Split(cfg["vstr"].(string), "/")
            vMap := List_strMap(vs)
            //must 必要参数
            if _, ok := vMap["must"]; ok { es = append(es, fmt.Sprintf("%s 为必要参数，请补全！", field)); continue }
            //默认值，只有不存在时才设置
            if vd, ok := cfg["default"]; ok { data[field] = vd }
        }
    }

    return data, es
}

func ValidOne(data interface{}, cfg map[string]interface{}) (d interface{}, e error) {
    //数据类型检测
    dtype := GetType(data)
    //类型转换，为了开发时方便，将其他类型转为需要的类型
    //转换 []interface{} 类型为真实的类型
    if dtype == "[]interface{}"      { _, data = I_list(data) }
    if cfg["type"].(string) == "int" { data = I_int(data)      }
    if cfg["type"].(string) == "str" { data = I_str(data)      }
    dtype = GetType(data)

    //设置默认值
    //列表空值，设置默认值
    if _, ok := cfg["default"]; ok && dtype == "[]interface{}" && len(data.([]interface{})) == 0 { data = cfg["default"] }

    //类型检查
    if cfg["type"].(string) == "int"  && dtype != "int"                   { return data, fmt.Errorf("数据类型不正确, 要求为整数！")   }
    if cfg["type"].(string) == "str"  && dtype != "string"                { return data, fmt.Errorf("数据类型不正确, 要求为字符串！") }
    if cfg["type"].(string) == "list" && !strings.HasPrefix(dtype, "[]")  { return data, fmt.Errorf("数据类型不正确, 要求为数组！")   }
    if cfg["type"].(string) == "map"  && !strings.HasPrefix(dtype, "map") { return data, fmt.Errorf("数据类型不正确, 要求为字典！")   }

    //if _, ok := cfg["vstr"]; !ok { return data, fmt.Errorf("未设置验证方式！") }        //验证方式可以不设置
    vMap := map[string]int{}
    if vstr, ok := cfg["vstr"]; ok {
        vs := strings.Split(vstr.(string), "/")
        vMap = List_strMap(vs)
    }
    delete(vMap, "must")
    delete(vMap, "default")
    //处理常规验证
    for fkey, _ := range vMap {
        var fun interface{}
        isListFun := false
        if strings.HasPrefix(fkey, "list_") {           //列表
            flkey := strings.TrimPrefix(fkey, "list_")
            fun, e = fMap(flkey)
            isListFun = true
        } else {
            fun, e = fMap(fkey)
        }
        if e != nil { continue; return data, fmt.Errorf("%s 验证函数不存在", fkey) }
        switch GetTypeFun(fun) {
        case "func(interface{})(int, error)":
            if isListFun {
                ds, rs := []interface{}{}, []int{}
                if dtype == "[]int"         { for _, v := range data.([]int)         { ds = append(ds, v) } }
                if dtype == "[]string"      { for _, v := range data.([]string)      { ds = append(ds, v) } }
                if dtype == "[]float64"     { for _, v := range data.([]float64)     { ds = append(ds, v) } }
                for _, d := range ds {
                    f := fun.(func(interface{})(int, error))
                    r, e := f(d)
                    if e != nil { return ds, e }
                    rs = append(rs, r)
                }
                data = rs
            } else {
                f := fun.(func(interface{})(int, error))
                data, e = f(data)
                if e != nil { return data, e }
            }
        case "func(interface{})([]int, error)":
            f := fun.(func(interface{})([]int, error))
            data, e = f(data)
            if e != nil { return data, e }
        case "func(interface{})(error)":
            if isListFun {
                ds := []interface{}{}
                if dtype == "[]int"         { for _, v := range data.([]int)         { ds = append(ds, v) } }
                if dtype == "[]string"      { for _, v := range data.([]string)      { ds = append(ds, v) } }
                if dtype == "[]float64"     { for _, v := range data.([]float64)     { ds = append(ds, v) } }
                for _, d := range ds {
                    f := fun.(func(interface{})(error))
                    e = f(d)
                    if e != nil { return data, e }
                }
            } else {
                f := fun.(func(interface{})(error))
                e = f(data)
                if e != nil { return data, e }
            }
        case "func(interface{}, interface{})(error)":
            f := fun.(func(interface{}, interface{})(error))
            v, ok := cfg[fkey+"_data"]
            if !ok { return data, fmt.Errorf("%s 未提供验证参数", fkey) }
            e = f(data, v)
            if e != nil { return data, e }
        default:
            //return data, fmt.Errorf("不支持的验证函数类型")
        }
    }
    //处理字典， 子节点验证
    if _, ok := cfg["child"]; ok && cfg["type"].(string) == "list" {
        switch dtype {
        case "[]interface{}":
            //ds := []interface{}{}
            ds := []map[string]interface{}{}
            for _, v := range data.([]interface{}) {
                d, es := Valid(v.(map[string]interface{}), cfg["child"].(map[string]interface{}))
                if len(es) > 0 { return data, fmt.Errorf(strings.Join(es, ", ")) }
                ds = append(ds, d)
            }
            return ds, nil
        case "[]map[string]interface{}":
            ds := []map[string]interface{}{}
            for _, v := range data.([]map[string]interface{}) {
                d, es := Valid(v, cfg["child"].(map[string]interface{}))
                if len(es) > 0 { return data, fmt.Errorf(strings.Join(es, ", ")) }
                ds = append(ds, d)
            }
            return ds, nil
        default:
        }
    }
    //处理字典
    return data, nil
}

//func main() {
//    d := map[string]interface{}{
//        "id": 26,
//        "ids": []int{6, 7},
//        "rules": []map[string]interface{}{
//            map[string]interface{}{
//                "action": "anticc",
//                "sort": 10,
//            },
//            map[string]interface{}{
//                "action": "xxxxx",
//                "sort": 1000,
//            },
//        },
//    }
//    cfgs := map[string]interface{}{
//        "id": map[string]interface{}{
//            "name": "id",
//            "vstr": "intval/range",
//            "type": "int",
//            "range_data": []int{35, 60},
//        },
//        "ids": map[string]interface{}{
//            "name": "ids",
//            "vstr": "listInt",
//            "type": "list",
//        },
//        "rules": map[string]interface{}{
//            "name": "rules",
//            "vstr": "notEmpty/unique",        //非空/去重
//            "type": "list",
//            "child": map[string]interface{}{
//                "action": map[string]interface{}{
//                    "type": "str",
//                },
//                "sort": map[string]interface{}{
//                    "type": "int",
//                },
//            },
//        },
//    }
//    dv, es := Valid(d, cfgs)
//    fmt.Println("dv--------", dv, es)
//}

