// File: internal/utils/validator.go

package utils

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateName validates first name and last name
func ValidateName(name string, fieldName string) error {
	name = strings.TrimSpace(name)

	if len(name) < 2 {
		return errors.New(fieldName + " must be at least 2 characters long")
	}

	if len(name) > 50 {
		return errors.New(fieldName + " must be less than 50 characters long")
	}

	// Only allow letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !nameRegex.MatchString(name) {
		return errors.New(fieldName + " can only contain letters, spaces, hyphens, and apostrophes")
	}

	return nil
}

// ValidateOTP validates OTP format
func ValidateOTP(otp string) error {
	if len(otp) != 6 {
		return errors.New("OTP must be exactly 6 digits")
	}

	otpRegex := regexp.MustCompile(`^\d{6}$`)
	if !otpRegex.MatchString(otp) {
		return errors.New("OTP must contain only numbers")
	}

	return nil
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) error {
	phone = strings.TrimSpace(phone)

	if len(phone) == 0 {
		return errors.New("phone number is required")
	}

	// Remove all non-digit characters except +
	cleanPhone := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")

	// Check if it starts with + and has country code
	if !strings.HasPrefix(cleanPhone, "+") {
		return errors.New("phone number must include country code (e.g., +91)")
	}

	// Remove the + for further validation
	digits := strings.TrimPrefix(cleanPhone, "+")

	// Check minimum and maximum length
	if len(digits) < 10 || len(digits) > 15 {
		return errors.New("phone number must be between 10 and 15 digits including country code")
	}

	// Check if all remaining characters are digits
	digitRegex := regexp.MustCompile(`^\d+$`)
	if !digitRegex.MatchString(digits) {
		return errors.New("phone number can only contain digits and country code")
	}

	// Basic validation for common country codes
	validCountryCodes := map[string]bool{
		"91":  true, // India
		"1":   true, // US/Canada
		"44":  true, // UK
		"33":  true, // France
		"49":  true, // Germany
		"86":  true, // China
		"81":  true, // Japan
		"82":  true, // South Korea
		"61":  true, // Australia
		"55":  true, // Brazil
		"7":   true, // Russia
		"39":  true, // Italy
		"34":  true, // Spain
		"31":  true, // Netherlands
		"46":  true, // Sweden
		"47":  true, // Norway
		"45":  true, // Denmark
		"41":  true, // Switzerland
		"43":  true, // Austria
		"32":  true, // Belgium
		"48":  true, // Poland
		"420": true, // Czech Republic
		"36":  true, // Hungary
		"351": true, // Portugal
		"30":  true, // Greece
		"358": true, // Finland
		"353": true, // Ireland
		"372": true, // Estonia
		"371": true, // Latvia
		"370": true, // Lithuania
		"65":  true, // Singapore
		"60":  true, // Malaysia
		"66":  true, // Thailand
		"84":  true, // Vietnam
		"62":  true, // Indonesia
		"63":  true, // Philippines
		"852": true, // Hong Kong
		"886": true, // Taiwan
		"971": true, // UAE
		"966": true, // Saudi Arabia
		"92":  true, // Pakistan
		"880": true, // Bangladesh
		"94":  true, // Sri Lanka
		"977": true, // Nepal
		"975": true, // Bhutan
		"960": true, // Maldives
		"27":  true, // South Africa
		"20":  true, // Egypt
		"234": true, // Nigeria
		"254": true, // Kenya
		"233": true, // Ghana
		"212": true, // Morocco
		"216": true, // Tunisia
		"213": true, // Algeria
		"218": true, // Libya
		"251": true, // Ethiopia
		"256": true, // Uganda
		"255": true, // Tanzania
		"52":  true, // Mexico
		"54":  true, // Argentina
		"56":  true, // Chile
		"57":  true, // Colombia
		"51":  true, // Peru
		"58":  true, // Venezuela
		"593": true, // Ecuador
		"595": true, // Paraguay
		"598": true, // Uruguay
	}

	// Check for valid country codes (basic check)
	foundValidCode := false
	for code := range validCountryCodes {
		if strings.HasPrefix(digits, code) {
			foundValidCode = true
			break
		}
	}

	if !foundValidCode {
		return errors.New("invalid or unsupported country code")
	}

	return nil
}
