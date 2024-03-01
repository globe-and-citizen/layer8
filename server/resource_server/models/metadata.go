package models

import "time"

type UserMetadata struct {
	ID        uint      `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	UserID    uint      `gorm:"column:user_id; not null" json:"user_id"`
	Key       string    `gorm:"column:key; not null" json:"key"`
	Value     string    `gorm:"column:value; not null" json:"value"`
	CreatedAt time.Time `gorm:"column:created_at; default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
