package alicloud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
)

// 阿里云客户端
type aliClient struct {
	// 人机验证客户端
	sigClt *sdk.Client
	// 短信服务客户端
	smsClt *dysmsapi.Client
}

type resp struct {
	Code int `json:"Code"`
}

var a aliClient

func init() {
	var err error
	a.sigClt, err = sdk.NewClientWithAccessKey(
		viper.GetString("Aliyun.SMS.RegionID"),
		viper.GetString("Aliyun.SMS.AccessKeyID"),
		viper.GetString("Aliyun.SMS.AccessKeySecret"),
	)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("阿里云人机验证客户端创建失败")
	}
	a.smsClt, err = dysmsapi.NewClientWithAccessKey(
		viper.GetString("Aliyun.SMS.RegionID"),
		viper.GetString("Aliyun.SMS.AccessKeyID"),
		viper.GetString("Aliyun.SMS.AccessKeySecret"),
	)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("阿里云短信客户端创建失败")
	}
}

// AuthenticateSig 滑块验证
func AuthenticateSig(sessionID,token,sig,scene,remoteIP string) (bool, error) {
	request := requests.NewCommonRequest()
	request.Method = http.MethodPost
	request.Scheme = "https"
	request.Domain = "afs.aliyuncs.com"
	request.Version = "2018-01-12"
	request.ApiName = "AuthenticateSig"
	// 会话ID。必填参数，从前端获取，不可更改
	request.QueryParams["SessionId"] = sessionID
	// 签名串。必填参数，从前端获取，不可更改
	request.QueryParams["Sig"] = sig
	// 请求唯一标识。必填参数，从前端获取，不可更改
	request.QueryParams["Token"] = token
	// 场景标识。必填参数，从前端获取，不可更改
	request.QueryParams["Scene"] = scene
	// 应用类型标识。必填参数，后端填写
	request.QueryParams["AppKey"] = viper.GetString("Aliyun.SMS.AppKey")
	// 客户端IP。必填参数，后端填写
	request.QueryParams["RemoteIp"] = remoteIP

	response,err := a.sigClt.ProcessCommonRequest(request)
	if err != nil {
		return false,err
	}

	log.Info().Msgf("aliyun send sms response: %v",response)
	
	respbyte := response.GetHttpContentBytes()
	var sigResp resp
	err = json.Unmarshal(respbyte,&sigResp)
	if err != nil {
		return false,err
	}
	if sigResp.Code == 100 {
		return true,nil
	}

	return false,nil
}

// SendSMS 发送验证码
func SendSMS(mobile,templateParam string) error {
	request := requests.NewCommonRequest()
	request.Method = http.MethodPost
	request.Scheme = "https"
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = viper.GetString("Aliyun.SMS.RegionID")
	// 接收短信的手机号码，必填
	request.QueryParams["PhoneNumbers"] = mobile
	// 短信签名名称，必填
	request.QueryParams["SignName"] = viper.GetString("Aliyun.SMS.SignName")
	// 短信模板ID，必填
	request.QueryParams["TemplateCode"] = viper.GetString("Aliyun.SMS.TemplateCode")
	// 短信模板对应的实际值：{"code":"1111"}
	request.QueryParams["TemplateParam"] = fmt.Sprintf(`{"code":"%s"}`,templateParam)

	response,err := a.smsClt.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	log.Info().Msgf("aliyun send sms response: %v",response)

	return nil
}
