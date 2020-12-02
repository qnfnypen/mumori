package myerror

var errmsg = map[ErrCode]string{
	Success:"成功",
	Fail:"失败",
	ErrParseParam:"请求参数错误",
	ErrTwoPassword:"两次输入的密码不一致",
	ErrCaptcha:"验证码错误",
	ErrRegister:"服务器错误，用户注册失败",
	ErrNameOrPassword:"用户名或密码错误",
	ErrPhone:"手机号错误，请检查手机号后重新获取验证码",
	ErrUserNameEmpty:"用户名不能为空",
	ErrPhoneEmpty:"手机号不能为空",
	ErrUserNameUsed:"该用户名已被使用",
	ErrPhoneUsed:"该手机号已被注册，请转到登录页面",
	ErrSendSMS:"验证码发送失败",
	ErrAuthenticateSig:"滑块验证码失败",
	ErrUpdatePassword:"修改密码失败",
}

// GetMsg 根据错误码获取错误信息
func GetMsg(code ErrCode) string {
	msg,ok := errmsg[code]
	if ok {
		return msg
	}

	return errmsg[Fail]
}