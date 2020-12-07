package controllers

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qnfnypen/mumori/dao/opmysql"
	"github.com/qnfnypen/mumori/dao/opredis"
	"github.com/qnfnypen/mumori/http/internal"
	"github.com/qnfnypen/mumori/models"
	"github.com/qnfnypen/mumori/pkg/alicloud"
	"github.com/qnfnypen/mumori/public/myerror"
	"github.com/qnfnypen/mumori/utils"
	"github.com/rs/zerolog/log"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

// Register 用户注册
// @Tags 登录
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
	captcha, _ := opredis.GetCaptcha(regParam.Phone)
	if captcha != regParam.Captcha {
		internal.ResBaseInfo(c, myerror.ErrCaptcha)
		return
	}

	// 将用户信息存入MySQL
	// 本打算使用加密算法生成UID的，现在直接使用数据库从2020开始递增
	user := models.User{
		UserName:      regParam.UserName,
		Phone:         regParam.Phone,
		Password:      regParam.Password,
		LastLoginIP:   c.ClientIP(),
		LastLoginTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := opmysql.StoreUserInfo(user)
	if err != nil {
		internal.ResBaseInfo(c, myerror.ErrRegister)
		log.Debug().Str("error", err.Error()).Msg("数据库写入用户信息失败")
		return
	}

	internal.ResBaseInfo(c, myerror.Success)
}

// Login 用户登录 -- 密码登录
// @Tags 登录
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
		err := opmysql.CheckUserInfo(lgParam.Account, lgParam.Password, "phone")
		if err == nil {
			uid := opmysql.GetUserUID("phone", isPhone)
			opmysql.UpdateUserInfo(uid, models.User{LastLoginIP: c.ClientIP(), LastLoginTime: time.Now().Format("2006-01-02 15:04:05")})
			internal.ResBaseInfo(c, myerror.Success)
			return
		}
	}
	isEmail, _ := regexp.MatchString(`\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}$`, account)
	if isEmail {
		err := opmysql.CheckUserInfo(lgParam.Account, lgParam.Password, "email")
		if err == nil {
			uid := opmysql.GetUserUID("email", isEmail)
			opmysql.UpdateUserInfo(uid, models.User{LastLoginIP: c.ClientIP(), LastLoginTime: time.Now().Format("2006-01-02 15:04:05")})
			internal.ResBaseInfo(c, myerror.Success)
			return
		}
	}
	err := opmysql.CheckUserInfo(lgParam.Account, lgParam.Password, "username")
	if err == nil {
		uid := opmysql.GetUserUID("username", account)
		opmysql.UpdateUserInfo(uid, models.User{LastLoginIP: c.ClientIP(), LastLoginTime: time.Now().Format("2006-01-02 15:04:05")})
		internal.ResBaseInfo(c, myerror.Success)
		return
	}

	internal.ResBaseInfo(c, myerror.ErrNameOrPassword)
}

// FastLogin 快捷登录 -- 手机验证码登录
// @Tags 登录
// @Summary 用户快捷登录
// @Description 用户使用手机和验证码快捷登录
// @Accept json
// @Produce json
// @Param fast_login_request body internal.FastloginRequest true "用户快捷登录参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/fast_login [post]
func FastLogin(c *gin.Context) {
	var flg internal.FastloginRequest
	if err := c.ShouldBindJSON(&flg); err != nil {
		internal.ResBaseInfo(c, myerror.ErrParseParam)
		return
	}
	captcha, err := opredis.GetCaptcha(flg.Phone)
	if err != nil {
		internal.ResBaseInfo(c, myerror.ErrPhone)
		return
	}
	if captcha != flg.Captcha {
		internal.ResBaseInfo(c, myerror.ErrCaptcha)
		return
	}

	// 判断该用户是否存在
	err = opmysql.CheckUserInfo(flg.Phone, "", "phone")
	if err != nil {
		username := fmt.Sprintf("user_%s", string([]rune(flg.Phone)[len([]rune(flg.Phone))-4:]))
		err = opmysql.StoreUserInfo(models.User{UserName: username, Phone: flg.Phone})
		if err != nil {
			internal.ResBaseInfo(c, myerror.ErrRegister)
			log.Debug().Str("error", err.Error()).Msg("数据库写入用户信息失败")
			return
		}
	}

	uid := opmysql.GetUserUID("phone", flg.Phone)
	opmysql.UpdateUserInfo(uid, models.User{LastLoginIP: c.ClientIP(), LastLoginTime: time.Now().Format("2006-01-02 15:04:05")})

	internal.ResBaseInfo(c, myerror.Success)
}

