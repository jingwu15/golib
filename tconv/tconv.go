package tconv

import (
    "fmt"
	"bytes"
	"strconv"
    "reflect"
    "strings"
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

func KeysMapStrInf(data map[string]interface{}) []string {
    keys := []string{}
    for key, _ := range data {
        keys = append(keys, key)
    }
    return keys
}

func KeysMapStrInt(data map[string]int) []string {
    keys := []string{}
    for key, _ := range data {
        keys = append(keys, key)
    }
    return keys
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
//result := tconv.MergeMap(&a, &b)
//fmt.Println(result)
func MergeMap(raw interface{}, ext interface{}) map[string]interface{} {
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
            keysRaw := KeysMapStrInf(rawData)
            keysExt := KeysMapStrInf(extData)
            keys := keysRaw
            for _,key := range keysExt {
                keys = append(keys, key)
            }
            keysMap := KeyStrToMap(keys)
            keys = KeysMapStrInt(keysMap)
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
                    rawData[k] = MergeMap(rawData[k], extData[k])
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

