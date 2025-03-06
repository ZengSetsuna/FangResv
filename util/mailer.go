package util

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type Mail interface {
	SendEmail(to, subject, body string) error
}

type Mailer struct {
	SmtpHost string
	SmtpPort int
	SmtpUser string
	SmtpPass string
}

// SendEmail 发送邮件
func (mailer *Mailer) SendEmail(to, subject, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", mailer.SmtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body) // 允许 HTML 格式的邮件内容

	d := gomail.NewDialer(mailer.SmtpHost, mailer.SmtpPort, mailer.SmtpUser, mailer.SmtpPass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		log.Printf("邮件发送失败: %v", err)
		return fmt.Errorf("邮件发送失败: %w", err)
	}
	log.Println("邮件发送成功")
	return nil
}
