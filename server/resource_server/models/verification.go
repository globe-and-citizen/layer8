package models

import "time"

type EmailVerificationData struct {
	ID               uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserId           string    `gorm:"column:user_id; unique; not null" json:"user_id"`
	VerificationCode string    `gorm:"column:verification_code; not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
}

func (EmailVerificationData) TableName() string {
	return "email_verification_data"
}
