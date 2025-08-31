package models

import "time"

type UserMetadata struct {
	ID                    uint      `gorm:"column:id; primaryKey; not null" json:"id"`
	DisplayName           string    `gorm:"column:display_name; not null" json:"display_name"`
	Color                 string    `gorm:"column:color; not null" json:"color"`
	Bio                   string    `gorm:"column:bio; not null" json:"bio"`
	IsEmailVerified       bool      `gorm:"column:is_email_verified; not null" json:"is_email_verified"`
	IsPhoneNumberVerified bool      `gorm:"column:is_phone_number_verified; not null" json:"is_phone_number_verified"`
	CreatedAt             time.Time `gorm:"column:created_at; default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at; default:CURRENT_TIMESTAMP; autoUpdateTime" json:"updated_at"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
