// File: internal/database/migrations.go

package database

import (
	"context"
	"log"
)

func RunMigrations() {
	db := GetDB()
	ctx := context.Background()

	// Create users table
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			phone VARCHAR(20),
			is_email_verified BOOLEAN DEFAULT FALSE,
			is_active BOOLEAN DEFAULT TRUE,
			google_id VARCHAR(255),
			avatar_url TEXT,
			auth_provider VARCHAR(50) DEFAULT 'email',
			last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	// Create refresh tokens table
	createRefreshTokensTable := `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			token_hash VARCHAR(255) NOT NULL,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			is_revoked BOOLEAN DEFAULT FALSE
		);
	`

	// Create OTP table
	createOTPTable := `
		CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			otp_code VARCHAR(10) NOT NULL,
			otp_type VARCHAR(50) NOT NULL, -- 'email_verification', 'password_reset', 'login'
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			is_used BOOLEAN DEFAULT FALSE
		);
	`

	// Create user activity logs table
	createUserActivityLogsTable := `
		CREATE TABLE IF NOT EXISTS user_activity_logs (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			activity_type VARCHAR(100) NOT NULL,
			ip_address INET,
			user_agent TEXT,
			metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	// Create indexes
	createIndexes := `
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
		CREATE INDEX IF NOT EXISTS idx_users_last_activity ON users(last_activity_at);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
		CREATE INDEX IF NOT EXISTS idx_otps_email ON otps(email);
		CREATE INDEX IF NOT EXISTS idx_otps_code ON otps(otp_code);
		CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps(expires_at);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON user_activity_logs(user_id);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_activity_type ON user_activity_logs(activity_type);
		CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON user_activity_logs(created_at);
	`

	// Create trigger for updating updated_at
	createTrigger := `
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS update_users_updated_at ON users;
		CREATE TRIGGER update_users_updated_at
			BEFORE UPDATE ON users
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`

	// Migration for existing installations - Add new columns to users table
	addNewColumnsToUsers := `
		DO $$ BEGIN
			-- Add phone column if it doesn't exist
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='phone') THEN
				ALTER TABLE users ADD COLUMN phone VARCHAR(20);
			END IF;

			-- Add last_activity_at column if it doesn't exist
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='last_activity_at') THEN
				ALTER TABLE users ADD COLUMN last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
			END IF;

			-- Add google_id column if it doesn't exist
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='google_id') THEN
				ALTER TABLE users ADD COLUMN google_id VARCHAR(255);
			END IF;

			-- Add avatar_url column if it doesn't exist
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='avatar_url') THEN
				ALTER TABLE users ADD COLUMN avatar_url TEXT;
			END IF;

			-- Add auth_provider column if it doesn't exist
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='auth_provider') THEN
				ALTER TABLE users ADD COLUMN auth_provider VARCHAR(50) DEFAULT 'email';
			END IF;

			-- Make password_hash nullable for Google OAuth users
			ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;
		END $$;
	`

	// Execute migrations in order
	queries := []string{
		createUsersTable,
		addNewColumnsToUsers,
		createRefreshTokensTable,
		createOTPTable,
		createUserActivityLogsTable,
		createIndexes,
		createTrigger,
	}

	for i, query := range queries {
		_, err := db.Exec(ctx, query)
		if err != nil {
			log.Fatalf("Failed to run migration %d: %v", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
	log.Println("✅ Enhanced authentication system ready with:")
	log.Println("   - Phone number support")
	log.Println("   - Google OAuth integration")
	log.Println("   - User activity tracking")
	log.Println("   - Multi-step authentication")
}
