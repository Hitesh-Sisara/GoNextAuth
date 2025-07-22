// File: internal/services/email_service.go

package services

import (
	"fmt"
	"net/smtp"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

// getBaseStyles returns common CSS styles for all emails
func (e *EmailService) getBaseStyles() string {
	return fmt.Sprintf(`
		<style>
			* {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
			}
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
				line-height: 1.6;
				color: #1f2937;
				background: linear-gradient(135deg, #f8fafc 0%%, #e2e8f0 100%%);
				padding: 20px;
			}
			.email-container {
				max-width: 600px;
				margin: 0 auto;
				background: #ffffff;
				border-radius: 16px;
				overflow: hidden;
				box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
			}
			.header {
				background: linear-gradient(135deg, %s 0%%, %s 100%%);
				padding: 40px 30px;
				text-align: center;
				color: white;
				position: relative;
			}
			.header::before {
				content: '';
				position: absolute;
				top: 0;
				left: 0;
				right: 0;
				bottom: 0;
				background: url('data:image/svg+xml;charset=UTF-8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="%%23ffffff" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="%%23ffffff" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="%%23ffffff" opacity="0.15"/><circle cx="20" cy="80" r="0.5" fill="%%23ffffff" opacity="0.15"/></pattern></defs><rect width="100%%" height="100%%" fill="url(%%23grain)"/></svg>');
				pointer-events: none;
			}
			.header h1 {
				font-size: 28px;
				font-weight: 700;
				margin-bottom: 8px;
				position: relative;
				z-index: 1;
			}
			.header p {
				font-size: 16px;
				opacity: 0.9;
				position: relative;
				z-index: 1;
			}
			.logo {
				width: 60px;
				height: 60px;
				margin: 0 auto 20px;
				background: rgba(255, 255, 255, 0.2);
				border-radius: 12px;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 24px;
				font-weight: bold;
				position: relative;
				z-index: 1;
			}
			.content {
				padding: 40px 30px;
			}
			.otp-container {
				background: linear-gradient(135deg, #f8fafc 0%%, #f1f5f9 100%%);
				border: 2px dashed %s;
				border-radius: 12px;
				padding: 30px;
				text-align: center;
				margin: 30px 0;
				position: relative;
			}
			.otp-label {
				font-size: 14px;
				color: #64748b;
				margin-bottom: 10px;
				text-transform: uppercase;
				letter-spacing: 1px;
				font-weight: 600;
			}
			.otp-code {
				font-size: 36px;
				font-weight: 800;
				color: %s;
				letter-spacing: 8px;
				font-family: 'Courier New', monospace;
				text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
			}
			.info-box {
				background: linear-gradient(135deg, #fef3cd 0%%, #fde68a 100%%);
				border-left: 4px solid #f59e0b;
				border-radius: 8px;
				padding: 20px;
				margin: 25px 0;
			}
			.info-box.security {
				background: linear-gradient(135deg, #fef2f2 0%%, #fecaca 100%%);
				border-left-color: #ef4444;
			}
			.info-box h3 {
				font-size: 16px;
				margin-bottom: 10px;
				color: #92400e;
			}
			.info-box.security h3 {
				color: #b91c1c;
			}
			.info-box ul {
				list-style: none;
				margin: 0;
				padding: 0;
			}
			.info-box li {
				margin: 8px 0;
				padding-left: 20px;
				position: relative;
			}
			.info-box li:before {
				content: '•';
				color: #f59e0b;
				font-size: 18px;
				position: absolute;
				left: 0;
			}
			.info-box.security li:before {
				color: #ef4444;
			}
			.feature-grid {
				display: grid;
				gap: 20px;
				margin: 30px 0;
			}
			.feature-card {
				background: linear-gradient(135deg, #f8fafc 0%%, #f1f5f9 100%%);
				border: 1px solid #e2e8f0;
				border-radius: 12px;
				padding: 24px;
				transition: all 0.3s ease;
			}
			.feature-card:hover {
				transform: translateY(-2px);
				box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
			}
			.feature-icon {
				font-size: 32px;
				margin-bottom: 12px;
				display: block;
			}
			.feature-card h3 {
				font-size: 18px;
				margin-bottom: 8px;
				color: #1f2937;
			}
			.feature-card p {
				color: #64748b;
				font-size: 14px;
			}
			.cta-button {
				display: inline-block;
				background: linear-gradient(135deg, %s 0%%, %s 100%%);
				color: white !important;
				text-decoration: none;
				padding: 14px 32px;
				border-radius: 8px;
				font-weight: 600;
				font-size: 16px;
				text-align: center;
				margin: 20px 0;
				box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
				transition: all 0.3s ease;
			}
			.cta-button:hover {
				transform: translateY(-1px);
				box-shadow: 0 8px 15px -3px rgba(0, 0, 0, 0.1);
			}
			.footer {
				background: #f8fafc;
				padding: 30px;
				text-align: center;
				border-top: 1px solid #e2e8f0;
			}
			.footer p {
				color: #64748b;
				font-size: 13px;
				margin: 5px 0;
			}
			.social-links {
				margin: 20px 0;
			}
			.social-links a {
				display: inline-block;
				width: 36px;
				height: 36px;
				background: %s;
				color: white;
				text-decoration: none;
				border-radius: 50%%;
				line-height: 36px;
				margin: 0 8px;
				font-size: 16px;
			}
			@media only screen and (max-width: 600px) {
				.email-container {
					margin: 10px;
					border-radius: 12px;
				}
				.header, .content, .footer {
					padding: 24px 20px;
				}
				.otp-code {
					font-size: 28px;
					letter-spacing: 4px;
				}
				.header h1 {
					font-size: 24px;
				}
			}
		</style>
	`, e.config.Branding.PrimaryColor, e.config.Branding.SecondaryColor,
		e.config.Branding.PrimaryColor, e.config.Branding.PrimaryColor,
		e.config.Branding.PrimaryColor, e.config.Branding.SecondaryColor,
		e.config.Branding.PrimaryColor)
}

// getEmailHeader returns the common header HTML
func (e *EmailService) getEmailHeader(title, subtitle string) string {
	logoSection := ""
	if e.config.Branding.LogoURL != "" {
		logoSection = fmt.Sprintf(`<img src="%s" alt="%s" style="width: 60px; height: 60px; margin-bottom: 20px;">`,
			e.config.Branding.LogoURL, e.config.Branding.Name)
	} else {
		logoSection = fmt.Sprintf(`<div class="logo">%s</div>`, string(e.config.Branding.Name[0]))
	}

	return fmt.Sprintf(`
		<div class="header">
			%s
			<h1>%s</h1>
			<p>%s</p>
		</div>
	`, logoSection, title, subtitle)
}

// getEmailFooter returns the common footer HTML
func (e *EmailService) getEmailFooter() string {
	return fmt.Sprintf(`
		<div class="footer">
			<div class="social-links">
				<a href="%s">🌐</a>
				<a href="mailto:%s">✉️</a>
			</div>
			<p><strong>%s</strong></p>
			<p>This is an automated email. Please do not reply to this message.</p>
			<p>If you need assistance, contact us at <a href="mailto:%s" style="color: %s;">%s</a></p>
			<p style="margin-top: 20px; font-size: 12px;">
				© 2025 %s. All rights reserved.
			</p>
		</div>
	`, e.config.Branding.Website, e.config.Branding.SupportEmail, e.config.Branding.Name,
		e.config.Branding.SupportEmail, e.config.Branding.PrimaryColor, e.config.Branding.SupportEmail,
		e.config.Branding.Name)
}

// SendEmailVerificationOTP sends an email verification OTP
func (e *EmailService) SendEmailVerificationOTP(email, otpCode string) error {
	subject := fmt.Sprintf("Welcome to %s - Verify Your Email", e.config.Branding.Name)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Email Verification</title>
			%s
		</head>
		<body>
			<div class="email-container">
				%s
				<div class="content">
					<h2 style="color: #1f2937; margin-bottom: 20px;">Verify Your Email Address</h2>
					<p style="font-size: 16px; color: #4b5563; margin-bottom: 25px;">
						Thank you for joining %s! We're excited to have you on board. Please verify your email address by entering the verification code below:
					</p>

					<div class="otp-container">
						<div class="otp-label">Your Verification Code</div>
						<div class="otp-code">%s</div>
					</div>

					<div class="info-box">
						<h3>📋 Important Information</h3>
						<ul>
							<li>This verification code will expire in %d minutes</li>
							<li>Keep this code confidential and don't share it with anyone</li>
							<li>If you didn't create an account with us, please ignore this email</li>
						</ul>
					</div>

					<p style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
						Having trouble? Contact our support team at 
						<a href="mailto:%s" style="color: %s; text-decoration: none;">%s</a>
					</p>
				</div>
				%s
			</div>
		</body>
		</html>
	`, e.getBaseStyles(),
		e.getEmailHeader("Email Verification", fmt.Sprintf("Complete your %s registration", e.config.Branding.Name)),
		e.config.Branding.Name, otpCode, e.config.OTP.Expiry,
		e.config.Branding.SupportEmail, e.config.Branding.PrimaryColor, e.config.Branding.SupportEmail,
		e.getEmailFooter())

	textBody := fmt.Sprintf(`
		Welcome to %s!

		Thank you for joining %s! We're excited to have you on board.

		Please verify your email address by entering the following verification code:

		%s

		Important Information:
		- This verification code will expire in %d minutes
		- Keep this code confidential and don't share it with anyone
		- If you didn't create an account with us, please ignore this email

		Having trouble? Contact our support team at %s

		This is an automated email from %s. Please do not reply to this message.
	`, e.config.Branding.Name, e.config.Branding.Name, otpCode, e.config.OTP.Expiry,
		e.config.Branding.SupportEmail, e.config.Branding.Name)

	return e.sendEmail(email, subject, htmlBody, textBody)
}

// SendPasswordResetOTP sends a password reset OTP
func (e *EmailService) SendPasswordResetOTP(email, otpCode string) error {
	subject := fmt.Sprintf("%s - Password Reset Request", e.config.Branding.Name)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Password Reset</title>
			%s
		</head>
		<body>
			<div class="email-container">
				%s
				<div class="content">
					<h2 style="color: #1f2937; margin-bottom: 20px;">Password Reset Request</h2>
					<p style="font-size: 16px; color: #4b5563; margin-bottom: 25px;">
						We received a request to reset your password for your %s account. If you made this request, please use the verification code below to reset your password:
					</p>

					<div class="otp-container">
						<div class="otp-label">Your Reset Code</div>
						<div class="otp-code">%s</div>
					</div>

					<div class="info-box security">
						<h3>🔐 Security Notice</h3>
						<ul>
							<li>This reset code will expire in %d minutes</li>
							<li>Never share this code with anyone, including %s support</li>
							<li>If you didn't request this password reset, please ignore this email</li>
							<li>Your account remains secure and no changes have been made</li>
						</ul>
					</div>

					<p style="color: #6b7280; font-size: 14px; margin-top: 25px;">
						<strong>Didn't request a password reset?</strong> If you continue to receive these emails, 
						please contact our security team at 
						<a href="mailto:%s" style="color: %s; text-decoration: none;">%s</a>
					</p>
				</div>
				%s
			</div>
		</body>
		</html>
	`, e.getBaseStyles(),
		e.getEmailHeader("Password Reset", "Secure your account"),
		e.config.Branding.Name, otpCode, e.config.OTP.Expiry, e.config.Branding.Name,
		e.config.Branding.SupportEmail, e.config.Branding.PrimaryColor, e.config.Branding.SupportEmail,
		e.getEmailFooter())

	textBody := fmt.Sprintf(`
		Password Reset Request

		We received a request to reset your password for your %s account.

		Please use the following verification code to reset your password:

		%s

		Security Notice:
		- This reset code will expire in %d minutes
		- Never share this code with anyone, including %s support
		- If you didn't request this password reset, please ignore this email
		- Your account remains secure and no changes have been made

		If you continue to receive these emails and didn't request them, please contact our security team at %s

		This is an automated email from %s. Please do not reply to this message.
	`, e.config.Branding.Name, otpCode, e.config.OTP.Expiry, e.config.Branding.Name,
		e.config.Branding.SupportEmail, e.config.Branding.Name)

	return e.sendEmail(email, subject, htmlBody, textBody)
}

