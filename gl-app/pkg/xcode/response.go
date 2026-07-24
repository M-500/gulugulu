package xcode

import "context"

type JsonResult struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func ErrHandler(err error) (int, any) {

	return 200, JsonResult{
		Code: 500101,
		Msg:  err.Error(),
	}
}

type Response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "请求成功",
		Data:    data,
	}
}

func OkHandler(_ context.Context, v interface{}) any {
	return Success(v)
}
