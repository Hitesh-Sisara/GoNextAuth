// File: internal/email/templates/password_reset.go

package templates

import "fmt"

type PasswordResetTemplate struct{}

func NewPasswordResetTemplate() *PasswordResetTemplate {
	return &PasswordResetTemplate{}
}

func (p *PasswordResetTemplate) GetTitle() string {
	return "Password Reset"
}

func (p *PasswordResetTemplate) GetSubtitle() string {
	return "Secure your account"
}

func (p *PasswordResetTemplate) RenderContent(data interface{}) string {
	d := data.(PasswordResetEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.BrandName, d.SupportEmail, d.PrimaryColor, d.SupportEmail)
}

func (p *PasswordResetTemplate) RenderText(data interface{}) string {
	d := data.(PasswordResetEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.BrandName, d.SupportEmail, d.BrandName)
}
