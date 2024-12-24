package models

type ZkSnarksKeyPair struct {
	ID           uint   `gorm:"primaryKey; unique; autoIncrement; not null" json:"id"`
	ProvingKey   []byte `gorm:"column:proving_key; not null" json:"proving_key"`
	VerifyingKey []byte `gorm:"column:verifying_key; not null" json:"verifying_key"`
}

func (ZkSnarksKeyPair) TableName() string {
	return "zk_snarks_key_pairs"
}
