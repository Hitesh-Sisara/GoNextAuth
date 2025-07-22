// File: internal/middleware/validation.go

package middleware

import (
	"net/http"
	"strings"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationErrorResponse represents validation error details
type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidateJSON validates JSON request body
func ValidateJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			var errors []ValidationErrorResponse

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrors {
					errors = append(errors, ValidationErrorResponse{
						Field:   strings.ToLower(e.Field()),
						Message: getValidationMessage(e),
						Value:   e.Value().(string),
					})
				}
			} else {
				errors = append(errors, ValidationErrorResponse{
					Message: "Invalid JSON format",
				})
			}

			c.JSON(http.StatusBadRequest, utils.APIResponse{
				Success: false,
				Message: "Validation failed",
				Error:   errors,
			})
			c.Abort()
			return
		}

		// Store the validated object in context
		c.Set("validated_data", obj)
		c.Next()
	}
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters long"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters long"
	case "len":
		return e.Field() + " must be exactly " + e.Param() + " characters long"
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	case "alphanum":
		return e.Field() + " must contain only letters and numbers"
	case "alpha":
		return e.Field() + " must contain only letters"
	case "numeric":
		return e.Field() + " must contain only numbers"
	default:
		return e.Field() + " is invalid"
	}
}

// SanitizeInput sanitizes input strings
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize common headers
		userAgent := c.GetHeader("User-Agent")
		if len(userAgent) > 512 {
			c.Header("User-Agent", userAgent[:512])
		}

		// Check for suspicious patterns
		suspiciousPatterns := []string{
			"<script", "</script>", "javascript:", "onload=", "onerror=",
			"eval(", "alert(", "document.cookie", "window.location",
		}

		for _, values := range c.Request.Header {
			for _, value := range values {
				lowerValue := strings.ToLower(value)
				for _, pattern := range suspiciousPatterns {
					if strings.Contains(lowerValue, pattern) {
						utils.SingleErrorResponse(c, http.StatusBadRequest, "Invalid request format")
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}

// ContentTypeMiddleware ensures correct content type
func ContentTypeMiddleware(expectedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, expectedType) {
				utils.SingleErrorResponse(c, http.StatusUnsupportedMediaType,
					"Content-Type must be "+expectedType)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequestSizeMiddleware limits request size
func RequestSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			utils.SingleErrorResponse(c, http.StatusRequestEntityTooLarge,
				"Request body too large")
			c.Abort()
			return
		}
		c.Next()
	}
}
