package dto

import "time"

type CreateScheduledMessageRequest struct {
	ChatID            uint      `json:"chat_id" validate:"required"`
	CreatedByUserID   *uint     `json:"created_by_user_id"`
	MessageText       string    `json:"message_text" validate:"required"`
	MessageType       string    `json:"message_type" validate:"required"`
	Timezone          string    `json:"timezone" validate:"required"`
	StartAt           time.Time `json:"start_at" validate:"required"`
	EndAt             *time.Time `json:"end_at"`
	FrequencyType     string    `json:"frequency_type" validate:"required"`
	FrequencyInterval *int      `json:"frequency_interval"`
	DaysOfWeek        *string   `json:"days_of_week"`
	DayOfMonth        *int      `json:"day_of_month"`
	TimeOfDay         *string   `json:"time_of_day"`
	CronExpression    *string   `json:"cron_expression"`
	IsActive          bool      `json:"is_active"`
	MentionUserIDs    []uint    `json:"mention_user_ids"`
}

type UpdateScheduledMessageRequest struct {
	MessageText       *string    `json:"message_text"`
	MessageType       *string    `json:"message_type"`
	Timezone          *string    `json:"timezone"`
	StartAt           *time.Time `json:"start_at"`
	EndAt             *time.Time `json:"end_at"`
	FrequencyType     *string    `json:"frequency_type"`
	FrequencyInterval *int       `json:"frequency_interval"`
	DaysOfWeek        *string    `json:"days_of_week"`
	DayOfMonth        *int       `json:"day_of_month"`
	TimeOfDay         *string    `json:"time_of_day"`
	CronExpression    *string    `json:"cron_expression"`
	IsActive          *bool      `json:"is_active"`
	MentionUserIDs    []uint     `json:"mention_user_ids"`
}
