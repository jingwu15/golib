package json

import (
	"fmt"
//	"github.com/json-iterator/go"
//	jparse "github.com/buger/jsonparser"
)

type Ren struct {
	Name string               `json:"name"`
	Age int64                 `json:"age"`
	Scores map[string]float64 `json:"scores"`
}

func Test_Marshal() {
	ren := Ren{
		Name: "hello",
		Age: 18,
		Scores: map[string]float64{
			"shuxue": 98.32,
			"yuwen": 99.2,
		},
	}
	result, err := Marshal(&ren)
	fmt.Println(string(result), err)
}

func Test_Unmarshal() {
	ren := Ren{}
	raw := []byte(`{"name":"hello","age":18,"scores":{"shuxue":98.32,"yuwen":99.2}}`)
	err := Unmarshal(raw, &ren)
	fmt.Println(ren, err)
}

func Test_ArrayEach() {
	raw := []byte(`["name", "hello", "age"]`)
	rows := []string{}
	ArrayEach(raw, func(value []byte, dataType ValueType, offset int, err error){
		if err == nil {
			rows = append(rows, string(value))
		}
	})
	fmt.Println(rows)
}

