package models

type UserMetadata struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

func (UserMetadata) TableName() string {
	return "user_metadata"
}
