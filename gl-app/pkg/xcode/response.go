package xcode

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
