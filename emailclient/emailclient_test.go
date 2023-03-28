package emailclient

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"net/smtp"
	"testing"
)

const (
	SMTPHost     = "smtp.gmail.com"
	SMTPPort     = ":587"
	SMTPUsername = "xxx@gmail.com"
	SMTPPassword = "xxxx"
)

func sendEmail(receiver string) {
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)
	msg := []byte("Subject: 这里是标题内容\r\n\r\n" + "这里是正文内容\r\n")
	err := smtp.SendMail(SMTPHost+SMTPPort, auth, SMTPUsername, []string{receiver}, msg)
	if err != nil {
		log.Fatal("failed to send email:", err)
	}
}

func sendHTMLEmail(receiver string, html []byte) {
	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)
	msg := append([]byte("Subject: 这里是标题内容\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"),
		html...)
	err := smtp.SendMail(SMTPHost+SMTPPort, auth, SMTPUsername, []string{receiver}, msg)
	if err != nil {
		log.Fatal("failed to send email:", err)
	}
}

func Test_send(t *testing.T) {
	sendHTMLEmail("接受者@gmail.com", []byte("<html><h2>这是网页内容</h2></html>"))
}

func Test_gomail(t *testing.T) {
	m := gomail.NewMessage()
	m.SetHeader("From", "dcs@dongs.top")
	// To、Cc（carbon copy）和Bcc（blind carbon copy），分别是收件人、抄送、密送
	m.SetHeader("To", "dongs@dongs.top")
	m.SetAddressHeader("Bcc", "dongs@dongs.top", "Dong")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Baby</b>!")
	//m.SetBody("text/plain", "Hello!")
	m.Attach("F:\\OneDrive - abc\\图片\\美女壁纸\\0d9d6dfc332c34d7bbde92072ab81f8e.jpg")

	d := gomail.NewDialer("smtp.qiye.aliyun.com", 465, "dcs@dongs.top", "1qaz@WSX")
	//d := gomail.Dialer{Host: "localhost", Port: 587}

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func Test_template(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	emailTemplate.ExecuteTemplate(b, "email_body_template_name_3", "dcs")
	fmt.Print(b.String())
}

func Test_sendTemplate(t *testing.T) {
	emailConfig := EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		Username: "xxx@gmail.com",
		Password: "xxxx",
	}
	New("", emailConfig)
	err := GetClient().SendTemplateEmail("dongs@dongs.top", "一封邮件来了", "email_body_template_name_3", nil)
	if err != nil {
		t.Log(err)
	}
}

func Test_LoadTemplate(t *testing.T) {
	LoadTemplate("templates/email/*.gohtml")
}
