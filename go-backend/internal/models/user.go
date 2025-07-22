// File: internal/models/user.go

package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              int       `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	Phone           *string   `json:"phone" db:"phone"`
	IsEmailVerified bool      `json:"is_email_verified" db:"is_email_verified"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	GoogleID        *string   `json:"-" db:"google_id"`
	AvatarURL       *string   `json:"avatar_url" db:"avatar_url"`
	AuthProvider    string    `json:"auth_provider" db:"auth_provider"`
	LastActivityAt  time.Time `json:"last_activity_at" db:"last_activity_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UserActivityLog represents user activity tracking
type UserActivityLog struct {
	ID           int                    `json:"id" db:"id"`
	UserID       int                    `json:"user_id" db:"user_id"`
	ActivityType string                 `json:"activity_type" db:"activity_type"`
	IPAddress    *string                `json:"ip_address" db:"ip_address"`
	UserAgent    *string                `json:"user_agent" db:"user_agent"`
	Metadata     map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// UserResponse represents user data returned to client (without sensitive info)
type UserResponse struct {
	ID              int       `json:"id"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           *string   `json:"phone"`
	IsEmailVerified bool      `json:"is_email_verified"`
	IsActive        bool      `json:"is_active"`
	AvatarURL       *string   `json:"avatar_url"`
	AuthProvider    string    `json:"auth_provider"`
	LastActivityAt  time.Time `json:"last_activity_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// SignupRequest represents signup request payload
type SignupRequest struct {
	Email     string `json:"email" binding:"required,email" example:"user@example.com"`
	Password  string `json:"password" binding:"required,min=8" example:"password123"`
	FirstName string `json:"first_name" binding:"required,min=2" example:"John"`
	LastName  string `json:"last_name" binding:"required,min=2" example:"Doe"`
	Phone     string `json:"phone" binding:"omitempty" example:"+911234567890"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// EmailLoginRequest represents email-only login request
type EmailLoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// OTPLoginRequest represents OTP login request
type OTPLoginRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// GoogleAuthRequest represents Google OAuth request
type GoogleAuthRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// CompleteSignupRequest represents final signup step
type CompleteSignupRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required,min=2"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// ForgotPasswordRequest represents forgot password request payload
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ResetPasswordRequest represents reset password request payload
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP         string `json:"otp" binding:"required,len=6" example:"123456"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newpassword123"`
}

// VerifyEmailRequest represents email verification request payload
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	OTP   string `json:"otp" binding:"required,len=6" example:"123456"`
}

// RefreshTokenRequest represents refresh token request payload
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int          `json:"expires_in"` // seconds
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsRevoked bool      `json:"is_revoked" db:"is_revoked"`
}

// Activity Types
const (
	ActivityLogin         = "login"
	ActivityLogout        = "logout"
	ActivitySignup        = "signup"
	ActivityEmailVerify   = "email_verify"
	ActivityPasswordReset = "password_reset"
	ActivityProfileUpdate = "profile_update"
	ActivityOTPLogin      = "otp_login"
	ActivityGoogleLogin   = "google_login"
)

// ToUserResponse converts User to UserResponse
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Phone:           u.Phone,
		IsEmailVerified: u.IsEmailVerified,
		IsActive:        u.IsActive,
		AvatarURL:       u.AvatarURL,
		AuthProvider:    u.AuthProvider,
		LastActivityAt:  u.LastActivityAt,
		CreatedAt:       u.CreatedAt,
	}
}
