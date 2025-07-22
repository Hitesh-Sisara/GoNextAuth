// File: internal/services/auth_service_google.go (Updated with Account Merging)

package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// GetGoogleClientID returns the Google client ID for frontend use
func (a *AuthService) GetGoogleClientID() string {
	return a.config.Google.ClientID
}

// GetGoogleRedirectURL returns the Google redirect URL
func (a *AuthService) GetGoogleRedirectURL() string {
	return a.config.Google.RedirectURL
}

// GoogleCallbackAuth handles Google OAuth callback with authorization code
// GoogleCallbackAuth handles Google OAuth callback with authorization code
func (a *AuthService) GoogleCallbackAuth(ctx context.Context, code string, c *gin.Context) (*models.AuthResponse, error) {
	fmt.Printf("Starting Google callback auth with code: %s...\n", code[:20])

	// Exchange code for access token
	tokenInfo, err := a.googleService.ExchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	fmt.Printf("Token exchange successful, access token: %s...\n", tokenInfo.AccessToken[:20])

	// Get user info using access token (this implicitly validates the token)
	googleUser, err := a.googleService.GetUserInfo(ctx, tokenInfo.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	fmt.Printf("User info retrieved for: %s\n", googleUser.Email)

	return a.processGoogleUser(ctx, googleUser, c)
}

// GoogleAuth handles Google OAuth authentication with access token
func (a *AuthService) GoogleAuth(ctx context.Context, req models.GoogleAuthRequest, c *gin.Context) (*models.AuthResponse, error) {
	fmt.Printf("Starting Google auth with access token: %s...\n", req.AccessToken[:20])

	// Verify Google access token first
	valid, err := a.googleService.VerifyAccessToken(ctx, req.AccessToken)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid Google access token")
	}

	// Get user info from Google
	googleUser, err := a.googleService.GetUserInfo(ctx, req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}

	return a.processGoogleUser(ctx, googleUser, c)
}

