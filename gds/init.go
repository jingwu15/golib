package gds

import (
    //"fmt"
)

//变量类型有三种方式：
//1. reflect.TypeOf 最全，但是性能最差
//2. switch-type    最快，但对动态类型只会返回 interface{}， 性能是 reflect 的 20倍
//3. value.(type)   较平衡，性能与if的多少及位置有关，但总体是 reflect 的 10 倍
//使用 switch-type 方式，检测变量的基本类型
func GetType(data interface{}) string {
    //case的前后位置对性能影响小
    switch data.(type) {
    case bool:                      return "bool"
    case int:                       return "int"
    case int8:                      return "int8"
    case int16:                     return "int16"
    //case int32:                     return "int32"
    case int64:                     return "int64"
    case uint:                      return "uint"
    //case uint8:                     return "uint8"
    case uint16:                    return "uint16"
    case uint32:                    return "uint32"
    case uint64:                    return "uint64"
    case uintptr:                   return "uintptr"
    case float32:                   return "float32"
    case float64:                   return "float64"
    case complex64:                 return "complex64"
    case complex128:                return "complex128"
    case byte:                      return "byte"             //alias uint8
    case rune:                      return "rune"             //alias int32
    case string:                    return "string"
    //数组
    case []bool:                    return "[]bool"
    case []int:                     return "[]int"
    case []int8:                    return "[]int8"
    case []int16:                   return "[]int16"
    //case []int32:                   return "[]int32"
    case []int64:                   return "[]int64"
    case []uint:                    return "[]uint"
    //case []uint8:                   return "[]uint8"
    case []uint16:                  return "[]uint16"
    case []uint32:                  return "[]uint32"
    case []uint64:                  return "[]uint64"
    case []uintptr:                 return "[]uintptr"
    case []float32:                 return "[]float32"
    case []float64:                 return "[]float64"
    case []complex64:               return "[]complex64"
    case []complex128:              return "[]complex128"
    case []byte:                    return "[]byte"           //alias uint8
    case []rune:                    return "[]rune"           //alias int32
    case []string:                  return "[]string"
    case []interface{}:             return "[]interface{}"
    //复合结构
    case map[string]string:         return "map[string]string"
    case map[string]interface{}:    return "map[string]interface{}"
    case []map[string]interface{}:  return "[]map[string]interface{}"
    default:          // interface{}
        return GetTypeVT(data)    //switch 不能检测试 动态的interface{}
    }
    return ""
}

//使用 value.(type) 方式检测变量类型，if更多更靠后，性能更差
func GetTypeVT(data interface{}) string {
    if _, ok := data.([]string);                  ok { return "[]string"                  }
    if _, ok := data.([]int);                     ok { return "[]int"                     }
    if _, ok := data.([]interface{});             ok { return "[]interface{}"             }
    if _, ok := data.([]map[string]interface{});  ok { return "[]map[string]interface{}"    }
    if _, ok := data.(map[string]interface{});    ok { return "map[string]interface{}"    }
    if _, ok := data.(map[string]string);         ok { return "map[string]string"         }
    if _, ok := data.(map[string]int);            ok { return "map[string]int"            }
    if _, ok := data.(map[string]float64);        ok { return "map[string]float64"        }
    if _, ok := data.(string);                    ok { return "string"                    }
    if _, ok := data.(int);                       ok { return "int"                       }
    if _, ok := data.(map[int]interface{});       ok { return "map[int]interface{}"       }
    if _, ok := data.(map[int]string);            ok { return "map[int]string"            }
    if _, ok := data.(func(interface{}, interface{})(error));          ok { return "func(interface{}, interface{})(error)"            }
    return "othervt"
}

//使用 value.(type) 方式检测变量类型，if更多更靠后，性能更差, 仅针对函数
func GetTypeFun(data interface{}) string {
    if _, ok := data.(func(interface{})(error)); ok { return "func(interface{})(error)" }
    if _, ok := data.(func(interface{})(int, error)); ok { return "func(interface{})(int, error)" }
    if _, ok := data.(func(interface{})([]int, error)); ok { return "func(interface{})([]int, error)" }
    if _, ok := data.(func(interface{}, interface{})(error)); ok { return "func(interface{}, interface{})(error)" }
    return ""
}

