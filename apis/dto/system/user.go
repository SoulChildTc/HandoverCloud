package system

type Register struct {
	Username string  `json:"username" binding:"required" msg:"用户名必填"`
	Mobile   string  `json:"mobile" binding:"required,mobile" msg:"手机号必填" mobile_err:"手机号格式有误"`
	Password string  `json:"password" binding:"required" msg:"密码必填"`
	Email    *string `json:"email" binding:"omitempty,email" email_err:"邮箱格式错误"`
}

type Login struct {
	Mobile   string `json:"mobile" binding:"required" required_err:"手机号不能为空"`
	Password string `json:"password" binding:"required" required_err:"密码不能为空"`
}
