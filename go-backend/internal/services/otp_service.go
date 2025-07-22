// File: internal/services/otp_service.go

package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
)

type OTPService struct {
	config       *config.Config
	emailService *EmailService
}

func NewOTPService(cfg *config.Config, emailService *EmailService) *OTPService {
	return &OTPService{
		config:       cfg,
		emailService: emailService,
	}
}

// GenerateOTP generates a random OTP code
func (o *OTPService) GenerateOTP() (string, error) {
	otpLength := o.config.OTP.Length
	if otpLength <= 0 {
		otpLength = 6
	}

	// Generate random OTP
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(int64(otpLength)), nil)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", otpLength)
	return fmt.Sprintf(format, n), nil
}

// SendOTP generates and sends an OTP to the specified email
func (o *OTPService) SendOTP(ctx context.Context, email, otpType string) error {
	// Generate OTP
	otpCode, err := o.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Calculate expiry time
	expiresAt := time.Now().Add(time.Duration(o.config.OTP.Expiry) * time.Minute)

	// Save OTP to database
	db := database.GetDB()
	query := `
		INSERT INTO otps (email, otp_code, otp_type, expires_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.Exec(ctx, query, email, otpCode, otpType, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Send email based on OTP type
	switch otpType {
	case models.OTPTypeEmailVerification:
		return o.emailService.SendEmailVerificationOTP(email, otpCode)
	case models.OTPTypePasswordReset:
		return o.emailService.SendPasswordResetOTP(email, otpCode)
	case models.OTPTypeLogin:
		return o.emailService.SendLoginOTP(email, otpCode)
	default:
		return fmt.Errorf("unsupported OTP type: %s", otpType)
	}
}

// VerifyOTP verifies an OTP code
func (o *OTPService) VerifyOTP(ctx context.Context, email, otpCode, otpType string) error {
	db := database.GetDB()

	// Find valid OTP
	var otp models.OTP
	query := `
		SELECT id, email, otp_code, otp_type, expires_at, created_at, is_used
		FROM otps
		WHERE email = $1 AND otp_code = $2 AND otp_type = $3 AND is_used = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := db.QueryRow(ctx, query, email, otpCode, otpType).Scan(
		&otp.ID, &otp.Email, &otp.OTPCode, &otp.OTPType,
		&otp.ExpiresAt, &otp.CreatedAt, &otp.IsUsed,
	)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	// Mark OTP as used
	updateQuery := `UPDATE otps SET is_used = TRUE WHERE id = $1`
	_, err = db.Exec(ctx, updateQuery, otp.ID)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as used: %w", err)
	}

	return nil
}

// CleanupExpiredOTPs removes expired OTPs from the database
func (o *OTPService) CleanupExpiredOTPs(ctx context.Context) error {
	db := database.GetDB()
	query := `DELETE FROM otps WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := db.Exec(ctx, query)
	return err
}

// ResendOTP resends an OTP after invalidating previous ones
func (o *OTPService) ResendOTP(ctx context.Context, email, otpType string) error {
	db := database.GetDB()

	// Invalidate previous OTPs of the same type for this email
	invalidateQuery := `
		UPDATE otps
		SET is_used = TRUE
		WHERE email = $1 AND otp_type = $2 AND is_used = FALSE
	`
	_, err := db.Exec(ctx, invalidateQuery, email, otpType)
	if err != nil {
		return fmt.Errorf("failed to invalidate previous OTPs: %w", err)
	}

	// Send new OTP
	return o.SendOTP(ctx, email, otpType)
}
