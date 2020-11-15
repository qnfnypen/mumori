package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

// CORSMiddleware 跨域请求中间件
func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: viper.GetStringSlice("Gin.CORS.AllowOrigins"),
		AllowMethods: viper.GetStringSlice("Gin.CORS.AllowMethods"),
		AllowHeaders: viper.GetStringSlice("Gin.CORS.AllowHeaders"),
		AllowCredentials: viper.GetBool("Gin.CORS.AllowCredentials"),
		MaxAge: viper.GetDuration("Gin.CORS.MaxAge") * time.Second,
	})
}