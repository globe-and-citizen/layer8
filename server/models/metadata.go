package models

type UserMetadata struct {
	ID                    uint   `gorm:"column:id; primaryKey; not null" json:"id"`
	DisplayName           string `gorm:"column:display_name; not null" json:"display_name"`
	Color                 string `gorm:"column:color; not null" json:"color"`
	Bio                   string `gorm:"column:bio; not null" json:"bio"`
	IsEmailVerified       bool   `gorm:"column:is_email_verified; not null" json:"is_email_verified"`
	IsPhoneNumberVerified bool   `gorm:"column:is_phone_number_verified; not null" json:"is_phone_number_verified"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
