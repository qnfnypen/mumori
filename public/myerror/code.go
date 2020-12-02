package myerror

// ErrCode 自定义错误类型
type ErrCode int

const (
	// ErrParseParam 解析请求参数错误
	ErrParseParam ErrCode = 100 + iota
	// ErrTwoPassword 两次输入的密码不一致
	ErrTwoPassword
	// ErrCaptcha 验证码错误
	ErrCaptcha
	// ErrRegister 服务器错误，用户注册失败
	ErrRegister
	// ErrNameOrPassword 用户名或密码错误
	ErrNameOrPassword
	// ErrPhone 手机号错误
	ErrPhone
	// ErrUserNameEmpty 用户名不能为空
	ErrUserNameEmpty
	// ErrPhoneEmpty 手机号不能为空
	ErrPhoneEmpty
	// ErrUserNameUsed 用户名已被使用
	ErrUserNameUsed
	// ErrPhoneUsed 手机号已被使用
	ErrPhoneUsed
	// ErrSendSMS 短信验证码发送失败
	ErrSendSMS
	// ErrAuthenticateSig 滑块验证失败
	ErrAuthenticateSig
	// ErrUpdatePassword 修改用户密码失败
	ErrUpdatePassword

	// Success 成功
	Success ErrCode = 200
	// Fail 失败，错误超过一个时返回
	Fail ErrCode = 400
)