// File: internal/email/templates/data.go

package templates

// Common email data structures
type VerificationEmailData struct {
	BrandName     string
	OTPCode       string
	ExpiryMinutes int
	SupportEmail  string
	PrimaryColor  string
}

type PasswordResetEmailData struct {
	BrandName     string
	OTPCode       string
	ExpiryMinutes int
	SupportEmail  string
	PrimaryColor  string
}

type LoginOTPEmailData struct {
	BrandName     string
	OTPCode       string
	ExpiryMinutes int
	SupportEmail  string
	PrimaryColor  string
}

type WelcomeEmailData struct {
	BrandName    string
	FirstName    string
	Website      string
	SupportEmail string
	PrimaryColor string
}
