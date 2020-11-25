package internal

// RegisterRequest 注册参数
type RegisterRequest struct {
	UserName        string `json:"user_name" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
	Captcha         string `json:"captcha" binding:"required"`
}

// CaptchaRequest 发送验证码参数
type CaptchaRequest struct {
	Phone string `json:"phone" binding:"required"`
	NC struct {
		Token     string `json:"token" binding:"required"`     // 滑块验证码token
		Sig       string `json:"sig" binding:"required"`       // 滑块验证码sig
		Scene     string `json:"scene" binding:"required"`     // 滑块验证码scene
		SessionID string `json:"sessionId" binding:"required"` // 滑块验证码session_id
	} `json:"nc" binding:"required"`
}

// LoginRequest 用户登录参数
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// FastloginRequest 用户快捷登录参数
type FastloginRequest struct {
	Phone   string `json:"phone" binding:"required"`
	Captcha string `json:"captcha" binding:"required"`
}
