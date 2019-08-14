package gds

import (
    "fmt"
    "strconv"
    "strings"
)

// []interface{} 转为真实数据类型列表
func I_list(data interface{}) (total int, d interface{}) {
    switch GetTypeVT(data) {
    case "[]interface{}":
        ds := data.([]interface{})
        total = len(ds)
        if total > 0 {
            //类型不一致，返回原始数据; 一致，返回新的数据
            dtype := GetType(ds[0])
            switch dtype {
            case "int":
                rs := []int{}
                for _, d := range ds {
                    if GetType(d) != dtype { return total, data }
                    rs = append(rs, d.(int))
                }
                return total, rs
            case "string":
                rs := []string{}
                for _, d := range ds {
                    if GetType(d) != dtype { return total, data }
                    rs = append(rs, d.(string))
                }
                return total, rs
            case "float64":
                rs := []float64{}
                for _, d := range ds {
                    if GetType(d) != dtype { return total, data }
                    rs = append(rs, d.(float64))
                }
                return total, rs
            default:
                return total, data
            }
        } else {
            return total, data
        }
    default:
        return total, data
    }
}

//转换其他基本数据类型为字符串，支持数组，接收参数后用 .(string) 或 .([]string) 处理下即可
//支持 string, []string, int, []int, float64, []float64, []interface{}
//format 转换的格式，占位符统一使用 ## ，针对不同的类型会自动替换
func I_str(data interface{}, replaces ...string) string {
    format, replace := "#", "#"
    if len(replaces) == 2 { format, replace = replaces[0], replaces[1] }
    key, tmpFormat := "", ""
    switch GetType(data) {
    case "string":
        tmpFormat = strings.Replace(format, replace, "%s", -1)
        key = fmt.Sprintf(tmpFormat, data)
    case "int", "int64":
        tmpFormat = strings.Replace(format, replace, "%d", -1)
        key = fmt.Sprintf(tmpFormat, data)
    case "float64":
        tmpFormat = strings.Replace(format, replace, "%.f", -1)
        key = fmt.Sprintf(tmpFormat, data)
    default:
    }
    return key
}

//转换其他基本数据类型为字符串，支持数组，接收参数后用 .(string) 或 .([]string) 处理下即可
//支持 string, []string, int, []int, float64, []float64, []interface{}
//format 转换的格式，占位符统一使用 ## ，针对不同的类型会自动替换
func I_strs(data interface{}, replaces ...string) []string {
    total := 0
    dtype := GetType(data)
    if dtype == "[]interface{}" {
        total, data = I_list(data)
        if total > 0 { dtype = GetType(data) }
    }

    format, replace := "#", "#"
    if len(replaces) == 2 { format, replace = replaces[0],  replaces[1] }
    tmpFormat, keys := "", []string{}

    switch GetType(data) {
    case "[]string":
        tmpFormat = strings.Replace(format, replace, "%s", -1)
        for _, v := range data.([]string) {
            keys = append(keys, fmt.Sprintf(tmpFormat, v))
        }
    case "[]int":
        tmpFormat = strings.Replace(format, replace, "%d", -1)
        for _, v := range data.([]int) {
            keys = append(keys, fmt.Sprintf(tmpFormat, v))
        }
    case "[]int64":
        tmpFormat = strings.Replace(format, replace, "%d", -1)
        for _, v := range data.([]int64) {
            keys = append(keys, fmt.Sprintf(tmpFormat, v))
        }
    case "[]float64":
        tmpFormat = strings.Replace(format, replace, "%.f", -1)
        for _, v := range data.([]float64) {
            keys = append(keys, fmt.Sprintf(tmpFormat, v))
        }
    default:    //[]interface{}
    }
    return keys
}

//任意类型转整型
func I_int(data interface{}) (d int) {
    switch GetType(data) {
    case "int":
        return data.(int)
    case "int64":
        d, _ = strconv.Atoi(fmt.Sprintf("%d", data))
    case "string":
        d, _ = strconv.Atoi(data.(string))
    case "float64":
        d, _ = strconv.Atoi(fmt.Sprintf("%.f", data))
    default:
    }
    return d
}

//任意类型转整型数组
func I_ints(data interface{}) (ds []int) {
    total := 0
    dtype := GetType(data)
    if dtype == "[]interface{}" {
        total, data = I_list(data)
        if total > 0 { dtype = GetType(data) }
    }

    ds = []int{}
    switch GetType(data) {
    case "[]int":
        return data.([]int)
    case "[]int64":
        for _, v := range data.([]int64) {
            d, _ := strconv.Atoi(fmt.Sprintf("%d", v))
            ds = append(ds, d)
        }
    case "[]string":
        for _, v := range data.([]string) {
            d, _ := strconv.Atoi(v)
            ds = append(ds, d)
        }
    case "[]float64":
        for _, v := range data.([]float64) {
            d, _ := strconv.Atoi(fmt.Sprintf("%.f", v))
            ds = append(ds, d)
        }
    default:
    }
    return ds
}