// SendLoginOTP sends a login OTP
func (e *EmailService) SendLoginOTP(email, otpCode string) error {
	subject := fmt.Sprintf("%s - Your Sign-in Code", e.config.Branding.Name)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Sign-in Code</title>
			%s
		</head>
		<body>
			<div class="email-container">
				%s
				<div class="content">
					<h2 style="color: #1f2937; margin-bottom: 20px;">Your Sign-in Code</h2>
					<p style="font-size: 16px; color: #4b5563; margin-bottom: 25px;">
						Someone requested a sign-in code for your %s account. Enter the code below to complete your sign-in:
					</p>

					<div class="otp-container">
						<div class="otp-label">Your Sign-in Code</div>
						<div class="otp-code">%s</div>
					</div>

					<div class="info-box">
						<h3>⏰ Quick Reminder</h3>
						<ul>
							<li>This sign-in code will expire in %d minutes</li>
							<li>Only use this code if you're trying to sign in right now</li>
							<li>If you didn't request this code, please ignore this email</li>
						</ul>
					</div>

					<p style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
						Need help signing in? Contact us at 
						<a href="mailto:%s" style="color: %s; text-decoration: none;">%s</a>
					</p>
				</div>
				%s
			</div>
		</body>
		</html>
	`, e.getBaseStyles(),
		e.getEmailHeader("Sign In", fmt.Sprintf("Access your %s account", e.config.Branding.Name)),
		e.config.Branding.Name, otpCode, e.config.OTP.Expiry,
		e.config.Branding.SupportEmail, e.config.Branding.PrimaryColor, e.config.Branding.SupportEmail,
		e.getEmailFooter())

	textBody := fmt.Sprintf(`
		Sign in to %s

		Someone requested a sign-in code for your %s account.

		Your sign-in code:

		%s

		Quick Reminder:
		- This sign-in code will expire in %d minutes
		- Only use this code if you're trying to sign in right now
		- If you didn't request this code, please ignore this email

		Need help signing in? Contact us at %s

		This is an automated email from %s. Please do not reply to this message.
	`, e.config.Branding.Name, e.config.Branding.Name, otpCode, e.config.OTP.Expiry,
		e.config.Branding.SupportEmail, e.config.Branding.Name)

	return e.sendEmail(email, subject, htmlBody, textBody)
}

// SendWelcomeEmail sends a welcome email after successful signup
func (e *EmailService) SendWelcomeEmail(email, firstName string) error {
	subject := fmt.Sprintf("Welcome to %s, %s! 🎉", e.config.Branding.Name, firstName)

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Welcome</title>
			%s
		</head>
		<body>
			<div class="email-container">
				%s
				<div class="content">
					<h2 style="color: #1f2937; margin-bottom: 20px;">🎉 Welcome Aboard, %s!</h2>
					<p style="font-size: 16px; color: #4b5563; margin-bottom: 30px;">
						Your email has been successfully verified and your %s account is now active! 
						We're thrilled to have you join our community.
					</p>

					<div class="feature-grid">
						<div class="feature-card">
							<span class="feature-icon">🚀</span>
							<h3>Get Started</h3>
							<p>Explore our platform and discover all the amazing features we have built for you.</p>
						</div>
						<div class="feature-card">
							<span class="feature-icon">📚</span>
							<h3>Learn & Grow</h3>
							<p>Access our comprehensive documentation, tutorials, and resources to maximize your experience.</p>
						</div>
						<div class="feature-card">
							<span class="feature-icon">💬</span>
							<h3>Get Support</h3>
							<p>Our dedicated support team is here to help you succeed every step of the way.</p>
						</div>
					</div>

					<div style="text-align: center; margin: 35px 0;">
						<a href="%s" class="cta-button">Start Exploring →</a>
					</div>

					<div class="info-box">
						<h3>🌟 What's Next?</h3>
						<ul>
							<li>Complete your profile to personalize your experience</li>
							<li>Explore our features and discover what works best for you</li>
							<li>Join our community and connect with other users</li>
							<li>Check out our getting started guide for tips and tricks</li>
						</ul>
					</div>

					<p style="color: #6b7280; font-size: 14px; text-align: center; margin-top: 30px;">
						Questions? We're here to help! Reach out to us at 
						<a href="mailto:%s" style="color: %s; text-decoration: none;">%s</a>
					</p>
				</div>
				%s
			</div>
		</body>
		</html>
	`, e.getBaseStyles(),
		e.getEmailHeader("Welcome!", fmt.Sprintf("You're now part of the %s family", e.config.Branding.Name)),
		firstName, e.config.Branding.Name, e.config.Branding.Website,
		e.config.Branding.SupportEmail, e.config.Branding.PrimaryColor, e.config.Branding.SupportEmail,
		e.getEmailFooter())

	textBody := fmt.Sprintf(`
		Welcome to %s!

		Hi %s, we're excited to have you on board!

		Your email has been successfully verified and your %s account is now active! 
		We're thrilled to have you join our community.

		What's Next?
		🚀 Get Started - Explore our platform and discover all the amazing features we have built for you
		📚 Learn & Grow - Access our comprehensive documentation, tutorials, and resources  
		💬 Get Support - Our dedicated support team is here to help you succeed

		Next Steps:
		- Complete your profile to personalize your experience
		- Explore our features and discover what works best for you
		- Join our community and connect with other users
		- Check out our getting started guide for tips and tricks

		Visit us at: %s

		Questions? We're here to help! Reach out to us at %s

		Thanks for choosing %s!
		This is an automated email from %s. Please do not reply to this message.
	`, e.config.Branding.Name, firstName, e.config.Branding.Name, e.config.Branding.Website,
		e.config.Branding.SupportEmail, e.config.Branding.Name, e.config.Branding.Name)

	return e.sendEmail(email, subject, htmlBody, textBody)
}

// sendEmail sends an email using AWS SES SMTP
func (e *EmailService) sendEmail(to, subject, htmlBody, textBody string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.AWS.SES.FromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", textBody)
	m.AddAlternative("text/html", htmlBody)

	// Create SMTP dialer with correct SMTP credentials
	d := gomail.NewDialer(
		e.config.AWS.SES.SMTPHost,
		e.config.AWS.SES.SMTPPort,
		e.config.AWS.SES.SMTPUsername, // Use SMTP username
		e.config.AWS.SES.SMTPPassword, // Use SMTP password
	)

	// Use STARTTLS
	d.TLSConfig = nil
	d.Auth = smtp.PlainAuth(
		"",
		e.config.AWS.SES.SMTPUsername, // Use SMTP username
		e.config.AWS.SES.SMTPPassword, // Use SMTP password
		e.config.AWS.SES.SMTPHost,
	)

	// Send email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
