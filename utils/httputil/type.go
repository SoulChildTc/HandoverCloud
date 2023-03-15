package httputil

type ResponseBody struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}
