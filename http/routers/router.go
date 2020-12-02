package routers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/qnfnypen/mumori/http/controllers"
	"github.com/qnfnypen/mumori/http/middleware"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

// GenerateEngine 生成路由引擎
func GenerateEngine() *gin.Engine {
	// 设置Gin运行模式
	gin.SetMode(viper.GetString("Gin.RunMode"))

	// 将gin日志输出到文件
	var logFile *os.File
	if runtime.GOOS == "windows" {
		path := viper.GetString("Gin.LogFile.Windows")
		dir := filepath.Dir(path)
		pwd, _ := os.Getwd()
		filedir := fmt.Sprintf("%s/%s", pwd, dir)
		os.MkdirAll(filedir, os.ModePerm)
		file := fmt.Sprintf("%s/%s", pwd, path)
		logFile, _ = os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	} else {
		path := viper.GetString("Gin.LogFile.Linux")
		dir := filepath.Dir(path)
		os.MkdirAll(dir, os.ModePerm)
		logFile, _ = os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	}
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(logFile)

	engine := gin.Default()

	// 是否启用跨域
	if viper.GetBool("Gin.CORS.Enable") {
		engine.Use(middleware.CORSMiddleware())
	}

	auth := engine.Group("/auth")
	{
		auth.POST("/send_captcha", controllers.SendCaptcha)
		auth.POST("/register_check_username",controllers.CheckUserName)
		auth.POST("/register_check_phone",controllers.CheckPhone)
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/fast_login", controllers.FastLogin)
		auth.POST("/forget_password",controllers.ForgetPassword)
	}

	return engine
}
