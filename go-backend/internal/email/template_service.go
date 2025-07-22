// File: internal/email/template_service.go

package email

import (
	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/email/templates"
)

type TemplateService struct {
	config       *config.Config
	baseTemplate *templates.BaseTemplate
}

func NewTemplateService(cfg *config.Config) *TemplateService {
	return &TemplateService{
		config:       cfg,
		baseTemplate: templates.NewBaseTemplate(cfg),
	}
}

// GenerateVerificationEmail generates HTML and text versions of verification email
func (t *TemplateService) GenerateVerificationEmail(otpCode string) (htmlBody, textBody string) {
	data := templates.VerificationEmailData{
		BrandName:     t.config.Branding.Name,
		OTPCode:       otpCode,
		ExpiryMinutes: t.config.OTP.Expiry,
		SupportEmail:  t.config.Branding.SupportEmail,
		PrimaryColor:  t.config.Branding.PrimaryColor,
	}

	htmlTemplate := templates.NewVerificationTemplate()
	htmlBody = t.baseTemplate.RenderHTML(htmlTemplate, data)
	textBody = htmlTemplate.RenderText(data)

	return htmlBody, textBody
}

// GeneratePasswordResetEmail generates HTML and text versions of password reset email
func (t *TemplateService) GeneratePasswordResetEmail(otpCode string) (htmlBody, textBody string) {
	data := templates.PasswordResetEmailData{
		BrandName:     t.config.Branding.Name,
		OTPCode:       otpCode,
		ExpiryMinutes: t.config.OTP.Expiry,
		SupportEmail:  t.config.Branding.SupportEmail,
		PrimaryColor:  t.config.Branding.PrimaryColor,
	}

	htmlTemplate := templates.NewPasswordResetTemplate()
	htmlBody = t.baseTemplate.RenderHTML(htmlTemplate, data)
	textBody = htmlTemplate.RenderText(data)

	return htmlBody, textBody
}

// GenerateLoginOTPEmail generates HTML and text versions of login OTP email
func (t *TemplateService) GenerateLoginOTPEmail(otpCode string) (htmlBody, textBody string) {
	data := templates.LoginOTPEmailData{
		BrandName:     t.config.Branding.Name,
		OTPCode:       otpCode,
		ExpiryMinutes: t.config.OTP.Expiry,
		SupportEmail:  t.config.Branding.SupportEmail,
		PrimaryColor:  t.config.Branding.PrimaryColor,
	}

	htmlTemplate := templates.NewLoginOTPTemplate()
	htmlBody = t.baseTemplate.RenderHTML(htmlTemplate, data)
	textBody = htmlTemplate.RenderText(data)

	return htmlBody, textBody
}

// GenerateWelcomeEmail generates HTML and text versions of welcome email
func (t *TemplateService) GenerateWelcomeEmail(firstName string) (htmlBody, textBody string) {
	data := templates.WelcomeEmailData{
		BrandName:    t.config.Branding.Name,
		FirstName:    firstName,
		Website:      t.config.Branding.Website,
		SupportEmail: t.config.Branding.SupportEmail,
		PrimaryColor: t.config.Branding.PrimaryColor,
	}

	htmlTemplate := templates.NewWelcomeTemplate()
	htmlBody = t.baseTemplate.RenderHTML(htmlTemplate, data)
	textBody = htmlTemplate.RenderText(data)

	return htmlBody, textBody
}
