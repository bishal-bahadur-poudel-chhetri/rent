package models

import "time"

type SystemSetting struct {
	SettingID    int       `json:"setting_id"`
	SettingKey   string    `json:"setting_key"`
	SettingValue bool      `json:"setting_value"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SystemSettingsResponse struct {
	EnableRegistration bool `json:"enable_registration"`
	EnableLogin        bool `json:"enable_login"`
}
