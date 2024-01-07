package httpx

import (
	"net/http"
	"strings"
)

type Method struct {
}

type MethodGet struct{}

func (MethodGet) Method() string {
	return http.MethodGet
}

type MethodHead struct{}

func (MethodHead) Method() string {
	return http.MethodHead
}

type MethodPost struct{}

func (MethodPost) Method() string {
	return http.MethodPost
}

type MethodPut struct{}

func (MethodPut) Method() string {
	return http.MethodPut
}

type MethodPatch struct{}

func (MethodPatch) Method() string {
	return http.MethodPatch
}

type MethodDelete struct{}

func (MethodDelete) Method() string {
	return http.MethodDelete
}

type MethodConnect struct{}

func (MethodConnect) Method() string {
	return http.MethodConnect
}

type MethodOptions struct{}

func (MethodOptions) Method() string {
	return http.MethodOptions
}

type MethodTrace struct{}

func (MethodTrace) Method() string {
	return http.MethodTrace
}

// MethodMulti return multiple kinds of methods
//
//	use tag `method` to specify the methods
//
// NOTE: rum-gonic will not convert lower-case to upper-case
// example:
//
//	type Index struct {
//		MethodMulti `route:"/index" methods:"GET,POST"`
//	}
type MethodMulti struct {
	Methods string
}

func (m MethodMulti) Method() string {
	return m.Methods
}

func MarshalMethods(ms []string) string {
	return strings.Join(ms, ",")
}
func UnmarshalMethods(str string) []string {
	return strings.Split(str, ",")
}

// MethodAny return all kind of methods
//
//	http.MethodConnect,
//	http.MethodDelete,
//	http.MethodGet,
//	http.MethodHead,
//	http.MethodOptions,
//	http.MethodPost,
//	http.MethodPut,
//	http.MethodPatch,
//	http.MethodTrace,
type MethodAny struct{}

func (m MethodAny) Method() string {
	return MarshalMethods(
		[]string{
			http.MethodConnect,
			http.MethodDelete,
			http.MethodGet,
			http.MethodHead,
			http.MethodOptions,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodTrace,
		},
	)
}
