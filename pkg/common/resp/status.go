package resp

type StatusResponse interface {
	Code() int
	Meta() interface{}
}

type statusResponse struct {
	code int
	meta interface{}
}

func (s *statusResponse) Code() int {
	return s.code
}

func (s *statusResponse) Meta() interface{} {
	return s.meta
}

func NewStatusResponse(code int, meta interface{}) StatusResponse {
	return &statusResponse{
		code: code,
		meta: meta,
	}
}

func IsStatusResponse(resp interface{}) (StatusResponse, bool) {
	sr, ok := resp.(StatusResponse)
	if ok {
		return sr, true
	}

	return nil, false
}
