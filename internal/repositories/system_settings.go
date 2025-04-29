package repositories

import (
	"database/sql"
	"renting/internal/models"
)

type SystemSettingsRepository struct {
	db *sql.DB
}

func NewSystemSettingsRepository(db *sql.DB) *SystemSettingsRepository {
	return &SystemSettingsRepository{db: db}
}

func (r *SystemSettingsRepository) GetSystemSettings() (*models.SystemSettingsResponse, error) {
	query := `
		SELECT setting_key, setting_value 
		FROM system_settings 
		WHERE setting_key IN ('enable_registration', 'enable_login')
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := &models.SystemSettingsResponse{
		EnableRegistration: true, // Default values
		EnableLogin:        true,
	}

	for rows.Next() {
		var key string
		var value bool
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		switch key {
		case "enable_registration":
			settings.EnableRegistration = value
		case "enable_login":
			settings.EnableLogin = value
		}
	}

	return settings, nil
}

func (r *SystemSettingsRepository) UpdateSystemSetting(key string, value bool) error {
	query := `
		UPDATE system_settings 
		SET setting_value = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE setting_key = $2
	`
	_, err := r.db.Exec(query, value, key)
	return err
}
