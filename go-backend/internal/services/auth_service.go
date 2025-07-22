// File: internal/services/auth_service.go

package services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	config          *config.Config
	jwtService      *JWTService
	otpService      *OTPService
	emailService    *EmailService
	googleService   *GoogleService
	activityService *ActivityService
}

func NewAuthService(cfg *config.Config, jwtService *JWTService, otpService *OTPService, emailService *EmailService, googleService *GoogleService, activityService *ActivityService) *AuthService {
	return &AuthService{
		config:          cfg,
		jwtService:      jwtService,
		otpService:      otpService,
		emailService:    emailService,
		googleService:   googleService,
		activityService: activityService,
	}
}

// InitiateEmailSignup starts the email signup process
func (a *AuthService) InitiateEmailSignup(ctx context.Context, email string) error {
	// Validate email
	if err := utils.ValidateEmail(email); err != nil {
		return err
	}

	// Check if user already exists
	var existingUserID int
	checkQuery := `SELECT id FROM users WHERE email = $1`
	err := database.GetDB().QueryRow(ctx, checkQuery, email).Scan(&existingUserID)
	if err == nil {
		return fmt.Errorf("user with this email already exists")
	}
	if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	// Send OTP for email verification
	return a.otpService.SendOTP(ctx, email, models.OTPTypeEmailVerification)
}

// CompleteSignup completes the signup process after email verification
func (a *AuthService) CompleteSignup(ctx context.Context, req models.CompleteSignupRequest, c *gin.Context) (*models.User, error) {
	db := database.GetDB()

	// Validate input
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(req.FirstName, "First name"); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(req.LastName, "Last name"); err != nil {
		return nil, err
	}

	// Validate phone if provided
	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			return nil, err
		}
	}

	// Check if user already exists
	var existingUserID int
	checkQuery := `SELECT id FROM users WHERE email = $1`
	err := db.QueryRow(ctx, checkQuery, req.Email).Scan(&existingUserID)
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}
	if err != pgx.ErrNoRows {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	var user models.User
	var phone *string
	if req.Phone != "" {
		phone = &req.Phone
	}

	insertQuery := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, is_email_verified, is_active, auth_provider)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, first_name, last_name, phone, is_email_verified, is_active, auth_provider, last_activity_at, created_at, updated_at
	`
	err = db.QueryRow(ctx, insertQuery, req.Email, hashedPassword, req.FirstName, req.LastName, phone, true, true, "email").Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.IsEmailVerified, &user.IsActive, &user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Log activity
	a.activityService.LogActivity(ctx, user.ID, models.ActivitySignup, c, map[string]interface{}{
		"auth_provider": "email",
	})

	return &user, nil
}

// Signup creates a new user account (legacy method for backward compatibility)
func (a *AuthService) Signup(ctx context.Context, req models.SignupRequest, c *gin.Context) (*models.User, error) {
	completeReq := models.CompleteSignupRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
	}
	return a.CompleteSignup(ctx, completeReq, c)
}

// Login authenticates a user and returns tokens
func (a *AuthService) Login(ctx context.Context, req models.LoginRequest, c *gin.Context) (*models.AuthResponse, error) {
	db := database.GetDB()

	// Validate input
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	// Get user by email
	var user models.User
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, is_email_verified, is_active, auth_provider, last_activity_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := db.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Phone,
		&user.IsEmailVerified, &user.IsActive, &user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate tokens
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	if err := a.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Log activity
	a.activityService.LogActivity(ctx, user.ID, models.ActivityLogin, c, map[string]interface{}{
		"auth_provider": "email",
		"login_method":  "password",
	})

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    a.jwtService.GetTokenExpiry(),
	}, nil
}

// InitiateEmailLogin starts the OTP-based email login process
func (a *AuthService) InitiateEmailLogin(ctx context.Context, req models.EmailLoginRequest) error {
	db := database.GetDB()

	// Validate email
	if err := utils.ValidateEmail(req.Email); err != nil {
		return err
	}

	// Check if user exists and is active
	var userExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_active = TRUE)`
	err := db.QueryRow(ctx, checkQuery, req.Email).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !userExists {
		return fmt.Errorf("no account found with this email address")
	}

	// Send OTP for login
	return a.otpService.SendOTP(ctx, req.Email, models.OTPTypeLogin)
}

