package tconv

import (
    "fmt"
    //"strconv"
)

//从列表中查找子节点
func ListFindChildMap(m map[string]interface{}, idField string, pidField string, pid interface{}) (childs map[string]interface{}) {
    childs = map[string]interface{}{}
    for key, row := range m {
        t := row.(map[string]interface{})
        if t[pidField] == pid {
            t["childs"] =  ListFindChildMap(m, idField, pidField, t[idField])
            childs[key] = t
        }
    }
    return childs
}

//从列表中查找子节点
func ListFindChildList(m map[string]interface{}, idField string, pidField string, pid interface{}) (childs []interface{}) {
    childs = []interface{}{}
    for _, row := range m {
        t := row.(map[string]interface{})
        if t[pidField] == pid {
            t["childs"] =  ListFindChildList(m, idField, pidField, t[idField])
            childs = append(childs, t)
        }
    }
    return childs
}

//列表转为kv结构，key为主键，且转为字符串
func ListToKeyMap(data interface{}, idField string) map[string]interface{} {
    ret := map[string]interface{}{}

    //处理多种类型, switch是一种好的处理数据类型的方式
    items := []map[string]interface{}{}
    switch data.(type) {
    case []map[string]interface{}:
        items = data.([]map[string]interface{})
    case []interface{}:
        items = data.([]map[string]interface{})
    case interface{}:
        items = data.([]map[string]interface{})
    default:
        items = data.([]map[string]interface{})
    }
    for _, row := range items {
        item := row
        idStr := ""
        switch item[idField].(type) {
        case int:
            idStr = Itoa(item[idField].(int))
        case int32:
            idStr = Itoa(int(item[idField].(int32)))
        case int64:
            idStr = Itoa(int(item[idField].(int64)))
        case float64:
            idStr = fmt.Sprintf("%.f", item[idField].(float64))
        case string:
            idStr = item[idField].(string)
        default:
        }
        ret[idStr] = item
    }
    return ret
}

//将列表转为树型结构 interface{}, map[string]interface{}, []interface{}
func ListToTreeMap(data interface{}, idField, pidField string) (ret map[string]map[string]interface{}) {
    m := ListToKeyMap(data, idField)

    ret = map[string]map[string]interface{}{}
    for idStr, row := range m {
        tmp := row.(map[string]interface{})
        isParent := false
        switch tmp[pidField].(type) {
        case float64:
            if tmp[pidField].(float64) == 0 { isParent = true }
        case int:
            if tmp[pidField].(int) == 0 { isParent = true }
        case int32:
            if tmp[pidField].(int32) == 0 { isParent = true }
        case int64:
            if tmp[pidField].(int64) == 0 { isParent = true }
        case string:
            if tmp[pidField].(string) == "0" { isParent = true }
        default:
        }
        if isParent {
            tmp["childs"] = ListFindChildMap(m, idField, pidField, tmp[idField])
            ret[idStr] = tmp
        }
    }
    return ret
}

//将列表转为树型结构 interface{}, map[string]interface{}, []interface{}
func ListToTreeList(data interface{}, idField, pidField string) (ret []map[string]interface{}) {
    m := ListToKeyMap(data, idField)

    ret = []map[string]interface{}{}
    for _, row := range m {
        tmp := row.(map[string]interface{})
        isParent := false
        switch tmp[pidField].(type) {
        case float64:
            if tmp[pidField].(float64) == 0 { isParent = true }
        case int:
            if tmp[pidField].(int) == 0 { isParent = true }
        case int32:
            if tmp[pidField].(int32) == 0 { isParent = true }
        case int64:
            if tmp[pidField].(int64) == 0 { isParent = true }
        case string:
            if tmp[pidField].(string) == "0" { isParent = true }
        default:
        }
        if isParent {
            tmp["childs"] = ListFindChildList(m, idField, pidField, tmp[idField])
            ret = append(ret, tmp)
        }
    }
    return ret
}

