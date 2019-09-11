package json

import (
//	"fmt"
    "bytes"
    "encoding/json"
	"github.com/json-iterator/go"
	jparse "github.com/buger/jsonparser"
)

var jsonN = jsoniter.ConfigCompatibleWithStandardLibrary

type ValueType = jparse.ValueType

func Marshal(v interface{}) ([]byte, error) {
    //处理html反转义问题
    bf := bytes.NewBuffer([]byte{})
    jsonEncoder := jsonN.NewEncoder(bf)
    jsonEncoder.SetEscapeHTML(false)
    err := jsonEncoder.Encode(v)
    return bf.Bytes(), err

	//return jsonN.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
    raw, err := Marshal(v)
    if err != nil {
        return nil, err
    }
    var buf bytes.Buffer
	err = json.Indent(&buf, raw, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	return jsonN.Unmarshal(data, v)
}

func Encode(v interface{}) ([]byte, error) {
    return Marshal(v)
	//return jsonN.Marshal(v)
}

func EncodeIndent(v interface{}, prefix, indent string) ([]byte, error) {
    return MarshalIndent(v, prefix, indent)
}

func Decode(data []byte, v interface{}) error {
	return jsonN.Unmarshal(data, v)
}

func ArrayEach(data []byte, callback func(value []byte, dataType ValueType, offset int, err error), keys ...string) (offset int, err error) {
	return jparse.ArrayEach(data, callback, keys...)
	//return jparse.ArrayEach(data, cb(value, dataType, offset, err), keys...)
}

func Delete(data []byte, keys ...string) []byte {
	return jparse.Delete(data, keys...)
}

func EachKey(data []byte, callback func(idx int, value []byte, dataType jparse.ValueType, err error), paths ...[]string) int {
	return jparse.EachKey(data, callback, paths...)
}

func Get(data []byte, keys ...string) (value []byte, dataType jparse.ValueType, offset int, err error) {
	return jparse.Get(data, keys...)
}

//封装Get方法，取行json字符串的子串
func GetByte(data []byte, keys ...string) (value []byte, err error) {
	value, _, _, err = jparse.Get(data, keys...)
    return value, err
}

//封装Get方法，取行json字符串的子串
func GetStrs(data []byte, keys ...string) (value []string, err error) {
    rows := []string{}
    callback := func(value []byte, dataType ValueType, offset int, err error) {
        if err == nil { rows = append(rows, string(value)) }
    }
    ArrayEach(data, callback, keys...)
    return rows, nil
}

func GetBoolean(data []byte, keys ...string) (val bool, err error) {
	return jparse.GetBoolean(data, keys...)
}

func GetFloat(data []byte, keys ...string) (val float64, err error) {
	return jparse.GetFloat(data, keys...)
}

func GetInt(data []byte, keys ...string) (val int64, err error) {
	return jparse.GetInt(data, keys...)
}

func GetString(data []byte, keys ...string) (val string, err error) {
	return jparse.GetString(data, keys...)
}

func GetUnsafeString(data []byte, keys ...string) (val string, err error) {
	return jparse.GetUnsafeString(data, keys...)
}

func ObjectEach(data []byte, callback func(key []byte, value []byte, dataType jparse.ValueType, offset int) error, keys ...string) (err error) {
	return jparse.ObjectEach(data, callback, keys...)
}

func Set(data []byte, setValue []byte, keys ...string) (value []byte, err error) {
	return jparse.Set(data, setValue, keys...)
}

//type Ren struct {
//	Name string               `json:"name"`
//	Age int64                 `json:"age"`
//	Scores map[string]float64 `json:"scores"`
//}
//
//func Test_Marshal() {
//	ren := Ren{
//		Name: "hello",
//		Age: 18,
//		Scores: map[string]float64{
//			"shuxue": 98.32,
//			"yuwen": 99.2,
//		},
//	}
//	result, err := Marshal(&ren)
//	fmt.Println(string(result), err)
//}
//
//func Test_Unmarshal() {
//	ren := Ren{}
//	raw := []byte(`{"name":"hello","age":18,"scores":{"shuxue":98.32,"yuwen":99.2}}`)
//	err := Unmarshal(raw, &ren)
//	fmt.Println(ren, err)
//}
//
//func Test_ArrayEach() {
//	raw := []byte(`["name", "hello", "age"]`)
//	rows := []string{}
//	ArrayEach(raw, func(value []byte, dataType ValueType, offset int, err error){
//		if err == nil {
//			rows = append(rows, string(value))
//		}
//	})
//	fmt.Println(rows)
//}

