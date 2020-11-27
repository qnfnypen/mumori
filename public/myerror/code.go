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

	// Success 成功
	Success ErrCode = 200
	// Fail 失败，错误超过一个时返回
	Fail ErrCode = 400
)