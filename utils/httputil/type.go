package httputil

type ResponseBody struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
}

type PageResponseBody struct {
	Status string    `json:"status"`
	Msg    string    `json:"msg"`
	Data   *PageResp `json:"data"`
}

type PageResp struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
	Total int `json:"total"`
	Items any `json:"items"`
}
