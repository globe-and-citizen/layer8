package models

type Client struct {
	ID                   string `gorm:"column:id; not null" json:"id"`
	Secret               string `gorm:"column:secret" json:"secret"`
	Name                 string `gorm:"column:name" json:"name"`
	RedirectURI          string `gorm:"column:redirect_uri" json:"redirect_uri"`
	BackendURI           string `gorm:"column:backend_uri" json:"backend_uri"`
	Username             string `gorm:"column:username; unique; not null" json:"username"`
	Salt                 string `gorm:"column:salt; not null" json:"salt"`
	X509CertificateBytes []byte `gorm:"column:x509_certificate_bytes" json:"x509_certificate_bytes"`
}

func CreateClient(id, secret, name, redirect_uri string) Client {
	return Client{
		ID:          id,
		Secret:      secret,
		Name:        name,
		RedirectURI: redirect_uri,
	}
}

func (c *Client) TableName() string {
	return "clients"
}
