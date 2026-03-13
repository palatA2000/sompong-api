package models

import "time"

type ScheduledMessage struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID            uint      `gorm:"not null;index" json:"chat_id"`
	Chat              Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	CreatedByUserID   *uint     `gorm:"index" json:"created_by_user_id"`
	CreatedByUser     *User     `gorm:"foreignKey:CreatedByUserID" json:"created_by_user,omitempty"`
	MessageText       string    `gorm:"type:text;not null" json:"message_text"`
	MessageType       string    `gorm:"type:varchar;not null" json:"message_type"`
	Timezone          string    `gorm:"type:varchar;not null" json:"timezone"`
	StartAt           time.Time `gorm:"not null" json:"start_at"`
	EndAt             *time.Time `gorm:"" json:"end_at"`
	FrequencyType     string    `gorm:"type:varchar;not null" json:"frequency_type"`
	FrequencyInterval *int      `gorm:"" json:"frequency_interval"`
	DaysOfWeek        *string   `gorm:"type:varchar" json:"days_of_week"`
	DayOfMonth        *int      `gorm:"" json:"day_of_month"`
	TimeOfDay         *string   `gorm:"type:varchar" json:"time_of_day"`
	CronExpression    *string   `gorm:"type:text" json:"cron_expression"`
	IsActive          bool      `gorm:"not null;default:true" json:"is_active"`
	NextRunAt         *time.Time `gorm:"" json:"next_run_at"`
	CreatedAt         time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time `gorm:"not null;default:now()" json:"updated_at"`

	Substitutions []ScheduledMessageSubstitution `gorm:"foreignKey:ScheduledMessageID" json:"substitutions,omitempty"`
}
