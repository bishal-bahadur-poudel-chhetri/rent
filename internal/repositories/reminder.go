package repositories

import (
	"context"
	"fmt"
	"log"
	"renting/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *ReminderRepository {
	return &ReminderRepository{db: db}
}

// CreateReminder creates a new reminder with calculated next due date
func (r *ReminderRepository) CreateReminder(ctx context.Context, reminder *models.Reminder) error {
	// Log the incoming reminder for debugging
	log.Printf("Incoming reminder: %+v", reminder)

	// Validate custom frequency
	if reminder.Frequency == models.ReminderFrequencyCustom && reminder.CustomInterval == nil {
		return fmt.Errorf("custom_interval is required for custom frequency")
	}

	// Reset NextDueDate if it matches StartDate for monthly, yearly, or custom
	if reminder.Frequency == models.ReminderFrequencyMonthly ||
		reminder.Frequency == models.ReminderFrequencyYearly ||
		reminder.Frequency == models.ReminderFrequencyCustom {
		if !reminder.NextDueDate.IsZero() && reminder.NextDueDate.Equal(reminder.StartDate) {
			log.Printf("Resetting preset NextDueDate from %s", reminder.NextDueDate.Format(time.RFC3339))
			reminder.NextDueDate = time.Time{} // Force calculation
		}
	}

	// Calculate next due date if not set
	if reminder.NextDueDate.IsZero() {
		reminder.NextDueDate = r.CalculateNextDueDate(reminder, reminder.StartDate)
		log.Printf("Calculated NextDueDate: %s", reminder.NextDueDate.Format(time.RFC3339))
	}

	return r.db.WithContext(ctx).Create(reminder).Error
}

// GetReminderByID retrieves a reminder by its ID
func (r *ReminderRepository) GetReminderByID(ctx context.Context, id int) (*models.Reminder, error) {
	var reminder models.Reminder
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&reminder).Error
	if err != nil {
		return nil, err
	}
	return &reminder, nil
}

// GetRemindersByVehicleID retrieves all reminders for a vehicle
func (r *ReminderRepository) GetRemindersByVehicleID(ctx context.Context, vehicleID int) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.WithContext(ctx).Where("vehicle_id = ?", vehicleID).Find(&reminders).Error
	return reminders, err
}

// GetDueReminders retrieves all reminders that are due
func (r *ReminderRepository) GetDueReminders(ctx context.Context, limit, offset int, reminderType string) ([]models.Reminder, error) {
	var reminders []models.Reminder
	today := time.Now().UTC().Truncate(24 * time.Hour)

	query := r.db.WithContext(ctx).
		Where("next_due_date <= ?", today)

	if reminderType != "" {
		query = query.Where("type = ?", reminderType)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(&reminders).Error
	return reminders, err
}

// UpdateNextDueDate updates the next due date for a reminder and updates the updated_at timestamp
func (r *ReminderRepository) UpdateNextDueDate(ctx context.Context, reminderID int, nextDueDate time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Reminder{}).
		Where("id = ?", reminderID).
		Updates(map[string]interface{}{
			"next_due_date": nextDueDate,
			"updated_at":    time.Now().UTC(),
		}).Error
}

// UpdateStartDateAndNextDueDate updates both the start date and next due date for a reminder
func (r *ReminderRepository) UpdateStartDateAndNextDueDate(ctx context.Context, reminderID int, startDate, nextDueDate time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.Reminder{}).
		Where("id = ?", reminderID).
		Updates(map[string]interface{}{
			"start_date":    startDate,
			"next_due_date": nextDueDate,
			"updated_at":    time.Now().UTC(),
		}).Error
}

