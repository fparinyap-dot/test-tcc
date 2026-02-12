package models

import "time"

type QueueTicket struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	QueueNumber string    `gorm:"type:varchar(2);not null;index" json:"queue_number"`
	IssuedAt    time.Time `gorm:"not null;autoCreateTime" json:"issued_at"`
}