// SendCaptcha 发送验证码
// @Tags 登录
// @Summary 发送验证码
// @Description 进行人机验证，发送验证码
// Accept json
// @Produce json
// @Param send_captcha_request body internal.CaptchaRequest true "发送验证码参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/send_captcha [post]
func SendCaptcha(c *gin.Context) {
	var cap internal.CaptchaRequest
	err := c.ShouldBindJSON(&cap)
	if err != nil {
		internal.ResBaseInfo(c, myerror.ErrParseParam)
		return
	}

	b, _ := alicloud.AuthenticateSig(cap.NC.SessionID, cap.NC.Token, cap.NC.Sig, cap.NC.Scene, c.Request.Host)
	if b {
		err := alicloud.SendSMS(cap.Phone, utils.GenerateCode(6))
		if err != nil {
			internal.ResBaseInfo(c, myerror.ErrSendSMS)
			return
		}

		internal.ResBaseInfo(c, myerror.Success)
		return
	}

	internal.ResBaseInfo(c, myerror.ErrAuthenticateSig)
}

// CheckUserName 检测用户名是否已被注册
// @Tags 登录
// @Summary 检测用户名是否已被占用
// @Description 检查已有用户，判断该用户名是否已经被占用
// @Accept json
// @Produce json
// @Param user_name path string true "用户名"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/register_check_username/{user_name} [post]
func CheckUserName(c *gin.Context) {
	name := c.Param("user_name")
	if name == "" {
		internal.ResBaseInfo(c, myerror.ErrUserNameEmpty)
		return
	}
	err := opmysql.CheckUserInfo(name, "", "username")
	if err != nil && err.Error() == "该用户不存在" {
		internal.ResBaseInfo(c, myerror.Success)
		return
	}

	internal.ResBaseInfo(c, myerror.ErrUserNameUsed)
}

// CheckPhone 检测手机号是否已被注册
// @Tags 登录
// @Summary 检测手机号是否已被注册
// @Description 检测注册过的手机号，判断该手机号是否已被注册
// @Accept json
// @Produce json
// @Param phone path string true "手机号"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/register_check_phone/{phone} [post]
func CheckPhone(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		internal.ResBaseInfo(c, myerror.ErrPhoneEmpty)
		return
	}
	err := opmysql.CheckUserInfo(phone, "", "phone")
	if err != nil && err.Error() == "该用户不存在" {
		internal.ResBaseInfo(c, myerror.Success)
		return
	}

	internal.ResBaseInfo(c, myerror.ErrPhoneUsed)
}

// ForgetPassword 忘记密码
// @Tags 登录
// @Summary 忘记密码
// @Description 用户修改登录密码
// @Accept json
// @Produce json
// @Param forget_password body internal.ForgetPasswordRequest true "忘记密码参数"
// @Success 200 {object} internal.ResponseBase  "{code:200,message:"成功"}"
// @Failure 400 {ogject} internal.ResponseBase  "{code:400,message:"失败"}"
// @Router /auth/forget_password [post]
func ForgetPassword(c *gin.Context) {
	var fp internal.ForgetPasswordRequest
	err := c.ShouldBindJSON(&fp)
	if err != nil {
		internal.ResBaseInfo(c, myerror.ErrParseParam)
		return
	}

	// 判断两次密码是否相同
	if fp.NewPassword != fp.PasswordConfirm {
		internal.ResBaseInfo(c, myerror.ErrTwoPassword)
		return
	}
	// 判断验证码是否正确
	captcha, _ := opredis.GetCaptcha(fp.Phone)
	if fp.Captcha != captcha {
		internal.ResBaseInfo(c, myerror.ErrCaptcha)
		return
	}

	uid := opmysql.GetUserUID("phone", fp.Phone)
	err = opmysql.UpdateUserInfo(uid, models.User{Password: fp.NewPassword})
	if err != nil {
		internal.ResBaseInfo(c, myerror.ErrUpdatePassword)
		return
	}

	internal.ResBaseInfo(c, myerror.Success)
}
