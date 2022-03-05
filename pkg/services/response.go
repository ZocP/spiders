package services

type Response struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"status"`
	Data     interface{} `json:"Data"`
}

func ErrorResponse(err error) *Response {
	return &Response{Code: -1, ErrorMsg: err.Error()}
}

func SuccessResponse(data interface{}) *Response {
	return &Response{Code: 0, ErrorMsg: "ok", Data: data}
}
