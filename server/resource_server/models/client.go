package models

type Client struct {
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	Name        string `json:"name"`
	RedirectURI string `json:"redirect_uri"`
	Username    string `gorm:"column:username; unique; not null" json:"username"`
	Password    string `gorm:"column:password; not null" json:"password"`
	Salt        string `gorm:"column:salt; not null" json:"salt"`
}

func (Client) TableName() string {
	return "clients"
}