package controllers

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/qnfnypen/mumori/dao/opmysql"
	"github.com/qnfnypen/mumori/dao/opredis"
	"github.com/qnfnypen/mumori/http/internal"
	"github.com/qnfnypen/mumori/models"
	"github.com/qnfnypen/mumori/public/myerror"
	"github.com/rs/zerolog/log"

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

	// 将用户信息存入MySQL
	// 本打算使用加密算法生成UID的，现在直接使用数据库从2020开始递增
	user := models.User{
		UserName: regParam.UserName,
		Phone: regParam.Phone,
		Password: regParam.Password,
	}
	err = opmysql.StoreUserInfo(user)
	if err != nil {
		internal.ResBaseInfo(c,myerror.ErrRegister)
		log.Debug().Str("error",err.Error()).Msg("数据库写入用户信息失败")
		return
	}

	internal.ResBaseInfo(c, myerror.Success)
}

// Login 用户登录 -- 密码登录
// @Summary 用户登录
// @Description 用户使用用户名/邮箱/手机，并使用密码进行登录
// @Accept json
// @Produce json
// @Param login_request body internal.LoginRequest true "用户登录参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var lgParam internal.LoginRequest
	if err := c.ShouldBindJSON(&lgParam); err != nil {
		internal.ResBaseInfo(c, myerror.ErrParseParam)
		return
	}
	// 判断输入的是手机、邮箱或用户名
	// 手机：^1\d{10}$
	// 邮箱：\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}$
	account := lgParam.Account
	isPhone, _ := regexp.MatchString(`^1\d{10}$`, account)
	if isPhone {

	}

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
// @Summary 发送验证码
// @Description 进行人机验证，发送验证码
// Accept json
// @Produce json
// @Param send_captcha_request body internal.CaptchaRequest "发送验证码参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/send_captcha [post]
func SendCaptcha(c *gin.Context) {

}

// CheckUserName 检测用户名是否已被注册
func CheckUserName(c *gin.Context) {

}

// CheckPhone 检测手机号是否已被注册
func CheckPhone(c *gin.Context) {

}
