package models

import "time"

type PasswordResetTokenData struct {
	ID        uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	Username  string    `gorm:"column:username; unique; not null" json:"username"`
	Token     []byte    `gorm:"column:token; unique; not null" json:"token"`
	ExpiresAt time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
}

func (PasswordResetTokenData) TableName() string {
	return "password_reset_tokens"
}
