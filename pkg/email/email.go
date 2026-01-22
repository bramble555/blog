package email

import (
	"html/template"
	"math/rand"
	"strings"
	"time"

	"github.com/bramble555/blog/global"
	"gopkg.in/gomail.v2"
)

// 全局随机数生成器
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// generate4DigitCode 生成4位随机验证码
func generate4DigitCode() int {
	return rng.Intn(9000) + 1000 // [1000, 9999]
}

// SendEmail 发送邮件
func SendEmail(toEmail string) (int, error) {
	// 模板内容
	templateContent := `
		<h1>欢迎来到 Bramble 博客</h1>
		<p>您的验证码是: <strong>{{.Code}}</strong></p>
		<p>有效时间: {{.Expiration}} 分钟</p>
	`

	// 动态数据
	code := generate4DigitCode()
	data := map[string]any{
		"Code":       code,
		"Expiration": 5, // 验证码有效期 5 分钟
	}

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", global.Config.Email.User) // 设置发件人
	m.SetHeader("To", toEmail)                    // 设置收件人
	m.SetHeader("Subject", "您的验证码")               // 设置主题

	// 使用模板生成邮件正文
	tmpl, err := template.New("emailTemplate").Parse(templateContent)
	if err != nil {
		global.Log.Errorf("解析模板失败: %s", err)
		return -1, err
	}

	// 填充模板数据
	body := new(strings.Builder)
	if err := tmpl.Execute(body, data); err != nil {
		global.Log.Errorf("执行模板失败: %s", err)
		return -1, err
	}
	m.SetBody("text/html", body.String()) // 设置邮件内容为 HTML

	// 设置邮件服务拨号器
	d := gomail.NewDialer(
		global.Config.Email.Host,
		global.Config.Email.Port,
		global.Config.Email.User,
		global.Config.Email.Password,
	)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		global.Log.Errorf("发送邮件失败: %s", err)
		return -1, err
	}

	// 日志记录成功发送
	global.Log.Infof("邮件发送成功，验证码为: %d", code)
	return code, nil
}
