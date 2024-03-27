package util

import (
	"github.com/mufe/golang-base/camp/xlog"
	"net/smtp"
)

func SendEmail(smtpServer, smtpPort, username, password, from, subject, body string, to []string) error {

	// 认证信息
	auth := smtp.PlainAuth("", username, password, smtpServer)

	// 构建电子邮件内容
	message := []byte(subject + body)

	// 使用发送邮件的邮箱地址、接收者列表、电子邮件内容，通过auth完成认证，并发送电子邮件
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}
