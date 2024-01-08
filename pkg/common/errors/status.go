package errors

import (
	"errors"
	"reflect"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type StatusError struct {
	// standard error
	Err error `json:"error"`

	// error message, it will be set by http status if it is empty
	Message string `json:"message"`

	// meta data, anything
	Meta any `json:"meta"`
}

func New(err error, meta any) *StatusError {
	return &StatusError{
		Err:  err,
		Meta: meta,
	}
}

func (e *StatusError) Error() string {
	return e.Err.Error()
}

// JSON creates a properly formatted JSON
// github.com/cloudwego/hertz/pkg/common/errors
func (e *StatusError) JSON() interface{} {
	jsonData := make(map[string]interface{})

	jsonData["error"] = e.Error()
	jsonData["message"] = e.Message

	if e.Meta != nil {
		value := reflect.ValueOf(e.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return e.Meta
		// case reflect.Map:
		// 	for _, key := range value.MapKeys() {
		// 		jsonData[key.String()] = value.MapIndex(key).Interface()
		// 	}
		default:
			jsonData["meta"] = e.Meta
		}
	}

	return jsonData
}

func (e *StatusError) Unwrap() error {
	return e.Err
}

func (e *StatusError) SetMessage(msg string) *StatusError {
	e.Message = msg
	return e
}

func (e *StatusError) SetCode(code int) {
	if e.Message == "" {
		e.SetMessage(consts.StatusMessage(code))
	}
}

func AsStatusError(err error) (*StatusError, bool) {

	se := &StatusError{}
	ok := errors.As(err, &se)

	if ok {
		return se, true
	}

	return nil, false
}
