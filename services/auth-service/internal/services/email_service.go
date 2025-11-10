package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strconv"
	"strings"
)

// IEmailService defines the interface for email operations
type IEmailService interface {
	SendPasswordResetEmail(email, resetToken, resetURL string) error
}

// emailService handles email sending via SMTP or console logging
type emailService struct {
	frontendURL  string
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromAddress  string
	fromName     string
	provider     string // "smtp", "sendgrid", "mailgun", "console"
}

// NewEmailService creates a new email service instance
func NewEmailService(frontendURL, smtpHost, smtpPort, smtpUser, smtpPassword, fromAddress, fromName, provider string) IEmailService {
	return &emailService{
		frontendURL:  frontendURL,
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUser:     smtpUser,
		smtpPassword: smtpPassword,
		fromAddress:  fromAddress,
		fromName:     fromName,
		provider:     provider,
	}
}

// SendPasswordResetEmail sends a password reset email to the user
func (e *emailService) SendPasswordResetEmail(email, resetToken, resetURL string) error {
	// Construct the full reset URL
	fullResetURL := resetURL
	if fullResetURL == "" {
		fullResetURL = fmt.Sprintf("%s/reset-password?token=%s", e.frontendURL, resetToken)
	}

	// Email content
	subject := "Reset Your Password - GearBox"
	plainTextContent := e.generatePlainTextEmail(fullResetURL)
	htmlContent := e.generateHTMLEmail(fullResetURL)

	// Send email based on provider
	switch strings.ToLower(e.provider) {
	case "smtp", "sendgrid", "mailgun":
		return e.sendViaSMTP(email, subject, plainTextContent, htmlContent)
	case "console":
		fallthrough
	default:
		// Log email content for development/testing
		log.Printf("📧 Password Reset Email (Console Mode)")
		log.Printf("   To: %s", email)
		log.Printf("   Subject: %s", subject)
		log.Printf("   Reset URL: %s", fullResetURL)
		log.Printf("   Token: %s", resetToken)
		log.Printf("   --- Email Content ---")
		log.Printf("%s", plainTextContent)
		log.Printf("   --------------------")
		return nil
	}
}

// sendViaSMTP sends email using SMTP
func (e *emailService) sendViaSMTP(to, subject, plainText, html string) error {
	if e.smtpHost == "" || e.smtpPort == "" {
		return fmt.Errorf("SMTP configuration is missing. Set EMAIL_SMTP_HOST and EMAIL_SMTP_PORT")
	}

	if e.smtpUser == "" || e.smtpPassword == "" {
		return fmt.Errorf("SMTP credentials are missing. Set EMAIL_SMTP_USER and EMAIL_SMTP_PASSWORD")
	}

	// Convert port to integer
	port, err := strconv.Atoi(e.smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	// Set up authentication
	auth := smtp.PlainAuth("", e.smtpUser, e.smtpPassword, e.smtpHost)

	// Build email message
	from := e.fromAddress
	if e.fromName != "" {
		from = fmt.Sprintf("%s <%s>", e.fromName, e.fromAddress)
	}

	// Create email headers and body
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "multipart/alternative; boundary=\"boundary123\""

	// Build message body
	var body bytes.Buffer
	for k, v := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	body.WriteString("\r\n")

	// Add plain text part
	body.WriteString("--boundary123\r\n")
	body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	body.WriteString("\r\n")
	body.WriteString(plainText)
	body.WriteString("\r\n")

	// Add HTML part
	body.WriteString("--boundary123\r\n")
	body.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	body.WriteString("\r\n")
	body.WriteString(html)
	body.WriteString("\r\n")
	body.WriteString("--boundary123--\r\n")

	// Determine if we need TLS
	addr := fmt.Sprintf("%s:%d", e.smtpHost, port)

	// For ports 465, use SSL/TLS directly
	if port == 465 {
		return e.sendViaSMTPS(addr, auth, from, to, body.Bytes())
	}

	// For ports 587 and others, use STARTTLS
	return e.sendViaSMTPWithSTARTTLS(addr, auth, from, to, body.Bytes())
}

// sendViaSMTPWithSTARTTLS sends email using STARTTLS (for port 587)
func (e *emailService) sendViaSMTPWithSTARTTLS(addr string, auth smtp.Auth, from, to string, body []byte) error {
	// Connect to SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Check if server supports STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: e.smtpHost}
		if err = client.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %w", err)
	}
	if _, err = writer.Write(body); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close data connection: %w", err)
	}

	// Quit
	if err = client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP connection: %w", err)
	}

	log.Printf("✅ Password reset email sent successfully to %s", to)
	return nil
}