// CompleteOTPLogin completes the OTP-based login
func (a *AuthService) CompleteOTPLogin(ctx context.Context, req models.OTPLoginRequest, c *gin.Context) (*models.AuthResponse, error) {
	db := database.GetDB()

	// Verify OTP
	if err := a.otpService.VerifyOTP(ctx, req.Email, req.OTP, models.OTPTypeLogin); err != nil {
		return nil, err
	}

	// Get user by email
	var user models.User
	query := `
		SELECT id, email, first_name, last_name, phone, is_email_verified, is_active, auth_provider, avatar_url, last_activity_at, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = TRUE
	`
	err := db.QueryRow(ctx, query, req.Email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.IsEmailVerified, &user.IsActive, &user.AuthProvider, &user.AvatarURL,
		&user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate tokens
	accessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	if err := a.storeRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Log activity
	a.activityService.LogActivity(ctx, user.ID, models.ActivityOTPLogin, c, map[string]interface{}{
		"auth_provider": "email",
		"login_method":  "otp",
	})

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    a.jwtService.GetTokenExpiry(),
	}, nil
}

// VerifyEmail verifies a user's email with OTP
func (a *AuthService) VerifyEmail(ctx context.Context, req models.VerifyEmailRequest, c *gin.Context) error {
	// Validate OTP
	if err := a.otpService.VerifyOTP(ctx, req.Email, req.OTP, models.OTPTypeEmailVerification); err != nil {
		return err
	}

	// Update user's email verification status
	db := database.GetDB()
	updateQuery := `UPDATE users SET is_email_verified = TRUE WHERE email = $1`
	result, err := db.Exec(ctx, updateQuery, req.Email)
	if err != nil {
		return fmt.Errorf("failed to update email verification status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	// Get user details for welcome email and activity logging
	var user models.User
	userQuery := `SELECT id, first_name FROM users WHERE email = $1`
	err = db.QueryRow(ctx, userQuery, req.Email).Scan(&user.ID, &user.FirstName)
	if err == nil {
		// Log activity
		a.activityService.LogActivity(ctx, user.ID, models.ActivityEmailVerify, c, nil)

		// Send welcome email (non-blocking)
		go func() {
			if err := a.emailService.SendWelcomeEmail(req.Email, user.FirstName); err != nil {
				fmt.Printf("Failed to send welcome email: %v\n", err)
			}
		}()
	}

	return nil
}

// ForgotPassword initiates password reset process
func (a *AuthService) ForgotPassword(ctx context.Context, req models.ForgotPasswordRequest) error {
	db := database.GetDB()

	// Check if user exists
	var userExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_active = TRUE)`
	err := db.QueryRow(ctx, checkQuery, req.Email).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	if !userExists {
		// Don't reveal if user exists or not for security
		return nil
	}

	// Send password reset OTP
	return a.otpService.SendOTP(ctx, req.Email, models.OTPTypePasswordReset)
}

// ResetPassword resets user's password using OTP
func (a *AuthService) ResetPassword(ctx context.Context, req models.ResetPasswordRequest, c *gin.Context) error {
	// Validate new password
	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	// Verify OTP
	if err := a.otpService.VerifyOTP(ctx, req.Email, req.OTP, models.OTPTypePasswordReset); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user's password
	db := database.GetDB()
	var userID int
	updateQuery := `UPDATE users SET password_hash = $1 WHERE email = $2 RETURNING id`
	err = db.QueryRow(ctx, updateQuery, hashedPassword, req.Email).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Revoke all refresh tokens for this user
	if err := a.revokeAllUserRefreshTokens(ctx, req.Email); err != nil {
		fmt.Printf("Failed to revoke refresh tokens: %v\n", err)
	}

	// Log activity
	a.activityService.LogActivity(ctx, userID, models.ActivityPasswordReset, c, nil)

	return nil
}

// RefreshToken generates new access token using refresh token
func (a *AuthService) RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := a.jwtService.ValidateToken(req.RefreshToken, "refresh")
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists and is not revoked
	if err := a.validateStoredRefreshToken(ctx, claims.UserID, req.RefreshToken); err != nil {
		return nil, err
	}

	// Get user details
	user, err := a.getUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Generate new access token
	newAccessToken, err := a.jwtService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := a.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token and store new one
	if err := a.revokeRefreshToken(ctx, req.RefreshToken); err != nil {
		fmt.Printf("Failed to revoke old refresh token: %v\n", err)
	}

	if err := a.storeRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    a.jwtService.GetTokenExpiry(),
	}, nil
}

// GetUserProfile returns user profile information
func (a *AuthService) GetUserProfile(ctx context.Context, userID int) (*models.UserResponse, error) {
	user, err := a.getUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// Logout revokes refresh token
func (a *AuthService) Logout(ctx context.Context, refreshToken string, c *gin.Context) error {
	// Get user ID from token for activity logging
	claims, err := a.jwtService.ValidateToken(refreshToken, "refresh")
	if err == nil {
		// Log activity (best effort)
		a.activityService.LogActivity(ctx, claims.UserID, models.ActivityLogout, c, nil)
	}

	return a.revokeRefreshToken(ctx, refreshToken)
}

// Helper functions

func (a *AuthService) storeRefreshToken(ctx context.Context, userID int, token string) error {
	db := database.GetDB()

	// Hash the token for storage
	tokenHash := a.hashToken(token)
	expiresAt := time.Now().Add(time.Duration(a.config.JWT.RefreshExpiry) * 24 * time.Hour)

	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := db.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (a *AuthService) validateStoredRefreshToken(ctx context.Context, userID int, token string) error {
	db := database.GetDB()
	tokenHash := a.hashToken(token)

	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM refresh_tokens
			WHERE user_id = $1 AND token_hash = $2 AND is_revoked = FALSE AND expires_at > CURRENT_TIMESTAMP
		)
	`
	err := db.QueryRow(ctx, query, userID, tokenHash).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if !exists {
		return fmt.Errorf("refresh token is invalid or expired")
	}

	return nil
}

func (a *AuthService) revokeRefreshToken(ctx context.Context, token string) error {
	db := database.GetDB()
	tokenHash := a.hashToken(token)

	query := `UPDATE refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1`
	_, err := db.Exec(ctx, query, tokenHash)
	return err
}

func (a *AuthService) revokeAllUserRefreshTokens(ctx context.Context, email string) error {
	db := database.GetDB()

	query := `
		UPDATE refresh_tokens
		SET is_revoked = TRUE
		WHERE user_id = (SELECT id FROM users WHERE email = $1)
	`
	_, err := db.Exec(ctx, query, email)
	return err
}

func (a *AuthService) getUserByID(ctx context.Context, userID int) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	query := `
		SELECT id, email, first_name, last_name, phone, is_email_verified, is_active, 
		       google_id, avatar_url, auth_provider, last_activity_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.IsEmailVerified, &user.IsActive, &user.GoogleID, &user.AvatarURL,
		&user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (a *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

func (a *AuthService) GetJWTSecret() string {
	return a.config.JWT.Secret
}

// GetJWTService returns the JWT service instance
func (a *AuthService) GetJWTService() *JWTService {
	return a.jwtService
}

// GetActivityService returns the Activity service instance
func (a *AuthService) GetActivityService() *ActivityService {
	return a.activityService
}
