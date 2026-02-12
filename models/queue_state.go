package models

import "time"

type QueueState struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CurrentLetter string    `gorm:"type:varchar(1);not null;default:''" json:"current_letter"`
	CurrentNumber int       `gorm:"not null;default:-1" json:"current_number"`
	UpdatedAt     time.Time `json:"updated_at"`
}