// CreateAcknowledgement creates a new reminder acknowledgement and updates the reminder's start date and next due date
func (r *ReminderRepository) CreateAcknowledgement(ctx context.Context, ack *models.ReminderAcknowledgement) error {
	// Start a transaction to ensure atomicity
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create the acknowledgement
	if err := tx.Create(ack).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Retrieve the reminder to get its frequency
	reminder, err := r.GetReminderByID(ctx, ack.ReminderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update the reminder's start date to the acknowledgement's AcknowledgedAt
	newStartDate := ack.AcknowledgedAt
	log.Printf("Setting new start date to acknowledgement date: %s", newStartDate.Format(time.RFC3339))

	// Calculate the new next due date based on the frequency and new start date
	newNextDueDate := r.CalculateNextDueDate(reminder, newStartDate)
	log.Printf("Calculated new next due date: %s", newNextDueDate.Format(time.RFC3339))

	// Update the reminder's start date and next due date
	if err := r.UpdateStartDateAndNextDueDate(ctx, ack.ReminderID, newStartDate, newNextDueDate); err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// GetAcknowledgements retrieves all acknowledgements for a reminder
func (r *ReminderRepository) GetAcknowledgements(ctx context.Context, reminderID int) ([]models.ReminderAcknowledgement, error) {
	var acks []models.ReminderAcknowledgement
	err := r.db.WithContext(ctx).
		Where("reminder_id = ?", reminderID).
		Order("acknowledged_at DESC").
		Find(&acks).Error
	return acks, err
}

// GetRemindersByTypeAndVehicle retrieves reminders by type and vehicle
func (r *ReminderRepository) GetRemindersByTypeAndVehicle(ctx context.Context, vehicleID int, reminderType string) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.WithContext(ctx).
		Where("vehicle_id = ? AND type = ?", vehicleID, reminderType).
		Find(&reminders).Error
	return reminders, err
}

// GetRemindersByVehicleIDWithVehicleDetails gets reminders with vehicle info
func (r *ReminderRepository) GetRemindersByVehicleIDWithVehicleDetails(ctx context.Context, vehicleID int) ([]models.ReminderWithVehicle, error) {
	var reminders []models.ReminderWithVehicle
	err := r.db.WithContext(ctx).
		Table("reminders").
		Select("reminders.*, vehicles.vehicle_name, vehicles.vehicle_model, vehicles.vehicle_registration_number, vehicles.image_name").
		Joins("JOIN vehicles ON reminders.vehicle_id = vehicles.vehicle_id").
		Where("reminders.vehicle_id = ?", vehicleID).
		Find(&reminders).Error
	return reminders, err
}

// GetDueRemindersWithVehicleDetails gets due reminders with vehicle info
func (r *ReminderRepository) GetDueRemindersWithVehicleDetails(ctx context.Context, limit, offset int, reminderType string) ([]models.ReminderWithVehicle, error) {
	var reminders []models.ReminderWithVehicle
	today := time.Now().UTC().Truncate(24 * time.Hour)

	query := r.db.WithContext(ctx).
		Table("reminders").
		Select("reminders.*, vehicles.vehicle_name, vehicles.vehicle_model, vehicles.vehicle_registration_number, vehicles.image_name").
		Joins("JOIN vehicles ON reminders.vehicle_id = vehicles.vehicle_id").
		Where("reminders.next_due_date <= ?", today)

	if reminderType != "" {
		query = query.Where("reminders.type = ?", reminderType)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(&reminders).Error

	return reminders, err
}

// CalculateNextDueDate computes the next due date based on frequency
func (r *ReminderRepository) CalculateNextDueDate(reminder *models.Reminder, fromDate time.Time) time.Time {
	switch reminder.Frequency {
	case models.ReminderFrequencyMonthly:
		return fromDate.AddDate(0, 1, 0) // Add 1 month
	case models.ReminderFrequencyYearly:
		return fromDate.AddDate(1, 0, 0) // Add 1 year
	case models.ReminderFrequencyCustom:
		if reminder.CustomInterval != nil {
			return fromDate.AddDate(0, 0, *reminder.CustomInterval) // Add custom days
		}
		return fromDate
	default:
		return fromDate
	}
}

// UpdateReminder updates an existing reminder
func (r *ReminderRepository) UpdateReminder(ctx context.Context, reminder *models.Reminder) error {
	return r.db.WithContext(ctx).Save(reminder).Error
}

// DeleteReminder deletes a reminder permanently
func (r *ReminderRepository) DeleteReminder(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&models.Reminder{}, id).Error
}
