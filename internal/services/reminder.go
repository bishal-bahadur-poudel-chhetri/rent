package services

import (
	"context"
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type ReminderService struct {
	reminderRepo *repositories.ReminderRepository
}

func NewReminderService(reminderRepo *repositories.ReminderRepository) *ReminderService {
	return &ReminderService{
		reminderRepo: reminderRepo,
	}
}

// CreateReminder creates a new reminder
func (s *ReminderService) CreateReminder(ctx context.Context, reminder *models.Reminder) error {
	// Set initial next_due_date to start_date
	reminder.NextDueDate = reminder.StartDate
	return s.reminderRepo.CreateReminder(ctx, reminder)
}

// AcknowledgeReminder acknowledges a reminder, updates the start date to now, and calculates the next due date
func (s *ReminderService) AcknowledgeReminder(ctx context.Context, reminderID, userID int) error {
	// Get the reminder
	reminder, err := s.reminderRepo.GetReminderByID(ctx, reminderID)
	if err != nil {
		return err
	}

	// Current time to use for acknowledgment and new start date
	currentTime := time.Now().UTC()

	// Create acknowledgement
	ack := &models.ReminderAcknowledgement{
		ReminderID:     reminderID,
		UserID:         userID,
		AcknowledgedAt: currentTime,
	}
	if err := s.reminderRepo.CreateAcknowledgement(ctx, ack); err != nil {
		return err
	}

	// Update the start date to current time and calculate the next due date from this new start date
	nextDueDate := s.reminderRepo.CalculateNextDueDate(reminder, currentTime)
	return s.reminderRepo.UpdateStartDateAndNextDueDate(ctx, reminderID, currentTime, nextDueDate)
}

// GetDueReminders gets all reminders that are due
func (s *ReminderService) GetDueReminders(ctx context.Context, limit, offset int, reminderType string) ([]models.Reminder, error) {
	return s.reminderRepo.GetDueReminders(ctx, limit, offset, reminderType)
}

// GetRemindersByVehicleID gets all reminders for a vehicle
func (s *ReminderService) GetRemindersByVehicleID(ctx context.Context, vehicleID int) ([]models.Reminder, error) {
	return s.reminderRepo.GetRemindersByVehicleID(ctx, vehicleID)
}

// GetRemindersByVehicleIDWithVehicleDetails gets all reminders for a vehicle with vehicle details
func (s *ReminderService) GetRemindersByVehicleIDWithVehicleDetails(ctx context.Context, vehicleID int) ([]models.ReminderWithVehicle, error) {
	return s.reminderRepo.GetRemindersByVehicleIDWithVehicleDetails(ctx, vehicleID)
}

// GetDueRemindersWithVehicleDetails gets all reminders that are due with vehicle details
func (s *ReminderService) GetDueRemindersWithVehicleDetails(ctx context.Context, limit, offset int, reminderType string) ([]models.ReminderWithVehicle, error) {
	return s.reminderRepo.GetDueRemindersWithVehicleDetails(ctx, limit, offset, reminderType)
}

// GetReminderHistory gets the acknowledgement history for a reminder
func (s *ReminderService) GetReminderHistory(ctx context.Context, reminderID int) ([]models.ReminderAcknowledgement, error) {
	return s.reminderRepo.GetAcknowledgements(ctx, reminderID)
}

// UpdateReminder updates an existing reminder
func (s *ReminderService) UpdateReminder(ctx context.Context, reminder *models.Reminder) error {
	return s.reminderRepo.UpdateReminder(ctx, reminder)
}

// DeleteReminder soft deletes a reminder
func (s *ReminderService) DeleteReminder(ctx context.Context, id int) error {
	return s.reminderRepo.DeleteReminder(ctx, id)
}
