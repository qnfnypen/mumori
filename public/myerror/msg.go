package myerror

var errmsg = map[ErrCode]string{
	Success:"成功",
	Fail:"失败",
	ErrParseParam:"请求参数错误",
	ErrTwoPassword:"两次输入的密码不一致",
	ErrCaptcha:"验证码错误",
}

// GetMsg 根据错误码获取错误信息
func GetMsg(code ErrCode) string {
	msg,ok := errmsg[code]
	if ok {
		return msg
	}

	return errmsg[Fail]
}