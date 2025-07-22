// File: internal/email/templates/verification.go

package templates

import "fmt"

type VerificationTemplate struct{}

func NewVerificationTemplate() *VerificationTemplate {
	return &VerificationTemplate{}
}

func (v *VerificationTemplate) GetTitle() string {
	return "Email Verification"
}

func (v *VerificationTemplate) GetSubtitle() string {
	return "Complete your registration"
}

func (v *VerificationTemplate) RenderContent(data interface{}) string {
	d := data.(VerificationEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.SupportEmail, d.PrimaryColor, d.SupportEmail)
}

func (v *VerificationTemplate) RenderText(data interface{}) string {
	d := data.(VerificationEmailData)

	return fmt.Sprintf(`
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
	`, d.BrandName, d.BrandName, d.OTPCode, d.ExpiryMinutes, d.SupportEmail, d.BrandName)
}
