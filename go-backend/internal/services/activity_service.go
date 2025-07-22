// File: internal/services/activity_service.go

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/Hitesh-Sisara/GoNextAuth/internal/database"
	"github.com/Hitesh-Sisara/GoNextAuth/internal/models"

	"github.com/gin-gonic/gin"
)

type ActivityService struct{}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

// LogActivity logs user activity
func (a *ActivityService) LogActivity(ctx context.Context, userID int, activityType string, c *gin.Context, metadata map[string]interface{}) error {
	db := database.GetDB()

	// Get IP address
	var ipAddress *string
	if c != nil {
		ip := getClientIP(c)
		if ip != "" {
			ipAddress = &ip
		}
	}

	// Get user agent
	var userAgent *string
	if c != nil {
		ua := c.GetHeader("User-Agent")
		if ua != "" {
			userAgent = &ua
		}
	}

	// Convert metadata to JSON
	var metadataJSON []byte
	if metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			metadataJSON = nil
		}
	}

	// Insert activity log
	query := `
		INSERT INTO user_activity_logs (user_id, activity_type, ip_address, user_agent, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := db.Exec(ctx, query, userID, activityType, ipAddress, userAgent, metadataJSON)
	if err != nil {
		return err
	}

	// Update user's last activity
	updateQuery := `UPDATE users SET last_activity_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err = db.Exec(ctx, updateQuery, userID)
	return err
}

// GetUserActivity gets user activity logs
func (a *ActivityService) GetUserActivity(ctx context.Context, userID int, limit int, offset int) ([]models.UserActivityLog, error) {
	db := database.GetDB()

	query := `
		SELECT id, user_id, activity_type, ip_address, user_agent, metadata, created_at
		FROM user_activity_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.UserActivityLog
	for rows.Next() {
		var activity models.UserActivityLog
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.IPAddress,
			&activity.UserAgent,
			&metadataJSON,
			&activity.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Parse metadata JSON
		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &activity.Metadata)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// CleanupOldActivity removes old activity logs (older than specified days)
func (a *ActivityService) CleanupOldActivity(ctx context.Context, olderThanDays int) error {
	db := database.GetDB()

	query := `DELETE FROM user_activity_logs WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'`
	_, err := db.Exec(ctx, fmt.Sprintf(query, olderThanDays))
	return err
}

// getClientIP gets the real client IP address
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" && net.ParseIP(xri) != nil {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil && net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
