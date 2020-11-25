package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qnfnypen/mumori/dao/opredis"
	"github.com/qnfnypen/mumori/http/internal"
	"github.com/qnfnypen/mumori/public/myerror"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

// Register 用户注册
// @Summary 用户注册
// @Description 用户根据手机号进行注册
// @Accept json
// @Produce json
// @Param register_request body internal.RegisterRequest true "用户注册参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var regParam internal.RegisterRequest
	if err := c.ShouldBindJSON(&regParam); err != nil {
		internal.ResBaseInfo(c, myerror.ErrParseParam)
		return
	}
	// 判断两次输入的密码是否相同
	if regParam.Password != regParam.PasswordConfirm {
		internal.ResBaseInfo(c, myerror.ErrTwoPassword)
		return
	}
	// 判断验证码是否正确
	captcha, err := opredis.GetCaptcha(regParam.Phone)
	if err != nil || captcha != regParam.Captcha {
		internal.ResBaseInfo(c, myerror.ErrCaptcha)
		return
	}

	internal.ResBaseInfo(c, myerror.Success)
}

// Login 用户登录 -- 密码登录
// @Summary 用户登录
// @Description 用户使用昵称/邮箱/手机，并使用密码进行登录
// @Accept json
// @Produce json
// @Param login_request body internal.LoginRequest true "用户登录参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/login [post]
func Login(c *gin.Context) {

}

// FastLogin 快捷登录 -- 手机验证码登录
// @Summary 用户快捷登录
// @Description 用户使用手机和验证码快捷登录
// @Accept json
// @Produce json
// @Param fast_login_request body internal.FastloginRequest "用户快捷登录参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/fast_login [post]
func FastLogin(c *gin.Context) {

}

// SendCaptcha 发送验证码
func SendCaptcha(c *gin.Context) {

}
