package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
)

type SystemSettingsService struct {
	repo *repositories.SystemSettingsRepository
}

func NewSystemSettingsService(repo *repositories.SystemSettingsRepository) *SystemSettingsService {
	return &SystemSettingsService{repo: repo}
}

func (s *SystemSettingsService) GetSystemSettings() (*models.SystemSettingsResponse, error) {
	return s.repo.GetSystemSettings()
}

func (s *SystemSettingsService) UpdateSystemSetting(key string, value bool) error {
	return s.repo.UpdateSystemSetting(key, value)
}
