// 发送电子邮件相关操作
// https://github.com/go-gomail/gomail

package emailclient

import (
	"bytes"
	"errors"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

var clients = make(map[string]*emailClient)

var emailTemplate *template.Template

func New(name string, config EmailConfig) *emailClient {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	client := &emailClient{dialer: dialer}
	clients[name] = client
	return client
}

// LoadTemplate 加载邮件模板
func LoadTemplate(path string) {
	tpl, err := template.ParseGlob(path)
	if err != nil {
		log.Println(err)
	}
	emailTemplate = tpl
}

func GetClientByName(name string) *emailClient {
	return clients[name]
}

func GetClient() *emailClient {
	return clients[""]
}

type emailClient struct {
	dialer *gomail.Dialer
}

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Email struct {
	From    string   // 发件人
	To      []string // 收件人
	Cc      []string // 抄送
	Bcc     []string // 秘密抄送
	Subject string   // 标题
	Body    string   // 内容
	Attach  []string // 附件地址
}

func (client *emailClient) SendEmail(email *Email) error {
	m := gomail.NewMessage()
	if email.From == "" {
		return errors.New("发件人不能为空！")
	} else {
		m.SetHeader("From", email.From)
	}
	// To、Cc（carbon copy）和Bcc（blind carbon copy），分别是收件人、抄送、密送
	if email.To == nil || len(email.To) == 0 {
		return errors.New("收件人不能为空！")
	} else {
		m.SetHeader("To", email.To...)
	}
	if email.Cc != nil && len(email.Cc) > 0 {
		m.SetHeader("Cc", email.Cc...)
	}
	if email.Bcc != nil && len(email.Bcc) > 0 {
		m.SetHeader("Bcc", email.Bcc...)
	}
	if email.Subject == "" {
		return errors.New("邮件标题不能为空！")
	}
	m.SetHeader("Subject", email.Subject)
	//m.SetBody("text/plain", "Hello!")
	m.SetBody("text/html", email.Body)
	for _, attach := range email.Attach {
		m.Attach(attach)
	}
	// Send the email to Bob, Cora and Dan.
	return client.dialer.DialAndSend(m)
}

// SendSampleEmail 发送简单邮件
func (client *emailClient) SendSampleEmail(to, title, content string) error {
	email := &Email{
		From:    client.dialer.Username,
		To:      []string{to},
		Subject: title,
		Body:    content,
	}
	return client.SendEmail(email)
}

// SendTemplateEmail 发送模板邮件
func (client *emailClient) SendTemplateEmail(to, title, templateName string, params interface{}) error {
	b := bytes.NewBuffer([]byte{})
	emailTemplate.ExecuteTemplate(b, templateName, params)
	return client.SendSampleEmail(to, title, b.String())
}
