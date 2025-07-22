// File: internal/email/templates/login_otp.go

package templates

import "fmt"

type LoginOTPTemplate struct{}

func NewLoginOTPTemplate() *LoginOTPTemplate {
	return &LoginOTPTemplate{}
}

func (l *LoginOTPTemplate) GetTitle() string {
	return "Sign-in Code"
}

func (l *LoginOTPTemplate) GetSubtitle() string {
	return "Access your account"
}

func (l *LoginOTPTemplate) RenderContent(data interface{}) string {
	d := data.(LoginOTPEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.SupportEmail, d.PrimaryColor, d.SupportEmail)
}

func (l *LoginOTPTemplate) RenderText(data interface{}) string {
	d := data.(LoginOTPEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.SupportEmail, d.BrandName)
}