// processGoogleUser processes Google user info with account merging support
func (a *AuthService) processGoogleUser(ctx context.Context, googleUser *GoogleUserInfo, c *gin.Context) (*models.AuthResponse, error) {
	db := database.GetDB()

	// Validate email domain if needed
	if !a.isValidEmailDomain(googleUser.Email) {
		return nil, fmt.Errorf("email domain not allowed")
	}

	var user models.User
	var userExists bool
	var isAccountMerge bool

	// Start transaction for data consistency
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// First try to find by Google ID
	googleIDQuery := `
		SELECT id, email, first_name, last_name, phone, is_email_verified, is_active, 
		       google_id, avatar_url, auth_provider, last_activity_at, created_at, updated_at
		FROM users WHERE google_id = $1
	`
	err = tx.QueryRow(ctx, googleIDQuery, googleUser.ID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.IsEmailVerified, &user.IsActive, &user.GoogleID, &user.AvatarURL,
		&user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == nil {
		userExists = true
	} else if err == pgx.ErrNoRows {
		// Try to find by email (account merging scenario)
		emailQuery := `
			SELECT id, email, first_name, last_name, phone, is_email_verified, is_active, 
			       google_id, avatar_url, auth_provider, last_activity_at, created_at, updated_at
			FROM users WHERE email = $1
		`
		err = tx.QueryRow(ctx, emailQuery, googleUser.Email).Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
			&user.IsEmailVerified, &user.IsActive, &user.GoogleID, &user.AvatarURL,
			&user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
		)

		if err == nil {
			userExists = true
			isAccountMerge = true

			// Merge Google auth method to existing account
			var newAuthProvider string
			if user.AuthProvider == "email" {
				newAuthProvider = "email,google" // Multiple auth methods
			} else if !strings.Contains(user.AuthProvider, "google") {
				newAuthProvider = user.AuthProvider + ",google"
			} else {
				newAuthProvider = user.AuthProvider // Google already linked
			}

			// Update existing user with Google info
			updateQuery := `
				UPDATE users 
				SET google_id = $1, 
				    avatar_url = COALESCE($2, avatar_url), 
				    is_email_verified = TRUE, 
				    auth_provider = $3,
				    first_name = COALESCE(NULLIF($4, ''), first_name),
				    last_name = COALESCE(NULLIF($5, ''), last_name),
				    last_activity_at = CURRENT_TIMESTAMP
				WHERE id = $6
			`
			_, err = tx.Exec(ctx, updateQuery,
				googleUser.ID,
				googleUser.Picture,
				newAuthProvider,
				googleUser.GivenName,
				googleUser.FamilyName,
				user.ID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to merge Google account: %w", err)
			}

			// Update user struct with new values
			user.GoogleID = &googleUser.ID
			if googleUser.Picture != "" {
				user.AvatarURL = &googleUser.Picture
			}
			user.IsEmailVerified = true
			user.AuthProvider = newAuthProvider

			// Update name if Google provides better info and current names are generic
			if googleUser.GivenName != "" && (user.FirstName == "User" || user.FirstName == "") {
				user.FirstName = googleUser.GivenName
			}
			if googleUser.FamilyName != "" && (user.LastName == "Google" || user.LastName == "") {
				user.LastName = googleUser.FamilyName
			}

		} else if err != pgx.ErrNoRows {
			return nil, fmt.Errorf("failed to check existing user: %w", err)
		}
	} else {
		return nil, fmt.Errorf("failed to check Google user: %w", err)
	}

	// If user doesn't exist, create new user
	if !userExists {
		// Parse first and last name from Google response
		firstName := googleUser.GivenName
		lastName := googleUser.FamilyName

		// Fallback to parsing full name if given/family names are empty
		if firstName == "" && googleUser.Name != "" {
			nameParts := strings.Fields(googleUser.Name)
			if len(nameParts) > 0 {
				firstName = nameParts[0]
				if len(nameParts) > 1 {
					lastName = strings.Join(nameParts[1:], " ")
				}
			}
		}

		// Ensure we have at least some name
		if firstName == "" {
			firstName = "User"
		}
		if lastName == "" {
			lastName = ""
		}

		insertQuery := `
			INSERT INTO users (email, first_name, last_name, is_email_verified, is_active, 
			                  google_id, avatar_url, auth_provider, last_activity_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
			RETURNING id, email, first_name, last_name, phone, is_email_verified, is_active, 
			          google_id, avatar_url, auth_provider, last_activity_at, created_at, updated_at
		`
		err = tx.QueryRow(ctx, insertQuery,
			googleUser.Email, firstName, lastName,
			true, true, googleUser.ID, googleUser.Picture, "google").Scan(
			&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
			&user.IsEmailVerified, &user.IsActive, &user.GoogleID, &user.AvatarURL,
			&user.AuthProvider, &user.LastActivityAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create Google user: %w", err)
		}
	} else {
		// Update last activity for existing user
		_, err = tx.Exec(ctx, "UPDATE users SET last_activity_at = CURRENT_TIMESTAMP WHERE id = $1", user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update last activity: %w", err)
		}
		user.LastActivityAt = time.Now()
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
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

	// Log activity with account merge info
	activityType := models.ActivityGoogleLogin
	metadata := map[string]interface{}{
		"auth_provider": "google",
		"login_method":  "oauth",
		"is_new_user":   !userExists,
		"user_agent":    c.GetHeader("User-Agent"),
		"ip_address":    c.ClientIP(),
	}

	if isAccountMerge {
		metadata["account_merged"] = true
		metadata["previous_auth_method"] = strings.Split(user.AuthProvider, ",")[0]
	}

	a.activityService.LogActivity(ctx, user.ID, activityType, c, metadata)

	return &models.AuthResponse{
		User:         user.ToUserResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    a.jwtService.GetTokenExpiry(),
	}, nil
}

// isValidEmailDomain validates email domain if restrictions are needed
func (a *AuthService) isValidEmailDomain(email string) bool {
	// Add your domain restrictions here if needed
	// For example, to only allow certain domains:
	// allowedDomains := []string{"yourdomain.com", "gmail.com"}
	// domain := strings.Split(email, "@")[1]
	// for _, allowed := range allowedDomains {
	//     if domain == allowed {
	//         return true
	//     }
	// }
	// return false

	// For now, allow all domains
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
