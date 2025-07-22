// File: internal/models/otp.go

package models

import "time"

// OTP represents an OTP in the system
type OTP struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	OTPCode   string    `json:"otp_code" db:"otp_code"`
	OTPType   string    `json:"otp_type" db:"otp_type"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
}

// OTP Types
const (
	OTPTypeEmailVerification = "email_verification"
	OTPTypePasswordReset     = "password_reset"
	OTPTypeLogin             = "login" // Added for OTP-based login
)

// ResendOTPRequest represents resend OTP request payload
type ResendOTPRequest struct {
	Email   string `json:"email" binding:"required,email" example:"user@example.com"`
	OTPType string `json:"otp_type" binding:"required,oneof=email_verification password_reset login" example:"email_verification"`
}
