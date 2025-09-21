package models

import "time"

type EmailVerificationData struct {
	ID               uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserId           uint      `gorm:"column:user_id; unique; not null" json:"user_id"`
	VerificationCode string    `gorm:"column:verification_code; not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
}

func (EmailVerificationData) TableName() string {
	return "email_verification_data"
}

type PhoneNumberVerificationData struct {
	ID               uint      `gorm:"primaryKey; autoIncrement; not null" json:"id"`
	UserId           uint      `gorm:"column:user_id; unique; not null" json:"user_id"`
	VerificationCode string    `gorm:"column:verification_code; not null" json:"verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at; not null" json:"expires_at"`
	ZkProof          []byte    `gorm:"column:zk_proof; not null" json:"zk_proof"`
	ZkPairID         uint      `gorm:"column:phone_number_zk_pair_id; not null" json:"zk_key_pair_id"`
}

func (PhoneNumberVerificationData) TableName() string {
	return "phone_number_verification_data"
}
