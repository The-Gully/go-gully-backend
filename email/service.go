package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      bool
	SSL      bool
}

var config *EmailConfig

func Initialize(cfg EmailConfig) {
	config = &cfg
}

func GetConfig() *EmailConfig {
	return config
}

func SendEmail(to, subject, body string) error {
	if config == nil {
		return fmt.Errorf("email config not initialized")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := &gomail.Dialer{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
	}

	if config.TLS {
		d.TLSConfig = &tls.Config{
			ServerName: config.Host,
		}
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

type VerificationEmailData struct {
	Username    string
	VerifyURL   string
	FrontendURL string
}

func SendVerificationEmail(to, username, token string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	data := VerificationEmailData{
		Username:    username,
		VerifyURL:   verifyURL,
		FrontendURL: frontendURL,
	}

	body, err := ParseVerificationTemplate(data)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	subject := "Verify your email address"

	return SendEmail(to, subject, body)
}

func ParseVerificationTemplate(data VerificationEmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 4px; }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Verify your email address</h2>
        <p>Hi {{.Username}},</p>
        <p>Thanks for registering! Please verify your email address by clicking the button below:</p>
        <p><a href="{{.VerifyURL}}" class="button">Verify Email</a></p>
        <p>Or copy and paste this link into your browser:</p>
        <p><a href="{{.VerifyURL}}">{{.VerifyURL}}</a></p>
        <p>This link will expire in 24 hours.</p>
        <div class="footer">
            <p>If you didn't create an account, you can safely ignore this email.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
