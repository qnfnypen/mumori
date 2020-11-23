package utils

import (
	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
	"gopkg.in/gomail.v2"

	"github.com/spf13/viper"
)

type mailClient struct {
	username string
	password string
	host     string
	port     int
}

// SendMail 发送邮件
func SendMail(to []string, subject, body string) error {
	cli := mailClient{
		username: viper.GetString("Alarm.Email.User"),
		password: viper.GetString("Alarm.Email.Password"),
		host:     viper.GetString("Alarm.Email.Host"),
		port:     viper.GetString("Alarm.Email.Port"),
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(cli.username, "木森云"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(cli.host, cli.port, cli.username, cli.password)
	err := d.DialAndSend(m)

	return err
}
