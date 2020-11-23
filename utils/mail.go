package utils

import (
	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type mailClient struct {
	user     string
	password string
	server   string
	port     int
}

// SendMail 发送邮件
func SendMail(mailto []string, subject, body string) error {
	client := mailClient{
		user:     viper.GetString("Alarm.Email.User"),
		password: viper.GetString("Alarm.Email.Password"),
		server:   viper.GetString("Alarm.Email.Server"),
		port:     viper.GetInt("Alarm.Email.Port"),
	}

	m := gomail.NewMessage()
	// 指定发件人
	// m.SetHeader("From", viper.GetString("Alarm.Email.User"))
	m.SetHeader("To", mailto...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// 添加附件
	// m.Attach("log.txt")

	d := gomail.NewDialer(client.server, client.port, client.user, client.password)
	err := d.DialAndSend(m)
	return err
}
