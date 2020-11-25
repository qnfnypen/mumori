package opredis

import (
	"context"
	"time"
)

// SetCaptchaFor5M 存储验证码并设置有效时长为5分钟
func SetCaptchaFor5M(phone, captcha string) error {
	_, err := rdb.SetNX(context.TODO(), phone, captcha, 5*time.Minute).Result()
	return err
}

// GetCaptcha 获取验证码
func GetCaptcha(phone string) (string, error) {
	return rdb.Get(context.TODO(), phone).Result()
}
