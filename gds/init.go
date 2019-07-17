package gds

import (
)

//公用方法
func GetType(data interface{}) string {
    switch data.(type) {
    case bool:                      return "bool"
    case int:                       return "int"
    case int8:                      return "int8"
    case int16:                     return "int16"
    case int32:                     return "int32"
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
    case byte:                      return "byte"               //alias uint8
    //case rune:                      return "rune"             //alias int32
    case string:                    return "string"
    case interface{}:               return "interface{}"
    //数组
    case []bool:                    return "[]bool"
    case []int:                     return "[]int"
    case []int8:                    return "[]int8"
    case []int16:                   return "[]int16"
    case []int32:                   return "[]int32"
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
    //case []rune:                    return "[]rune"           //alias int32
    case []string:                  return "[]string"
    case []interface{}:             return "[]interface{}"
    //复合结构
    case map[string]string:         return "map[string]string"
    case map[string]interface{}:    return "map[string]interface{}"
    default:                        return ""
    }
    return ""
}
