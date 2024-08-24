package models

type ZkSnarksKeyPair struct {
	ID              uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	ProvingKey      []byte `gorm:"column:proving_key; not null" json:"proving_key"`
	VerificationKey []byte `gorm:"column:verification_key; not null" json:"verification_key"`
}

func (ZkSnarksKeyPair) TableName() string {
	return "zk_snarks_key_pairs"
}
