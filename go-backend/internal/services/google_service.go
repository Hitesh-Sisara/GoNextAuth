// File: internal/services/google_service.go

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/config"
)

type GoogleService struct {
	config *config.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleTokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

type GoogleTokenValidationResponse struct {
	Audience  string `json:"aud"`
	ClientID  string `json:"client_id"`
	ExpiresIn string `json:"expires_in"`
	Scope     string `json:"scope"`
	Email     string `json:"email"`
	Error     string `json:"error,omitempty"`
}

type GoogleErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewGoogleService(cfg *config.Config) *GoogleService {
	return &GoogleService{config: cfg}
}

// GetUserInfo gets user information from Google using access token
func (g *GoogleService) GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	fmt.Printf("Getting user info from Google with token: %s...\n", accessToken[:20])

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Google user info response status: %d\n", resp.StatusCode)
	fmt.Printf("Google user info response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		var errorResp GoogleErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
			return nil, fmt.Errorf("google API error: %s - %s", errorResp.Error, errorResp.ErrorDescription)
		}
		return nil, fmt.Errorf("google API returned status: %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// Validate that email is verified
	if !userInfo.VerifiedEmail {
		return nil, fmt.Errorf("google email is not verified")
	}

	fmt.Printf("Successfully retrieved user info for: %s\n", userInfo.Email)
	return &userInfo, nil
}

// VerifyAccessToken verifies if the access token is valid
func (g *GoogleService) VerifyAccessToken(ctx context.Context, accessToken string) (bool, error) {
	fmt.Printf("Verifying Google access token: %s...\n", accessToken[:20])

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v1/tokeninfo", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to verify token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Token verification response status: %d\n", resp.StatusCode)
	fmt.Printf("Token verification response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Token verification failed with status: %d\n", resp.StatusCode)
		return false, nil
	}

	var validation GoogleTokenValidationResponse
	if err := json.Unmarshal(body, &validation); err != nil {
		return false, fmt.Errorf("failed to decode validation response: %w", err)
	}

	if validation.Error != "" {
		fmt.Printf("Token validation error: %s\n", validation.Error)
		return false, nil
	}

	// Verify the token is for our client
	isValid := validation.Audience == g.config.Google.ClientID || validation.ClientID == g.config.Google.ClientID
	fmt.Printf("Token validation result: %t (aud: %s, client_id: %s, our_client: %s)\n",
		isValid, validation.Audience, validation.ClientID, g.config.Google.ClientID)

	return isValid, nil
}

// ExchangeCodeForToken exchanges authorization code for access token
func (g *GoogleService) ExchangeCodeForToken(ctx context.Context, code string) (*GoogleTokenInfo, error) {
	fmt.Printf("Exchanging code for token: %s...\n", code[:20])
	fmt.Printf("Using redirect URI: %s\n", g.config.Google.RedirectURL)

	data := url.Values{}
	data.Set("client_id", g.config.Google.ClientID)
	data.Set("client_secret", g.config.Google.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", g.config.Google.RedirectURL)

	fmt.Printf("Token exchange request data: %s\n", data.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Token exchange response status: %d\n", resp.StatusCode)
	fmt.Printf("Token exchange response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		// Try to get error details
		var errorResp GoogleErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error != "" {
			// Provide more specific error messages
			switch errorResp.Error {
			case "invalid_grant":
				if strings.Contains(errorResp.ErrorDescription, "Code was already redeemed") ||
					strings.Contains(errorResp.ErrorDescription, "Bad Request") {
					return nil, fmt.Errorf("authorization code has already been used or expired")
				}
				return nil, fmt.Errorf("invalid or expired authorization code")
			case "invalid_client":
				return nil, fmt.Errorf("invalid client credentials")
			case "invalid_request":
				return nil, fmt.Errorf("invalid request parameters")
			default:
				return nil, fmt.Errorf("google token exchange failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
			}
		}
		return nil, fmt.Errorf("google token exchange failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenInfo GoogleTokenInfo
	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	fmt.Printf("Successfully exchanged code for access token: %s...\n", tokenInfo.AccessToken[:20])
	return &tokenInfo, nil
}
