package models

type User struct {
	ID       uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	Username string `gorm:"column:username; unique; not null" json:"username"`

	EmailVerificationCode string `gorm:"column:verification_code; not null" json:"verification_code"`
	EmailZkProof          []byte `gorm:"column:email_proof; not null" json:"email_proof"`
	EmailZkKeyPairId      uint   `gorm:"column:zk_key_pair_id; not null" json:"zk_key_pair_id"`

	PhoneNumberVerificationCode string `gorm:"column:phone_number_verification_code; not null" json:"phone_number_verification_code"`
	PhoneNumberZkProof          []byte `gorm:"column:phone_number_zk_proof; not null" json:"phone_number_zk_proof"`
	PhoneNumberZkPairID         uint   `gorm:"column:phone_number_zk_pair_id; not null" json:"phone_number_zk_pair_id"`

	PublicKey []byte `gorm:"column:public_key; not null" json:"public_key"`

	Salt           string `gorm:"column:salt; not null" json:"salt"`
	IterationCount int    `gorm:"column:iteration_count;" json:"iteration_count"`
	ServerKey      string `gorm:"column:server_key;" json:"server_key"`
	StoredKey      string `gorm:"column:stored_key;" json:"stored_key"`

	TelegramSessionIDHash []byte `gorm:"column:telegram_session_id_hash; not null" json:"telegram_session_id_hash"`
}

func (User) TableName() string {
	return "users"
}
