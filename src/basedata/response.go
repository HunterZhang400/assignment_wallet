package basedata

const (
	ServerUnavailable = "server unavailable"
	NotLogIn          = "You may not login"
	InvalidAmount     = "invalid amount money"
	UserBusy          = "your transaction is busy"
	Success           = "success"
	Fail              = "fail"
)

func NewErrorResponse(msg string) Response {
	return Response{Msg: msg}
}
func NewResponse(data any) Response {
	return Response{Data: data}
}

type Response struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
