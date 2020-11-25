package internal

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
	"github.com/qnfnypen/mumori/public/myerror"
)

// ResponseBase 基础响应
type ResponseBase struct {
	Code    myerror.ErrCode `json:"code"`
	Message string          `json:"message"`
}

// resJSON 返回响应
func resJSON(c *gin.Context, code int, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		log.Debug().Str("error", err.Error()).Msg("响应体json解析失败")
	}
	c.Data(code, "application/json;charset=utf-8", buf)
	c.Abort()
}

// ResBaseInfo 返回基础信息
func ResBaseInfo(c *gin.Context, errcode myerror.ErrCode) {
	statusCode := 200
	if errcode != myerror.Success {
		statusCode = 400
	}
	base := ResponseBase{
		Code:    errcode,
		Message: myerror.GetMsg(errcode),
	}

	resJSON(c, statusCode, base)
}
