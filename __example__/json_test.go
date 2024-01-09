package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func Test_JsonUmarshal(t *testing.T) {
	jsonStr := `{
		"name": "jarvis",
		"age": 18,
		"address": {
			"home": "home",
			"school": "school"
		}
	}`

	JsonUmmarshal(jsonStr)
	JsonDecoder(jsonStr)
}

func Test_JsonMarshal2(t *testing.T) {
	jsonstr := `{
    "name": "jarvis",
    "age": 18,
    "home": "home",
    "school": "school"
}`

	JsonUmmarshal(jsonstr)
	JsonDecoder(jsonstr)
}

func JsonUmmarshal(jsonStr string) {
	var p *MyPerston
	fmt.Printf("JsonUmmarshal First: %+v\n", p)
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JsonUmmarshal Second: %+v\n", p)
}

func JsonDecoder(jsonStr string) {
	var p *MyPerston
	fmt.Printf("JsonDecoder Frist: %+v\n", p)
	buf := bytes.NewBuffer([]byte(jsonStr))
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JsonDecoder Second: %+v\n", p)
}

type MyPerston struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address `json:",inline"`
}

type MyAddress struct {
	Home   string `json:"home"`
	School string `json:"school"`
}

func TestMarshal(t *testing.T) {
	my := &MyPerston{
		Name: "jarvis",
		Age:  18,
		Address: Address{
			Home:   "home",
			School: "school",
		},
	}

	b, _ := json.Marshal(my)
	fmt.Printf("%s\n", b)
}
