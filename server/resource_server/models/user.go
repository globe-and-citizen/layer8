package models

type User struct {
	ID               uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	Username         string `gorm:"column:username; unique; not null" json:"username"`
	Salt             string `gorm:"column:salt; not null" json:"salt"`
	EmailProof       []byte `gorm:"column:email_proof; not null" json:"email_proof"`
	VerificationCode string `gorm:"column:verification_code; not null" json:"verification_code"`
	ZkKeyPairId      uint   `gorm:"column:zk_key_pair_id; not null" json:"zk_key_pair_id"`
	PublicKey        []byte `gorm:"column:public_key; not null" json:"public_key"`
	IterationCount   int    `gorm:"column:iteration_count;" json:"iteration_count"`
	ServerKey        string `gorm:"column:server_key;" json:"server_key"`
	StoredKey        string `gorm:"column:stored_key;" json:"stored_key"`
}

func (User) TableName() string {
	return "users"
}