// sendViaSMTPS sends email using SSL/TLS directly (for port 465)
func (e *emailService) sendViaSMTPS(addr string, auth smtp.Auth, from, to string, body []byte) error {
	// Create TLS connection
	tlsConfig := &tls.Config{
		ServerName: e.smtpHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, e.smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender and recipient
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %w", err)
	}
	if _, err = writer.Write(body); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}
	if err = writer.Close(); err != nil {
		return fmt.Errorf("failed to close data connection: %w", err)
	}

	// Quit
	if err = client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP connection: %w", err)
	}

	log.Printf("✅ Password reset email sent successfully to %s", to)
	return nil
}

// generatePlainTextEmail generates plain text email content
func (e *emailService) generatePlainTextEmail(resetURL string) string {
	return fmt.Sprintf(`Hello,

You have requested to reset your password for your GearBox account.

Please click on the link below to reset your password:

%s

This link will expire in 1 hour.

If you did not request a password reset, please ignore this email. Your password will remain unchanged.

If you continue to have problems, please contact our support team.

Best regards,
GearBox Team`, resetURL)
}

// generateHTMLEmail generates HTML email content
func (e *emailService) generateHTMLEmail(resetURL string) string {
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Reset Your Password</title>
	<style>
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
			line-height: 1.6;
			color: #333;
			max-width: 600px;
			margin: 0 auto;
			padding: 20px;
			background-color: #f4f4f4;
		}
		.container {
			background-color: #ffffff;
			border-radius: 8px;
			padding: 40px;
			box-shadow: 0 2px 4px rgba(0,0,0,0.1);
		}
		.header {
			text-align: center;
			margin-bottom: 30px;
		}
		.logo {
			font-size: 24px;
			font-weight: bold;
			color: #2563eb;
			margin-bottom: 10px;
		}
		.title {
			font-size: 24px;
			font-weight: 600;
			margin-bottom: 20px;
			color: #1f2937;
		}
		.content {
			margin-bottom: 30px;
		}
		.button {
			display: inline-block;
			padding: 12px 30px;
			background-color: #2563eb;
			color: #ffffff;
			text-decoration: none;
			border-radius: 6px;
			font-weight: 600;
			margin: 20px 0;
		}
		.button:hover {
			background-color: #1d4ed8;
		}
		.footer {
			margin-top: 30px;
			padding-top: 20px;
			border-top: 1px solid #e5e7eb;
			font-size: 14px;
			color: #6b7280;
			text-align: center;
		}
		.warning {
			background-color: #fef3c7;
			border-left: 4px solid #f59e0b;
			padding: 12px;
			margin: 20px 0;
			border-radius: 4px;
		}
		.link {
			color: #2563eb;
			word-break: break-all;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<div class="logo">GearBox</div>
		</div>
		<h1 class="title">Reset Your Password</h1>
		<div class="content">
			<p>Hello,</p>
			<p>You have requested to reset your password for your GearBox account.</p>
			<p>Click the button below to reset your password:</p>
			<div style="text-align: center;">
				<a href="{{.ResetURL}}" class="button">Reset Password</a>
			</div>
			<p>Or copy and paste this link into your browser:</p>
			<p><a href="{{.ResetURL}}" class="link">{{.ResetURL}}</a></p>
			<div class="warning">
				<strong>⚠️ Important:</strong> This link will expire in 1 hour.
			</div>
			<p>If you did not request a password reset, please ignore this email. Your password will remain unchanged.</p>
			<p>If you continue to have problems, please contact our support team.</p>
		</div>
		<div class="footer">
			<p>Best regards,<br>The GearBox Team</p>
			<p style="font-size: 12px; color: #9ca3af; margin-top: 20px;">
				This is an automated email. Please do not reply to this message.
			</p>
		</div>
	</div>
</body>
</html>`

	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		// Fallback to simple HTML if template parsing fails
		return fmt.Sprintf(`<html><body><h1>Reset Your Password</h1><p>Click <a href="%s">here</a> to reset your password.</p><p>This link expires in 1 hour.</p></body></html>`, resetURL)
	}

	var buf bytes.Buffer
	data := struct {
		ResetURL string
	}{
		ResetURL: resetURL,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		// Fallback to simple HTML if template execution fails
		return fmt.Sprintf(`<html><body><h1>Reset Your Password</h1><p>Click <a href="%s">here</a> to reset your password.</p><p>This link expires in 1 hour.</p></body></html>`, resetURL)
	}

	return buf.String()
}
