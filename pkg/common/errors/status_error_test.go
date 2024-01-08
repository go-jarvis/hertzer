package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

var Err3 = errors.New("Origin Error")

func TestErrir(t *testing.T) {
	// err := errors.New("Origin Error")

	err1 := New(Err3, nil)
	err1 = err1.SetMessage("Error Message")

	// t.Log(err1)

	err2 := &StatusError{}
	if !errors.As(err1, &err2) {
		t.Error("err1 is not StatusError")
	}

	fmt.Println("err1 is StatusError")
	fmt.Println("Error2: ", err2.Message)

	if !errors.Is(err1, Err3) {
		t.Error("err1 is not Err3")
	}

	fmt.Println("err1 is Err3")
	fmt.Println("Error2: ", err2.Message)
}

func TestJSON2(t *testing.T) {

	data := map[string]interface{}{
		"key1":  "value1",
		"key2":  "value2",
		"error": "map error",
	}
	se := New(Err3, data)
	se.SetMessage("Error Message")
	outputJson(se)

	outputJson(se.JSON())
}

func outputJson(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))
}
