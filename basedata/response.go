package basedata

const ServerUnavailable = "server unavailable"
const NotLogIn = "You may not login"
const InvalidAmount = "invalid amount money"
const UserBusy = "your transaction is busy"
const Success = "success"
const Fail = "fail"

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
